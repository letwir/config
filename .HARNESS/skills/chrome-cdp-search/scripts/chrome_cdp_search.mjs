#!/usr/bin/env node
/**
 * chrome_cdp_search.mjs
 * 
 * Google Chrome DevTools Protocol (CDP) Search & Web ETL Engine.
 * Enables CUI querying of local files (~/.codex, ~/.opencode, ~/.gemini, etc.)
 * AND live Web search/scraping via Chrome's headless V8 and DOM rendering runtime,
 * returning structured JSON.
 * 
 * Prerequisites: Node.js >= 22.0.0 (Native WebSocket support, zero external dependencies).
 */

import { spawn } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';

// --- ADV-07: Node.js Version Guard ---
const majorVersion = parseInt(process.versions.node.split('.')[0], 10);
if (majorVersion < 22) {
  console.error(JSON.stringify({
    status: "error",
    error: `Node.js >= 22.0.0 is required (found v${process.versions.node}) for native WebSocket support.`
  }, null, 2));
  process.exit(1);
}

// Find Chrome executable across standard Windows installation paths
function findChromeExecutable() {
  const candidates = [
    process.env.CHROME_BIN,
    "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
    "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
    path.join(process.env.LOCALAPPDATA || '', "Google\\Chrome\\Application\\chrome.exe")
  ].filter(Boolean);

  for (const p of candidates) {
    if (fs.existsSync(p)) return p;
  }
  return null;
}

function wait(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// --- Ephemeral Port Discovery & TOCTOU Guard with Backoff ---
async function waitForDevToolsActivePort(userDataDir, timeoutMs = 12000) {
  const portFile = path.join(userDataDir, 'DevToolsActivePort');
  const startTime = Date.now();
  let delay = 50;

  while (Date.now() - startTime < timeoutMs) {
    if (fs.existsSync(portFile)) {
      try {
        const raw = fs.readFileSync(portFile, 'utf8').trim();
        const lines = raw.split('\n').map(l => l.trim()).filter(Boolean);
        if (lines.length >= 2) {
          const port = parseInt(lines[0], 10);
          const wsPath = lines[1];
          if (!isNaN(port) && port > 0 && wsPath.startsWith('/devtools/browser/')) {
            try {
              const res = await fetch(`http://127.0.0.1:${port}/json/version`);
              if (res.ok) {
                const data = await res.json();
                return {
                  port,
                  wsPath,
                  browserWsUrl: data.webSocketDebuggerUrl || `ws://127.0.0.1:${port}${wsPath}`
                };
              }
            } catch (_) {}
          }
        }
      } catch (_) {}
    }
    await wait(delay);
    delay = Math.min(delay * 1.3, 300);
  }
  throw new Error(`Timeout waiting for DevToolsActivePort in ${userDataDir}`);
}

// CDP WebSocket Client with Event Subscription & Request/Response handling
class CDPClient {
  constructor(wsUrl) {
    this.wsUrl = wsUrl;
    this.ws = null;
    this.id = 1;
    this.callbacks = new Map();
    this.eventListeners = new Map();
  }

  async connect(timeoutMs = 6000) {
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        if (this.ws) {
          try { this.ws.close(); } catch (_) {}
        }
        reject(new Error(`WebSocket connection timeout to ${this.wsUrl}`));
      }, timeoutMs);

      try {
        this.ws = new WebSocket(this.wsUrl);
      } catch (err) {
        clearTimeout(timer);
        return reject(err);
      }

      this.ws.onopen = () => {
        clearTimeout(timer);
        resolve();
      };

      this.ws.onerror = (err) => {
        clearTimeout(timer);
        reject(err);
      };

      this.ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data);
          if (msg.id && this.callbacks.has(msg.id)) {
            const { resolve: reqResolve, reject: reqReject } = this.callbacks.get(msg.id);
            this.callbacks.delete(msg.id);
            if (msg.error) {
              reqReject(new Error(msg.error.message || JSON.stringify(msg.error)));
            } else {
              reqResolve(msg.result);
            }
          } else if (msg.method) {
            const listeners = this.eventListeners.get(msg.method) || [];
            listeners.forEach(fn => {
              try { fn(msg.params); } catch (_) {}
            });
          }
        } catch (_) {}
      };
    });
  }

  on(eventName, listener) {
    if (!this.eventListeners.has(eventName)) {
      this.eventListeners.set(eventName, []);
    }
    this.eventListeners.get(eventName).push(listener);
  }

  off(eventName, listener) {
    if (this.eventListeners.has(eventName)) {
      const filtered = this.eventListeners.get(eventName).filter(l => l !== listener);
      this.eventListeners.set(eventName, filtered);
    }
  }

  // --- ADV-1: Reliable Page Navigation with Event Synchronization ---
  async navigateAndWait(url, timeoutMs = 15000) {
    return new Promise(async (resolve, reject) => {
      let settled = false;
      const timer = setTimeout(() => {
        if (!settled) {
          settled = true;
          resolve({ timedOut: true });
        }
      }, timeoutMs);

      const onLoad = () => {
        if (!settled) {
          settled = true;
          clearTimeout(timer);
          this.off('Page.loadEventFired', onLoad);
          resolve({ timedOut: false });
        }
      };

      this.on('Page.loadEventFired', onLoad);

      try {
        await this.send('Page.navigate', { url });
      } catch (err) {
        if (!settled) {
          settled = true;
          clearTimeout(timer);
          this.off('Page.loadEventFired', onLoad);
          reject(err);
        }
      }
    });
  }

  async send(method, params = {}) {
    const id = this.id++;
    return new Promise((resolve, reject) => {
      this.callbacks.set(id, { resolve, reject });
      this.ws.send(JSON.stringify({ id, method, params }));
    });
  }

  close() {
    if (this.ws) {
      try {
        this.ws.close();
      } catch (_) {}
      this.ws = null;
    }
  }
}

// Global registry for process cleanup
const activeSearchers = new Set();

function cleanupAllActive() {
  for (const s of activeSearchers) {
    try { s.cleanupSync(); } catch (_) {}
  }
}
process.on('exit', cleanupAllActive);
process.on('SIGINT', () => { cleanupAllActive(); process.exit(130); });
process.on('SIGTERM', () => { cleanupAllActive(); process.exit(143); });
process.on('uncaughtException', (err) => {
  cleanupAllActive();
  console.error(JSON.stringify({ status: "fatal_error", error: err.message, stack: err.stack }));
  process.exit(1);
});

export class ChromeCDPSearcher {
  constructor(options = {}) {
    this.chromePath = options.chromePath || findChromeExecutable();
    this.tempDir = null;
    this.chromeProcess = null;
    this.port = null;
    this.pageClient = null;
    this.isClosed = false;
  }

  async start() {
    if (!this.chromePath) {
      throw new Error("Google Chrome executable not found. Please install Chrome or set CHROME_BIN.");
    }

    this.tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'chrome_cdp_search_'));
    activeSearchers.add(this);

    this.chromeProcess = spawn(this.chromePath, [
      '--headless=new',
      '--remote-debugging-port=0',
      '--remote-allow-origins=*',
      '--allow-file-access-from-files',
      '--disable-gpu',
      '--disable-extensions',
      '--disable-background-networking',
      '--disable-sync',
      '--no-first-run',
      '--no-default-browser-check',
      `--user-data-dir=${this.tempDir}`,
      'about:blank'
    ], { stdio: 'ignore' });

    const { port } = await waitForDevToolsActivePort(this.tempDir);
    this.port = port;

    const res = await fetch(`http://127.0.0.1:${port}/json/list`);
    const pages = await res.json();
    const targetPage = pages.find(p => p.type === 'page') || pages[0];

    if (!targetPage || !targetPage.webSocketDebuggerUrl) {
      throw new Error("No valid CDP page target found in Chrome instance.");
    }

    this.pageClient = new CDPClient(targetPage.webSocketDebuggerUrl);
    await this.pageClient.connect();

    await this.pageClient.send('Page.enable');
    await this.pageClient.send('Runtime.enable');
    await this.pageClient.send('DOM.enable');
  }

  // --- Local or Direct URL Search ---
  async searchFileOrUrl(targetPathOrUrl, query, options = {}) {
    const isWebUrl = /^https?:\/\//i.test(targetPathOrUrl);
    let targetUrl = targetPathOrUrl;
    let absolutePath = targetPathOrUrl;

    if (!isWebUrl) {
      absolutePath = path.resolve(targetPathOrUrl);
      if (!fs.existsSync(absolutePath)) {
        return { file: absolutePath, error: "File not found", matches: [] };
      }
      targetUrl = 'file:///' + absolutePath.replace(/\\/g, '/');
    }

    await this.pageClient.navigateAndWait(targetUrl, isWebUrl ? 15000 : 4000);
    await wait(options.waitMs || 100);

    const selector = options.selector || null;

    const expression = `(() => {
      const MAX_DEPTH = 32;
      const query = ${JSON.stringify(query ? query.toLowerCase() : '')};
      const selector = ${JSON.stringify(selector)};
      
      const isJson = window.location.href.endsWith('.json') || 
                     window.location.href.endsWith('.jsonc') || 
                     document.contentType === 'application/json';
      
      const text = document.body ? document.body.innerText : document.documentElement.innerText;
      const lines = text.split('\\n');
      
      const lineMatches = [];
      if (query) {
        lines.forEach((line, idx) => {
          if (line.toLowerCase().includes(query)) {
            lineMatches.push({
              lineNumber: idx + 1,
              content: line.trim()
            });
          }
        });
      }

      let jsonMatches = null;
      let parsedJson = null;
      if (isJson) {
        try {
          parsedJson = JSON.parse(text);
          const seen = new Set();
          
          function traverseJson(obj, currentPath, depth) {
            if (depth > MAX_DEPTH) return [];
            if (obj === null || obj === undefined) return [];
            
            const results = [];
            if (typeof obj === 'string') {
              if (query && obj.toLowerCase().includes(query)) {
                results.push({ path: currentPath, matchType: 'value', value: obj });
              }
            } else if (typeof obj === 'number' || typeof obj === 'boolean') {
              if (query && String(obj).toLowerCase().includes(query)) {
                results.push({ path: currentPath, matchType: 'value', value: obj });
              }
            } else if (typeof obj === 'object') {
              if (seen.has(obj)) return [];
              seen.add(obj);

              if (Array.isArray(obj)) {
                obj.forEach((item, index) => {
                  results.push(...traverseJson(item, currentPath + '[' + index + ']', depth + 1));
                });
              } else {
                for (const [key, val] of Object.entries(obj)) {
                  const keyPath = currentPath ? currentPath + '.' + key : key;
                  if (query && key.toLowerCase().includes(query)) {
                    results.push({
                      path: keyPath,
                      matchType: 'key',
                      value: typeof val === 'object' && val !== null ? (Array.isArray(val) ? '[Array]' : '{Object}') : val
                    });
                  }
                  results.push(...traverseJson(val, keyPath, depth + 1));
                }
              }
            }
            return results;
          }

          jsonMatches = traverseJson(parsedJson, '', 0);
        } catch (_) {}
      }

      const domMatches = [];
      if (selector) {
        try {
          const elements = document.querySelectorAll(selector);
          elements.forEach((el, idx) => {
            if (idx < 50) {
              domMatches.push({
                index: idx,
                tagName: el.tagName.toLowerCase(),
                id: el.id || null,
                className: el.className || null,
                innerText: el.innerText ? el.innerText.slice(0, 300) : ''
              });
            }
          });
        } catch (_) {}
      }

      return {
        totalLines: lines.length,
        lineMatches: lineMatches,
        jsonMatches: jsonMatches,
        domMatches: domMatches,
        isJson: !!parsedJson,
        title: document.title || null
      };
    })()`;

    const evalRes = await this.pageClient.send('Runtime.evaluate', {
      expression,
      returnByValue: true
    });

    const val = evalRes.result?.value || {};
    return {
      file: !isWebUrl ? absolutePath : null,
      url: targetUrl,
      title: val.title || null,
      query: query || null,
      selector: selector || null,
      totalLines: val.totalLines || 0,
      lineMatchCount: (val.lineMatches || []).length,
      lineMatches: val.lineMatches || [],
      jsonMatchCount: val.jsonMatches ? val.jsonMatches.length : 0,
      jsonMatches: val.jsonMatches || null,
      domMatchCount: (val.domMatches || []).length,
      domMatches: val.domMatches || null
    };
  }

  // --- ADV-2: Multi-Engine Web Search ETL with Semantic Fallback ---
  async searchWeb(keywords, options = {}) {
    const encoded = encodeURIComponent(keywords);
    const engine = options.engine || 'bing';
    let searchUrl = `https://www.bing.com/search?q=${encoded}`;
    if (engine === 'google') {
      searchUrl = `https://www.google.com/search?q=${encoded}&hl=ja`;
    } else if (engine === 'duckduckgo') {
      searchUrl = `https://duckduckgo.com/?q=${encoded}`;
    }

    await this.pageClient.navigateAndWait(searchUrl, 15000);
    await wait(options.waitMs || 800);

    const selectorOverride = options.selector || null;

    const expression = `(() => {
      const selectorOverride = ${JSON.stringify(selectorOverride)};
      const results = [];

      if (selectorOverride) {
        document.querySelectorAll(selectorOverride).forEach((el, idx) => {
          if (idx < 20) {
            results.push({
              rank: idx + 1,
              title: el.textContent ? el.textContent.slice(0, 150).trim() : '',
              url: el.querySelector('a') ? el.querySelector('a').href : (el.href || ''),
              snippet: el.textContent ? el.textContent.slice(0, 300).trim() : ''
            });
          }
        });
        return results;
      }

      // 1. Bing layout parsing (li.b_algo)
      const bingItems = document.querySelectorAll('li.b_algo');
      if (bingItems && bingItems.length > 0) {
        bingItems.forEach((el, idx) => {
          if (idx < 15) {
            const h2 = el.querySelector('h2');
            const link = el.querySelector('h2 a') || el.querySelector('a');
            const snippet = el.querySelector('.b_caption, p, .b_snippet');
            if (h2 && link && link.href) {
              results.push({
                rank: results.length + 1,
                title: h2.textContent.trim(),
                url: link.href,
                snippet: snippet ? snippet.textContent.trim() : ''
              });
            }
          }
        });
      }

      // 2. Google SERP layout fallback
      if (results.length === 0) {
        const gItems = document.querySelectorAll('div.g, div[data-hveid]');
        gItems.forEach((el, idx) => {
          if (idx < 15) {
            const titleEl = el.querySelector('h3');
            const linkEl = el.querySelector('a');
            const snippetEl = el.querySelector('div[style*="-webkit-line-clamp"], div.VwiC3b, span.aCOpRe, p');
            if (titleEl && linkEl && linkEl.href) {
              results.push({
                rank: results.length + 1,
                title: titleEl.textContent.trim(),
                url: linkEl.href,
                snippet: snippetEl ? snippetEl.textContent.trim() : ''
              });
            }
          }
        });
      }

      // 3. DuckDuckGo / Semantic Fallback (<article>, <h2> > a, <h3> > a)
      if (results.length === 0) {
        const fallbackHeadings = document.querySelectorAll('article, li[data-layout="organic"], h2 a, h3 a');
        fallbackHeadings.forEach((el, idx) => {
          if (idx < 15) {
            const a = el.tagName === 'A' ? el : el.querySelector('a');
            if (a && a.href && !a.href.startsWith('javascript:')) {
              results.push({
                rank: results.length + 1,
                title: (el.querySelector('h2, h3, [data-testid="result-title-a"]') || a).textContent.trim(),
                url: a.href,
                snippet: el.textContent ? el.textContent.slice(0, 200).trim() : ''
              });
            }
          }
        });
      }

      return results;
    })()`;

    const evalRes = await this.pageClient.send('Runtime.evaluate', {
      expression,
      returnByValue: true
    });

    const parsedItems = evalRes.result?.value || [];
    return {
      mode: "web",
      engine: engine,
      searchUrl: searchUrl,
      query: keywords,
      resultsCount: parsedItems.length,
      results: parsedItems
    };
  }

  cleanupSync() {
    if (this.isClosed) return;
    this.isClosed = true;
    activeSearchers.delete(this);
    if (this.pageClient) {
      try { this.pageClient.close(); } catch (_) {}
    }
    if (this.chromeProcess && !this.chromeProcess.killed) {
      try { this.chromeProcess.kill(); } catch (_) {}
    }
    if (this.tempDir && fs.existsSync(this.tempDir)) {
      try { fs.rmSync(this.tempDir, { recursive: true, force: true }); } catch (_) {}
    }
  }

  async close() {
    this.cleanupSync();
    await wait(100);
  }
}

// File collector
function collectFiles(targetPath, exts, maxFiles = 50) {
  const stats = fs.statSync(targetPath);
  if (!stats.isDirectory()) {
    return [targetPath];
  }

  const collected = [];
  const extSet = new Set(exts.map(e => e.toLowerCase().replace(/^\./, '')));

  function walk(dir) {
    if (collected.length >= maxFiles) return;
    try {
      const entries = fs.readdirSync(dir, { withFileTypes: true });
      for (const entry of entries) {
        if (collected.length >= maxFiles) break;
        const full = path.join(dir, entry.name);
        if (entry.isDirectory()) {
          if (!['node_modules', '.git', '.cache', 'dist', 'build'].includes(entry.name)) {
            walk(full);
          }
        } else if (entry.isFile()) {
          const ext = path.extname(entry.name).toLowerCase().replace(/^\./, '');
          if (extSet.has(ext)) {
            collected.push(full);
          }
        }
      }
    } catch (_) {}
  }

  walk(targetPath);
  return collected;
}

// CLI Execution Entry Point
async function main() {
  const args = process.argv.slice(2);
  let query = "";
  let webQuery = "";
  let webEngine = "bing";
  let selector = null;
  let target = path.join(os.homedir(), ".gemini");
  let exts = ['json', 'jsonc', 'md', 'toml', 'js', 'ts', 'html', 'txt', 'css'];
  let maxFiles = 30;
  let jsonOnly = false;

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '--query' || arg === '-q') query = args[++i];
    else if (arg === '--web-query' || arg === '-w') webQuery = args[++i];
    else if (arg === '--web-engine') webEngine = args[++i];
    else if (arg === '--target' || arg === '-t') target = args[++i];
    else if (arg === '--selector' || arg === '-s') selector = args[++i];
    else if (arg === '--ext') exts = args[++i].split(',').map(s => s.trim());
    else if (arg === '--max-files' || arg === '-m') maxFiles = parseInt(args[++i], 10);
    else if (arg === '--json-only') jsonOnly = true;
    else if (arg === '--help' || arg === '-h') {
      console.log(`
Usage: node chrome_cdp_search.mjs [options]

Modes:
  1. Local File / Directory Search:
     node chrome_cdp_search.mjs -t ~/.gemini -q "model"
  2. Live Web Search ETL:
     node chrome_cdp_search.mjs --web-query "Chrome DevTools Protocol"
  3. Direct Web URL Scraping:
     node chrome_cdp_search.mjs -t "https://example.com" -s "h1, p"

Options:
  -t, --target <path|url>   Target file, directory, or http(s) URL (default: ~/.gemini)
  -q, --query <string>      Search term for local lines or JSON key-paths
  -w, --web-query <string>  Search keywords to query live Web search engines
  --web-engine <engine>     Search engine for --web-query (bing | google | duckduckgo, default: bing)
  -s, --selector <css>      DOM CSS Selector to extract matching HTML elements
  --ext <extensions>        Comma-separated extensions to include (default: json,jsonc,md,toml,js,ts,html,txt,css)
  -m, --max-files <n>       Maximum files to search in directory (default: 30)
  --json-only               Output clean JSON without status logging
  -h, --help                Show this help message
      `);
      process.exit(0);
    }
  }

  if (!query && !webQuery && !selector && !/^https?:\/\//i.test(target)) {
    console.error(JSON.stringify({ status: "error", message: "Either --query, --web-query, --selector, or a web URL --target is required." }, null, 2));
    process.exit(1);
  }

  const searcher = new ChromeCDPSearcher();
  await searcher.start();

  try {
    // Branch A: Web Search Mode
    if (webQuery) {
      if (!jsonOnly) {
        console.error(`[CDP Web ETL] Searching "${webQuery}" via ${webEngine}...`);
      }
      const webResult = await searcher.searchWeb(webQuery, { engine: webEngine, selector });
      const output = {
        status: "success",
        engine: "Google Chrome DevTools Protocol (CDP)",
        mode: "web",
        timestamp: new Date().toISOString(),
        query: webQuery,
        searchEngine: webEngine,
        resultsCount: webResult.resultsCount,
        results: webResult.results
      };
      console.log(JSON.stringify(output, null, 2));
      return;
    }

    // Branch B: Direct Web URL Scraping
    if (/^https?:\/\//i.test(target)) {
      if (!jsonOnly) {
        console.error(`[CDP Web ETL] Navigating to URL: ${target}...`);
      }
      const urlRes = await searcher.searchFileOrUrl(target, query, { selector });
      const output = {
        status: "success",
        engine: "Google Chrome DevTools Protocol (CDP)",
        mode: "url_scrape",
        timestamp: new Date().toISOString(),
        targetUrl: target,
        result: urlRes
      };
      console.log(JSON.stringify(output, null, 2));
      return;
    }

    // Branch C: Local File / Directory Search Mode
    const resolvedTarget = path.resolve(target.replace(/^~[\\/]/, os.homedir() + path.sep));
    if (!fs.existsSync(resolvedTarget)) {
      console.error(JSON.stringify({ status: "error", message: `Target path does not exist: ${resolvedTarget}` }, null, 2));
      process.exit(1);
    }

    if (!jsonOnly) {
      console.error(`[CDP Local] Searching for "${query || selector}" in ${resolvedTarget}...`);
    }

    const files = collectFiles(resolvedTarget, exts, maxFiles);
    const matchedResults = [];
    for (const f of files) {
      const res = await searcher.searchFileOrUrl(f, query, { selector });
      if (res.lineMatchCount > 0 || res.jsonMatchCount > 0 || res.domMatchCount > 0) {
        matchedResults.push(res);
      }
    }

    const finalOutput = {
      status: "success",
      engine: "Google Chrome DevTools Protocol (CDP)",
      mode: "local",
      timestamp: new Date().toISOString(),
      query: query || null,
      selector: selector || null,
      targetPath: resolvedTarget,
      totalScannedFiles: files.length,
      matchedFilesCount: matchedResults.length,
      results: matchedResults
    };

    console.log(JSON.stringify(finalOutput, null, 2));
  } finally {
    await searcher.close();
  }
}

if (process.argv[1] && (process.argv[1].endsWith('chrome_cdp_search.mjs') || process.argv[1].endsWith('chrome_cdp_search'))) {
  main().catch(err => {
    console.error(JSON.stringify({ status: "error", error: err.message, stack: err.stack }, null, 2));
    process.exit(1);
  });
}
