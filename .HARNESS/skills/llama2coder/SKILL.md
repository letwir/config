---
name: llama2coder
description: Generate clean, stdout-only source code using a llama.cpp compatible OpenAI REST API.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/PIDGIN::ρ→max|target:local-llama-codegen|legibility:LLM≫human] 💻code 🌊stream
🔎 code: source-code generation request ∧ stream: token-by-token request. ⊢ Resolve `$bin` to first existing configured `coder-*` executable; none ⇒ STOP. Mor(code): local llama.cpp-compatible OpenAI REST API client → stdout source code only, ¬Markdown fences. Args: `$bin -p "<prompt>" -l "<lang>" [-u <url>] [-m <model>]`
Mor(stream): local llama server → stdout token stream only. Args: `$bin -p "<prompt>" -l "<lang>" -s`
⛔ Other generation-failure behavior unspecified; preserve uncertainty. ✅ Success: requested output modality returned.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=llama2coder
R|bin.resolve|task:llama2coder|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/llama2coder/coder-*,~/.harness/skills/llama2coder/coder-*.exe,./coder-*,./coder-*.exe); missing=>STOP
R|use.generate|task:llama2coder|MUST|RO_LOCAL|cmds=code,stream; effect=read-only-code-generation; stdout=pure-code-without-markdown
R|use.redirect|task:llama2coder&flag:pipe|MUST|LW_SCOPE|write=pipe-to-target-file-only
R|guard.cred|task:llama2coder|MUST_NOT|CRED|deny=credentials,raw-keys; allow=prompts,code
R|guard.truth|task:llama2coder|MUST|RO_LOCAL|unverified=fact-forbidden; source=local-llama-server
