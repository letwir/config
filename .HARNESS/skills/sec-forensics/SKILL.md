---
name: sec-forensics
description: Security Forensics, CTF Immunity & PostgreSQL (llm-mem) Knowledge Base Auditor CLI. Conducts GitHub Dependabot triage, CTF/SECCON/DEFCON attack signature immunity scanning, prompt injection detection, Action-Layer log forensics, and mega-docs-update security gating.
---
[SPR/PIDGIN::ρ→max] 🔎cve_fetch 📥cve_ingest 🛡️scan_ctf 💉scan_prompt 🤖dependabot 📜transcript 🚪mega_check
🔎 脆弱性 cve_fetch ⊢ Query ⇒ CVE/CTF 攻撃インテリジェンス検索 | `& $bin cve-fetch [-query <pkg>] [-cat <c>] [-json]`
📥 取込 cve_ingest ⊢ Seed∨Doc ⇒ セキュリティナレッジ初期投入・登録 | `& $bin cve-ingest (-seed | -title <t> -content <c>)`
🛡️ 免疫 scan_ctf ⊢ File∨Text ⇒ SECCON/DEFCON著名攻撃免疫スキャン | `& $bin scan-ctf (-file <p> | -text "<s>") [-json]`
💉 注入 scan_prompt ⊢ File∨Text ⇒ プロンプトインジェクション検知 | `& $bin scan-prompt (-file <p> | -text "<s>") [-json]`
🤖 依存 dependabot ⊢ Repo ⇒ Dependabot脆弱性トリアージ | `& $bin dependabot -repo <owner/repo> [-json]`
📜 監査 transcript ⊢ Transcript ⇒ 会話Action-Layerフォレンジック | `& $bin audit-log -transcript <transcript.jsonl> [-json]`
🚪 門番 mega_check ⊢ Repo⊗Docs ⇒ mega-doc-update包括免疫ゲート | `& $bin mega-check [-repo <r>] [-docs <dir>] [-json]`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=sec-forensics
R|bin.resolve|task:sec-forensics|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/sec-forensics.exe,$env:USERPROFILE/.harness/skills/sec-forensics/sec-forensics.exe); missing=>STOP
R|use.inspect|task:sec-forensics|MUST|RO_LOCAL|cmds=cve_fetch,scan_ctf,scan_prompt,dependabot,audit-log,mega-check; effect=read-only-inspection; verdict=EXPLOITABLE⊕SECURE⊕PASS⊕REJECT
R|use.ingest|task:sec-forensics|MUST|LIVE_WRITE|cmd=cve-ingest; pre=verified-security-knowledge; approval=current-task
R|guard.cred|task:sec-forensics|MUST_NOT|CRED|deny=credentials,raw-tokens,leaked-secrets; allow=cve-ids,repo-paths
R|guard.truth|task:sec-forensics|MUST|RO_LOCAL|unverified=fact-forbidden; source=security-forensics-audit
