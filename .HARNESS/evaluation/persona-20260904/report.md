# Subagent persona evaluation

Evaluation cutoff: 2026-09-04 23:59:59 Asia/Tokyo. Accessed and evaluated locally: 2026-09-05.

## Sources and selection

| Source snapshot | Reference | Adopted idea |
| --- | --- | --- |
| harnessworks/harness-starter-kit, 62437bec, 2026-06-18 | https://github.com/harnessworks/harness-starter-kit/blob/62437bec264b2deed83353e8209660d645e86828/AGENTS.md | Adapt minimal pieces to the existing environment and record validation. |
| mastersof-ai/harness, 3f5c484, 2026-03-19 | https://github.com/mastersof-ai/harness/blob/3f5c4846a2f7ecdf48198d6d1116204132e967a5/README.md | Explicit role identity and concise evidence-oriented output. |
| madebywild/agent-harness, 2ecf44f, 2026-08-29 | https://github.com/madebywild/agent-harness/blob/2ecf44f97ab6ffd76f5c1a58ac57b930e2ccfa17/docs/architecture.md | Keep prompt, agent and configuration responsibilities distinct. |

Researcher found five candidates; the main agent independently opened the three adopted fixed snapshots and their commit pages. BA-CalderonMorales commit 8d3bc4b was excluded: 2026-09-04 18:08 UTC is 2026-09-05 JST. None of these sources establishes a causal benefit from dense emoji or mathematical symbols. No external prompts or runtime were installed.

## Design and constraints

Each role retains its selector and role string. A single archetype word follows one country flag and the combining symbol. Flags denote chosen language context, not exclusive national origin. English terms are used in US English; no claim of a uniquely American archetype is made.

| Role | Character |
| --- | --- |
| Explorer | 🇯🇵ギャル |
| Researcher | 🇯🇵博士 |
| Proposer | 🇺🇸Strategist |
| Auditor | 🇯🇵判事 |
| Worker | 🇯🇵職人 |
| Refactorer | 🇺🇸Surgeon |
| Verifier | 🇯🇵お嬢様 |
| Critic | 🇺🇸Skeptic |
| Blackhat | 🇯🇵怪盗 |

Emoji mark role imagery and concepts. ∧ combines traits, ⊕ layers related meanings, → orders presentation, ≠ distinguishes concepts, and ↔ compares a claim with evidence. Plain-language text carries the same meaning. These are presentation cues, not a new grammar, policy, execution plan or proof of improved reasoning.

Existing Explorer/Kobe/gyaru and Verifier/aristocratic identity are retained in localized wording. Explorer's previous epistemology line is folded into presentation language instead of directing experiments from CSS. Existing Blackhat voice change is retained in spirit. Authority, effects, routing, TOML, and role Markdown are untouched.

## Verification

Status: independent Verifier PASS for static compatibility and presentation/authority separation; integrated into active agents/. Baseline includes all pre-existing user edits to these nine files. Integration must first compare every active CSS hash with baseline, then copy only approved candidates, then verify exact candidate/active equality. No Git remote is configured; merge means local file integration, not a GitHub PR merge.

## Uncertainty and impact

This is a presentation improvement hypothesis. No before/after task-quality, latency or token-cost experiment has been performed; maximum effectiveness is not established. A review and illustrative persona samples are not a benchmark of runtime behavior across all nine roles.

Existing issue: SUBAGENTS.lrf evaluation-write path prompt-ref/eval/subagents_model_eval.json is absent; SIGMA reads evaluation/subagents_model_eval.json. This task records bounded evaluation locally without repairing unrelated routing.

PromptDefect: no task-blocking defect established; date cutoff interpreted in the configured Asia/Tokyo timezone. AgentDefect: researcher initially included a UTC date crossing the JST cutoff; excluded before adoption.
Research precedent checks: EPUB tool found, cache sync exceeded 30 seconds without a result; llm-mem search failed with missing search_document column (SQLSTATE 42703). No precedent result was inferred. Requested research model: gpt-5.6-luna / medium; effective provider model and reasoning not independently observable from the returned run evidence.

Commit dates above are UTC. Main-agent public GitHub API verification (curl exit 0): 62437bec = 2026-06-18T11:34:31Z; 3f5c484 = 2026-03-19T20:39:51Z; 2ecf44f = 2026-08-29T23:49:41Z. All precede the JST cutoff.

Completed checks: Verifier independently checked 9/9 UTF-8/no-BOM, one CSS rule, four declarations, original selector/role, one character word, absence of authority overrides, and baseline/config hashes. Main integration rechecked all baseline hashes before writing, copied only nine CSS files, confirmed active/candidate SHA-256 equality 9/9, and ran git diff --check (exit 0). Full fresh-session persona behavior across nine roles and A/B effectiveness remain unverified.
