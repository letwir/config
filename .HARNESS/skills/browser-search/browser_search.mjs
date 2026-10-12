#!/usr/bin/env node
/**
 * browser_search.mjs — Playwright Core → Chromium Headless Web ETL engine (Termux/Android aarch64).
 *
 *   crabcode/Qwen → browser-search SKILL → Playwright Core → Chromium Headless
 *     → DuckDuckGo検索 | JavaScript実行 | DOM解析 | ページ遷移 | 出典URL抽出
 *
 * Zero browser download: drives the Termux system Chromium via executablePath
 * and a per-run ephemeral profile (PLAYWRIGHT_BROWSERS_PATH is irrelevant here).
 * Requires: node >= 22, playwright-core (local node_modules).
 */

import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { chromium } from 'playwright-core';

const PREFIX = process.env.PREFIX || '/data/data/com.termux/files/usr';
const DEFAULT_UA = 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36';

const OPTS = {
  mode: null, target: null, query: null,
  engine: 'duckduckgo', count: 5, timeout: 30000, waitUntil: 'domcontentloaded',
  selector: null, attr: 'auto', extract: 'text', evalExpr: null,
  waitFor: null, waitMs: 0, hops: 1, locale: 'ja-JP', ua: null,
  json: true, headless: true, keepTmp: false,
};

// flag → {key} takes a value; flag → {mode} sets mode; else boolean
const FLAGS = {
  '--search': { mode: 'search' }, '--fetch': { mode: 'fetch' }, '--follow': { mode: 'follow' }, '--eval': { mode: 'eval' },
  '--engine': 'engine', '-n': 'count', '--count': 'count', '--timeout': 'timeout', '--wait-until': 'waitUntil',
  '--selector': 'selector', '--attr': 'attr', '--extract': 'extract', '--expr': 'evalExpr',
  '--wait-for': 'waitFor', '--wait-ms': 'waitMs', '--hops': 'hops', '--locale': 'locale', '--ua': 'ua',
  '--json': 'jsonTrue', '--text': 'jsonFalse', '--headful': 'headful', '--keep-tmp': 'keepTmp',
};

function usage() {
  console.log(`browser_search.mjs — Playwright Core → Chromium Headless Web ETL

検索      --search "<query>" [--engine duckduckgo|bing|google] [-n 5] [--locale ja-JP]
抽出      --fetch "<url>" [--selector "<css>" [--attr text|href|html]] [--extract text|html|links|meta|links-meta]
遷移      --follow "<url>" [--hops N] [--selector "<css>"]
JS実行    --eval "<js expression>" [--fetch "<url>"]
待機制御  --wait-for "<css>"  --wait-ms <ms>  --wait-until domcontentloaded|load|networkidle  --timeout <ms>
出力      --json (既定) | --text        他: --headful  --keep-tmp  -h

例:
  node browser_search.mjs --search "東京タワー 高さ"
  node browser_search.mjs --fetch https://example.com --selector "h1"
  node browser_search.mjs --fetch https://example.com --extract links
  node browser_search.mjs --eval "document.title" --fetch https://example.com`);
}

function fail(msg) { out({ status: 'error', error: msg }); process.exit(1); }

function out(obj) {
  if (OPTS.json) { console.log(JSON.stringify(obj, null, 2)); return; }
  for (const [k, v] of Object.entries(obj)) {
    if (Array.isArray(v)) {
      console.log(`${k}:`);
      for (const it of v) console.log('  ' + (typeof it === 'object' ? Object.entries(it).map(([a, b]) => `${a}=${b}`).join('  ') : it));
    } else console.log(`${k}: ${v}`);
  }
}

function parseArgs(argv) {
  const positional = [];
  for (let i = 0; i < argv.length; i++) {
    const a = argv[i];
    if (a === '-h' || a === '--help') { usage(); process.exit(0); }
    const spec = FLAGS[a];
    if (!spec) { positional.push(a); continue; }
    if (typeof spec === 'string') {
      const v = argv[++i];
      if (v === undefined || (v.startsWith('--') && !/^--$/.test(v))) fail(`${a} requires a value`);
      OPTS[spec] = v;
    } else if (spec.mode) OPTS.mode = spec.mode;
    else if (spec === 'jsonTrue') OPTS.json = true;
    else if (spec === 'jsonFalse') OPTS.json = false;
    else if (spec === 'headful') OPTS.headless = false;
    else if (spec === 'keepTmp') OPTS.keepTmp = true;
  }
  if (!OPTS.mode) OPTS.mode = 'search';
  if (OPTS.mode === 'search') {
    if (!OPTS.query) OPTS.query = positional.join(' ');
  } else {
    if (!OPTS.target) OPTS.target = positional[0];
    if (OPTS.mode === 'eval' && !OPTS.evalExpr) OPTS.evalExpr = positional[OPTS.target ? 1 : 0];
  }
  if ((OPTS.mode === 'search' && !OPTS.query) || (OPTS.mode !== 'search' && !OPTS.target)) fail('query (--search) or url (--fetch/--follow/--eval) is required');
  OPTS.count = parseInt(OPTS.count, 10);
  OPTS.timeout = parseInt(OPTS.timeout, 10);
  OPTS.waitMs = parseInt(OPTS.waitMs, 10) || 0;
  OPTS.hops = parseInt(OPTS.hops, 10) || 1;
  if (!Number.isInteger(OPTS.count) || OPTS.count <= 0) fail('-n must be a positive integer');
  if (!Number.isFinite(OPTS.timeout) || OPTS.timeout <= 0) fail('--timeout must be a positive integer');
}

// ---------- Chromium / context ----------
let ctx = null;
let tmpDir = null;

function findChromium() {
  const candidates = [
    process.env.CHROMIUM_PATH,
    path.join(PREFIX, 'lib/chromium/headless_shell'),
    path.join(PREFIX, 'lib/chromium/chrome'),
    path.join(PREFIX, 'bin/chromium-browser'),
  ].filter(Boolean);
  for (const p of candidates) if (fs.existsSync(p)) return p;
  return null;
}

async function newPage() {
  const executablePath = findChromium();
  if (!executablePath) fail('chromium not found (pkg install chromium)');
  tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'bs-'));
  ctx = await chromium.launchPersistentContext(tmpDir, {
    executablePath,
    headless: OPTS.headless,
    timeout: OPTS.timeout,
    locale: OPTS.locale,
    userAgent: OPTS.ua || DEFAULT_UA,
    viewport: { width: 1366, height: 900 },
    javaScriptEnabled: true,
    args: [
      '--no-sandbox', '--disable-dev-shm-usage', '--disable-gpu', '--disable-software-rasterizer',
      '--no-first-run', '--no-default-browser-check', '--disable-extensions', '--hide-scrollbars',
      '--disable-crash-reporter', '--noerrdialogs', '--disable-sync', '--mute-audio',
    ],
  });
  const page = ctx.pages()[0] || (await ctx.newPage());
  page.setDefaultTimeout(OPTS.timeout);
  page.setDefaultNavigationTimeout(OPTS.timeout);
  return page;
}

// ---------- URL normalization (出典URLの追跡パラメータ除去) ----------
function cleanUrl(u) {
  let url; try { url = new URL(u); } catch { return u; }
  try {
    if (url.hostname.includes('duckduckgo.com') && url.pathname.startsWith('/l/')) {
      const g = url.searchParams.get('uddg'); if (g) return decodeURIComponent(g);
    }
    if (url.hostname.includes('bing.com') && url.pathname.startsWith('/ck/')) {
      const a = url.searchParams.get('u');
      if (a) return Buffer.from(a.replace(/^a\d/, ''), 'base64url').toString('utf8');
    }
    if (url.hostname.includes('google.') && url.searchParams.has('url')) return url.searchParams.get('url');
  } catch { return url.toString(); }
  url.hash = '';
  for (const k of ['utm_source', 'utm_medium', 'utm_campaign', 'utm_term', 'utm_content', 'ref', 'src', 'gclid', 'ocid', 'mc_eid', 'usg', 'ved', 'source', 'sig2', 'ved'])
    url.searchParams.delete(k);
  return url.toString();
}

const firstLine = (m) => String(m || '').split('\n')[0].slice(0, 200);

// ---------- engines ----------
const mapLinks = (els) => els.map((e) => ({ title: (e.textContent || '').trim().replace(/\s+/g, ' '), url: e.href }));
const ENGINES = {
  duckduckgo: {
    name: 'duckduckgo-html',
    url: (q) => `https://html.duckduckgo.com/html/?q=${encodeURIComponent(q)}&kl=${OPTS.locale === 'ja-JP' ? 'jp-jp' : 'us-en'}`,
    item: 'a.result__a', map: mapLinks,
  },
  'duckduckgo-lite': { name: 'duckduckgo-lite', url: (q) => `https://lite.duckduckgo.com/lite/?q=${encodeURIComponent(q)}`, item: 'a.result-link', map: mapLinks },
  'duckduckgo-sqe': { name: 'duckduckgo-sqe', url: (q) => `https://duckduckgo.com/?q=${encodeURIComponent(q)}&ia=web`, item: '[data-testid="result-title-a"]', map: mapLinks },
  bing: { name: 'bing', url: (q) => `https://www.bing.com/search?q=${encodeURIComponent(q)}&setlang=${OPTS.locale === 'ja-JP' ? 'ja-JP' : 'en-US'}`, item: 'li.b_algo h2 a', map: mapLinks },
  google: { name: 'google', url: (q) => `https://google.com/search?q=${encodeURIComponent(q)}&hl=${OPTS.locale === 'ja-JP' ? 'ja' : 'en'}&num=${Math.min(OPTS.count, 10)}`, item: 'div.g a:has(h3)', map: mapLinks },
};

async function doSearch(page, query) {
  const primary = ENGINES[OPTS.engine] || ENGINES[OPTS.engine === 'duckduckgo' ? 'duckduckgo' : OPTS.engine];
  const chain = [primary, ENGINES['duckduckgo-lite'], ENGINES['duckduckgo-sqe']].filter(Boolean);
  const errors = [];
  for (const eng of chain) {
    const t0 = Date.now();
    try {
      const resp = await page.goto(eng.url(query), { waitUntil: OPTS.waitUntil, timeout: OPTS.timeout });
      await page.waitForSelector(eng.item, { timeout: Math.min(OPTS.timeout, 15000) }).catch(() => {});
      const raw = await page.$$eval(eng.item, (els, n) => els.slice(0, n).map(function (e) { return { title: (e.textContent || '').trim().replace(/\s+/g, ' '), url: e.href }; }), OPTS.count);
      const seen = new Set();
      const results = raw.map((r) => ({ title: r.title, url: cleanUrl(r.url) })).filter((r) => { if (!r.url || seen.has(r.url)) return false; seen.add(r.url); return true; });
      if (results.length) return { engine: eng.name, query, count: results.length, results, ms: Date.now() - t0 };
      errors.push(`${eng.name}:0-results(status=${resp ? resp.status() : 'null'})`);
    } catch (e) {
      errors.push(`${eng.name}:${firstLine(e.message)}`);
    }
  }
  throw new Error(errors.join(' | ') || 'no engine produced results');
}

// ---------- fetch / DOM解析 / JS実行 ----------
async function readPage(page) {
  const res = { url: page.url(), status: null, title: await page.title().catch(() => null) };
  if (OPTS.evalExpr) {
    res.eval = await page.evaluate(`(() => { try { const r = (${OPTS.evalExpr}); return r === undefined ? null : r; } catch (e) { return 'ERROR: ' + e.message; } })()`);
  }
  if (OPTS.selector) {
    res.selector = OPTS.selector;
    res.items = await page.$$eval(OPTS.selector, (els, attr, n) => els.slice(0, n).map(function (e) {
      if (attr === 'html') return e.outerHTML;
      if (attr === 'href') return e.href || e.getAttribute('href');
      return (e.textContent || '').trim();
    }), OPTS.attr, OPTS.count);
  } else {
    const m = OPTS.extract;
    if (m === 'html') res.content = await page.content();
    else if (m === 'text') res.content = (await page.evaluate(() => document.body ? document.body.innerText : '')).slice(0, 200000);
    else if (m === 'links' || m === 'links-meta') {
      res.links = await page.$$eval('a[href]', (els) => Array.from(new Set(Array.from(els, (e) => e.href))).filter((u) => u && !u.startsWith('javascript')).slice(0, 2000).map((u) => ({ url: u })));
      if (m === 'links-meta') res.meta = await page.evaluate(() => Array.from(document.querySelectorAll('meta[name],meta[property]')).map((x) => ({ key: x.getAttribute('name') || x.getAttribute('property'), content: x.getAttribute('content') })));
    } else if (m === 'meta') res.meta = await page.evaluate(() => Array.from(document.querySelectorAll('meta[name],meta[property]')).map((x) => ({ key: x.getAttribute('name') || x.getAttribute('property'), content: x.getAttribute('content') })));
  }
  return res;
}

async function doFetch(page, target) {
  const t0 = Date.now();
  const resp = await page.goto(target, { waitUntil: OPTS.waitUntil, timeout: OPTS.timeout });
  if (OPTS.waitFor) await page.waitForSelector(OPTS.waitFor, { timeout: OPTS.timeout }).catch(() => {});
  if (OPTS.waitMs) await page.waitForTimeout(OPTS.waitMs);
  const res = await readPage(page);
  return { mode: OPTS.mode, status: resp ? resp.status() : null, ...res, ms: Date.now() - t0 };
}

async function doFollow(page, startUrl) {
  const hops = [];
  let url = startUrl;
  for (let i = 0; i < OPTS.hops; i++) {
    const t0 = Date.now();
    try {
      const resp = await page.goto(url, { waitUntil: OPTS.waitUntil, timeout: OPTS.timeout });
      if (OPTS.waitFor) await page.waitForSelector(OPTS.waitFor, { timeout: OPTS.timeout }).catch(() => {});
      if (OPTS.waitMs) await page.waitForTimeout(OPTS.waitMs);
      const hop = await readPage(page);
      hop.status = resp ? resp.status() : null;
      hop.ms = Date.now() - t0;
      hops.push({ hop: i + 1, ...hop });
      const next = await page.$$eval(OPTS.selector || 'a[href]', (els) => (els[0] ? els[0].href : null)).catch(() => null);
      if (!next || next === hop.url) break;
      url = next;
    } catch (e) {
      hops.push({ hop: i + 1, url, error: firstLine(e.message), ms: Date.now() - t0 });
      break;
    }
  }
  return { mode: 'follow', hops };
}

// ---------- main ----------
parseArgs(process.argv.slice(2));

try {
  const page = await newPage();
  const result = OPTS.mode === 'search' ? await doSearch(page, OPTS.query)
    : OPTS.mode === 'follow' ? await doFollow(page, OPTS.target)
      : await doFetch(page, OPTS.target);
  out({ status: 'ok', accessed: new Date().toISOString(), ...result });
} catch (e) {
  out({ status: 'error', mode: OPTS.mode, error: firstLine(e.message) });
  process.exitCode = 1;
} finally {
  try { if (ctx) await ctx.close(); } catch { }
  if (tmpDir && !OPTS.keepTmp) { try { fs.rmSync(tmpDir, { recursive: true, force: true }); } catch { } }
}
