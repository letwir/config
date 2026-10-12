import fs from 'node:fs'; import os from 'node:os'; import path from 'node:path';
import { chromium } from 'playwright-core';
const exe='/data/data/com.termux/files/usr/lib/chromium/chrome';
const UA='Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/154.0.0.0 Safari/537.36';
const ARGS=['--no-sandbox','--disable-dev-shm-usage','--disable-gpu','--disable-software-rasterizer','--no-first-run','--no-default-browser-check','--disable-extensions','--hide-scrollbars','--disable-crash-reporter','--noerrdialogs','--disable-sync','--mute-audio'];
const q=encodeURIComponent('東京タワー 高さ');
const targets=[
 ['html','a.result__a',`https://html.duckduckgo.com/html/?q=${q}&kl=jp-jp`],
 ['lite','a.result-link',`https://lite.duckduckgo.com/lite/?q=${q}`],
 ['sqe','[data-testid="result-title-a"]',`https://duckduckgo.com/?q=${q}&ia=web`],
];
for (const [name,sel,url] of targets){
  const tmpDir=fs.mkdtempSync(path.join(os.tmpdir(),'bs-'));
  const ctx=await chromium.launchPersistentContext(tmpDir,{executablePath:exe,headless:true,timeout:30000,locale:'ja-JP',userAgent:UA,viewport:{width:1366,height:900},javaScriptEnabled:true,bypassCSP:false,args:ARGS});
  const page=ctx.pages()[0]||(await ctx.newPage());
  page.setDefaultTimeout(30000); page.setDefaultNavigationTimeout(30000);
  try{
    const r=await page.goto(url,{waitUntil:'domcontentloaded',timeout:30000});
    await page.waitForTimeout(2500);
    const n=await page.$$eval(sel, e=>e.length).catch(e=>'ERR:'+e.message.slice(0,40));
    const txt=await page.evaluate(()=>document.body?document.body.innerText.replace(/\s+/g,' ').slice(0,120):'NOBODY');
    console.log(name,'status=',r&&r.status(),'sel_count=',n,'final=',page.url().slice(0,80));
    console.log('   body:',txt);
  }catch(e){ console.log(name,'NAV-ERR',e.message.split('\n')[0].slice(0,100)); }
  await ctx.close(); fs.rmSync(tmpDir,{recursive:true,force:true});
  await new Promise(r=>setTimeout(r,2000));
}
