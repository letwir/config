# Implementation Plan: POLYGLYPH instruction format and lint

## Purpose

Apply the user-provided `.HARNESS/POLYGLYPH.md` codec to repository LLM instruction/contract Markdown, and add a dependency-free, read-only structural linter with deterministic tests.

## Fixed source and scope

Baseline Git revision: `6a8889674bd01d867e9374430dbab46ebd3115b8` (all 59 selected paths are tracked at this revision).

Include exactly these 59 paths, discovered with POSIX-relative, case-sensitive, NFC-normalized UTF-8 path ordering and no symlink traversal:
- `.HARNESS/AGENTS.md` (1)
- `.HARNESS/agents/*.md` except `README.md` (16)
- `.HARNESS/skills/**/SKILL.md`, excluding every `node_modules` path (40)
- `.HARNESS/rules/README-forLLM.md` (1)
- `.HARNESS/opencode/agents/*.md` (1)

The machine-readable manifest stores the 59 paths relative to `HARNESS_ROOT` (the directory containing `POLYGLYPH.md`, i.e. `.HARNESS`); strip the leading `.HARNESS/` from the inventory above. The preservation fixture uses the same harness-relative keys. Baseline retrieval prefixes each key with `.HARNESS/` when invoking `git show` from the Git root. The manifest MUST enumerate the exact discovered paths. Any missing, duplicate, untracked, symlink, extra, or case/NFC-colliding path fails closed. New matching files not listed by the manifest are lint errors; no additional target is inferred.

Exclude human README/INSTALL/HARDWARE/state/reference/USAGE docs; `.toml`, JSON and YAML configs; source/scripts; standalone authoritative LRF; ETL sequence DSL; all `node_modules`; generated task artifacts; binaries. Keep existing unrelated workspace edits unchanged.

## Normative file serialization to add to `POLYGLYPH.md`

### Frontmatter and marker

- Inputs must be valid UTF-8 without BOM. LF and CRLF line terminators may coexist because the fixed source corpus already has mixed LF/CRLF files; do not normalize or reject that mixture. A bare CR, decode error, or BOM fails closed. Preserve each immutable region's original line-ending bytes exactly.
- YAML frontmatter exists iff the first physical line is exactly `---`. It ends at the first subsequent physical line exactly `---`. The complete byte range from byte 0 through the closing line terminator is immutable. Missing close delimiter fails closed; no YAML reparsing or normalization is performed.
- After optional frontmatter, zero or more empty physical lines are permitted; the first nonblank physical line MUST equal exactly `@meta=text/POLYGLYPH×SIGMA/1`. With no frontmatter, this marker is required on physical line 1. No other nonblank prefix is allowed.
- The next three physical lines MUST be exactly, in this order:
  `@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true`
  `@authority=SIGMA/1+routed-LRF-only;PG=representation-only`
  `@decode_failure=STOP;unknown=preserve;action=blocked;ask=true`
  These fixed bootstrap lines are the only material a consumer may act on before reading the codec. Resolve the codec by ascending parent directories from the file; if absent/unreadable, block the affected action and report the unknown. If multiple materially different readings remain for any Intent/Relation/Effect, Constraint, Exception, Scope, Authority, Decision/Threshold, Evidence, temporal, polarity, modality, or stop/success field, expose `PG_DECODE=AMBIGUOUS`, list the unresolved field(s), do not select a reading or take dependent action, and ask for resolution. A consumer may continue only with independent actions whose meaning is unambiguous.
- From `@meta` to EOF, content is PG except fenced opaque regions and an optional LRF/1 block. An LRF/1 block is immutable but does not necessarily end the PG body; if later body text exists, it resumes after the LRF block.

### Opaque fenced regions

Use CommonMark fence recognition deterministically: at most 3 leading spaces; opening run is >=3 identical backticks or tildes; a backtick fence's info string MUST NOT contain backticks; mixed-character fence runs are not delimiters. Closing run uses the same character, length >= opening, with only spaces/tabs following. Fences do not nest. The complete opening-through-closing bytes, including line endings, are opaque and MUST be preserved byte-for-byte and in order. Unterminated recognized fence fails closed. A marker inside a fence is not an LRF block.

### LRF/1 block

At most one standalone exact `#SIGMA LRF/1` line outside fences may occur. It begins an immutable LRF/1 block consisting of one exact three-field header `@lrf=1|aud=<nonempty>|scope=<nonempty>`, zero or more six-field `R`/`L` records, and blank separator lines. Following the header, every nonblank line beginning `R|` or `L|` MUST validate as a record; malformed rows fail closed. A second `@lrf=` header or second `#SIGMA LRF/1` marker before body resume fails closed. The block ends at the first nonblank line that is neither a header nor begins `R|`, `L|`, `@lrf=`, or `#SIGMA LRF/1`; that line resumes PG body, and fence parsing continues. Body lines beginning with these reserved prefixes must be escaped/rewritten unambiguously or the file fails closed. If no later prose exists, the block extends to EOF. Multiple markers fail. Each record is six pipe-delimited fields, nonempty ID/guard/effect/payload, modality ∈ `{MUST,MUST_NOT,MAY}`. Blank lines within/between header and record lines are permitted as part of the block; blank separator lines immediately preceding resumed PG prose are preserved as part of the immutable block. Malformed header/record/modal fails closed. The block bytes and order MUST match its source baseline exactly. Existing LRF `MAY` records remain untouched.

### Authority and wrappers

`POLYGLYPH.md` remains the decoder spec; read it before interpreting PG bodies. PG is representation only and never authorizes an effect. SIGMA/1 and routed LRF remain the only policy sources; conflicts are exposed and resolved, never silently overridden. Preserve runtime YAML frontmatter, every opaque code fence, every LRF/1 block, and exact handoff wrapper bytes, including ordered `Role`, `Target`, `Acceptance`, `Scope`, `Known facts`, exact suffix `SIGMA/1`, and following `</Γ>`.

## Deterministic baseline and linter

Before changing any target, Workstream A captures source regions from `git show 6a8889674bd01d867e9374430dbab46ebd3115b8:.HARNESS/<harness-relative-path>` and writes a frozen SHA-256 manifest at `.HARNESS/scripts/tests/fixtures/polyglyph-preservation.json`. For every target this records frontmatter (or absent), ordered opaque-fence hashes, and LRF/1 block hashes (or absent), keyed by path relative to `HARNESS_ROOT`. The capture implementation MUST independently recompute every saved signature from the same `git show` bytes and assert all 59 entries match before any source target is changed; a mismatch aborts and the fixture is rejected. The linter verifies current extracted regions against this frozen fixture; it does not use a moving HEAD or regenerate hashes during checks. Any region count/order/bytes change fails.

Add `.HARNESS/scripts/lint-polyglyph.py`, `.HARNESS/scripts/tests/test_lint_polyglyph.py`, `.HARNESS/scripts/tests/fixtures/polyglyph-targets.json`, `.HARNESS/scripts/tests/fixtures/polyglyph-preservation.json`, and standard-library-only tests. The linter is read-only and checks:
- exact eligible inventory vs manifest, existence/tracked regular-file/no-symlink/canonical-path guarantees;
- UTF-8/newline/frontmatter/PG-marker/bootstrap boundaries and fixed metadata lines;
- spec section uniqueness/order `§0` through `§8` and required conflict/uncertainty/authority invariants;
- complete/ordered fences, exact preservation signatures, and trailer recognition outside fences;
- exact LRF header/record arity, `R|L`, and `MUST|MUST_NOT|MAY` preservation;
- SIGMA serialized fixtures exactly `[SIGMA/1]\n<∑>\n<Contract>\n</∑>`, U+2211 not U+03A3, one matched pair, no extra header/footer, forbidden closing marker rejection;
- positive/negative `<Γ>` wrapper tests for closure, field order, and final `SIGMA/1` placement.

Lint verdict is structural only. It is never semantic proof or cross-model agreement. Before acceptance, write `.HARNESS/memo/20261018-polygllyph-format-lint/semantic-review.md` with representative source→PG comparisons read from the fixed baseline revision: (1) agent handoff, (2) a read-only skill, (3) a LIVE_WRITE skill, (4) a skill with code/config snippets and exceptions, (5) OpenCode agent. For each explicitly assess Intent, authority/effects, polarity/negation, modality incl MAY, exception, scope, temporal conditions, evidence/uncertainty, stop/success, and decision boundary. Mark PASS/FAIL with source/target paths and evidence; no credentials or irrelevant literal values. Any unresolved meaning fails acceptance; retain the source wording inside PG or stop that file, never guess.

## Disjoint worker workstreams and order

A (first): one context-selected GPT-family worker owns `.HARNESS/POLYGLYPH.md`, `.HARNESS/AGENTS.md`, `.HARNESS/rules/README-forLLM.md`, the manifest, frozen preservation fixture, linter, and tests. It MUST generate the baseline fixture before modifying any target. Validate the codec/linter scaffolding before any migration worker starts.

Then run independent workstreams B–E concurrently, with at most 4 active children (global policy maximum 6):
- B owns the exact 16 `.HARNESS/agents/*.md` files except README.md.
- C owns the first 20 paths in the manifest's NFC-normalized sorted `.HARNESS/skills/**/SKILL.md` list.
- D owns the remaining 20 skill paths.
- E owns `.HARNESS/opencode/agents/orchestrator.md`.
Each workstream uses exactly one currently available AGENT_ROUTER-selected worker, selected by comparable measured context consumption or the dated nonnormative fallback; no overlap, no substitute backend, no concurrent writes to the same file. User authorizes up to 20 children, but the applicable router caps parallelism at 6. Main owns final real diff and independent checks. No VCS write.

## Acceptance

PASS only when the exact 59-file inventory is migrated; baseline frontmatter/fences/LRF bytes pass; SIGMA and handoff fixtures include positive/negative cases and pass; semantic-review.md records all five cases with every listed invariant PASS; linter tests and LRF preflight pass; final diff is only the allowlisted implementation files plus `.HARNESS/memo/20261018-polygllyph-format-lint/{plan,knowledge,diary,walkthrough,evaluation,semantic-review,record,memory-receipts}.(md|json)`.

No deployment, commit, push, history rewrite, production dependency, live service write, or credential use. If a worker fails or any ambiguity/baseline check fails, stop the affected workstream and report it; never claim completion from partial lint.
