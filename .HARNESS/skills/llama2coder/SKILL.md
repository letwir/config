---
name: llama2coder
description: Generate clean, stdout-only source code using a llama.cpp compatible OpenAI REST API.
---
[SPR/PIDGIN::ρ→max] 💻code 🌊stream
💻 コード code ⊢ Prompt⊗Language ⇒ stdout直接純粋コード生成(MDフェンス除去済) | `& $bin -p "<prompt>" -l "<lang>" [-u <url>] [-m <model>]`
🌊 逐次 stream ⊢ Prompt⊗Language ⇒ トークン直接標準出力ストリーミング | `& $bin -p "<prompt>" -l "<lang>" -s`

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=llama2coder
R|bin.resolve|task:llama2coder|MUST|RO_LOCAL|bin=first-existing((lsd.exe $env:USERPROFILE/.harness/skills/llama2coder/coder-*.exe)); missing=>STOP
R|use.generate|task:llama2coder|MUST|RO_LOCAL|cmds=code,stream; effect=read-only-code-generation; stdout=pure-code-without-markdown
R|use.redirect|task:llama2coder&flag:pipe|MUST|LW_SCOPE|write=pipe-to-target-file-only
R|guard.cred|task:llama2coder|MUST_NOT|CRED|deny=credentials,raw-keys; allow=prompts,code
R|guard.truth|task:llama2coder|MUST|RO_LOCAL|unverified=fact-forbidden; source=local-llama-server
