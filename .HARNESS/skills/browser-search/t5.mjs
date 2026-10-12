import fs from 'node:fs'; import os from 'node:os'; import path from 'node:path';
import { chromium } from 'playwright-core';
const exe = '/data/data/com.termux/files/usr/lib/chromium/chrome';
const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(),'bs-'));
const ctx = await chromium.launchPersistentContext(tmpDir, {
  executablePath: exe, headless: true, timeout: 30000, locale: 'ja-JP',
  userAgent: 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36',
  viewport: { width: 1366, height: 900 }, javaScriptEnabled: true, bypassCSP: false,
  args: ['--no-sandbox','--disable-dev-shm-usage','--disable-gpu','--disable-software-rasterizer','--no-first-run','--no-default-browser-check','--disable-extensions','--hide-scrollbars','--disable-crash-reporter','--noerrdialogs','--disable-sync','--mute-audio'] });
const page = ctx.pages()[0] || (await ctx.newPage());
page.setDefaultTimeout(30000); page.setDefaultNavigationTimeout(30000);
const resp = await page.goto('https://html.duckduckgo.com/html/?q=' + encodeURIComponent('東京タワー 高さ') + '&kl=jp-jp', { waitUntil: 'domcontentloaded', timeout: 30000 });
console.log('status', resp && resp.status());
console.log('headers', JSON.stringify(resp && resp.headers()).slice(0,600));
const body = await page.evaluate(() => document.body ? document.body.innerText.slice(0,300) : 'NOBODY');
console.log('BODY:', body.replace(/\s+/g,' '));
await ctx.close(); fs.rmSync(tmpDir,{recursive:true,force:true});
