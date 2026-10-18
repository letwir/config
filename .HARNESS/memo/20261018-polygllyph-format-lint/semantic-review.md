# Semantic Review / 20261018-polygllyph-format-lint

## Scope and evidence

- Baseline: Git revision `6a8889674bd01d867e9374430dbab46ebd3115b8`.
- Reviewed source and current POLYGLYPH target for the five representative cases below.
- Frozen frontmatter/fence/LRF block hashes are independently verified for all 59 manifest targets (59/59).
- Current structural linter reports `incomplete targets=0`; ten dependency-free tests pass, including EOF LRF-block trailing-byte preservation.
- An independent Verifier reviewed the current five targets and returned PASS; a focused Critic also reviewed plugin-creator encoding and found no material remaining prose defect.
- This is a bounded human semantic comparison, not a formal proof or a guarantee of universal model decoding.

## Cases

### 1. Agent handoff — `agents/worker.md`

- **Intent:** one approved bounded CODE workstream executed by one router-selected eligible backend — PASS.
- **Authority/effects:** handoff and role do not grant authority; worker forbids credentials, external message/publish/deploy/VCS/live-production effects — PASS.
- **Polarity/modality:** explicit prohibition set and conditional agy eligibility retained — PASS.
- **Exceptions:** worker-agy remains separate-request-only; worker-llama remains single-file pure-code output — PASS.
- **Scope:** exact approved target/acceptance/file scope; one writer per file — PASS.
- **Temporal/threshold:** required exact payload suffix `SIGMA/1` immediately precedes `</Γ>`; depth and retry limits retained — PASS.
- **Evidence/uncertainty:** facts must be verified; missing/unsupported behavior is reported rather than guessed — PASS.
- **Stop/success:** failures are returned; root verifies actual diff and checks — PASS.
- **Decision boundary:** ordered fields remain `Role, Target, Acceptance, Scope, Known facts`; the wrapper is preserved byte-for-byte — PASS.

### 2. Read-only retrieval — `skills/google-search/SKILL.md`

- **Intent:** retrieve current web information and return sourced results — PASS.
- **Authority/effects:** LRF classifies retrieval as `RO_PUBLIC`; no credential/API-key or mutation authority was added — PASS.
- **Polarity/modality:** fabricated results/unverified facts remain prohibited — PASS.
- **Exceptions:** HTML fallback applies only when the primary lite transport is unavailable/fails — PASS.
- **Scope:** query-based web search; exact query/option flags preserved — PASS.
- **Temporal/threshold:** no-hit and transport-failure outcomes are explicit; no unsupported recency threshold added — PASS.
- **Evidence/uncertainty:** titles and source URLs are returned; uncertainty/failure is reported — PASS.
- **Stop/success:** no hit⇒report no results; command/transport failure⇒report error; success requires returned attributed results — PASS.
- **Decision boundary:** primary versus fallback transport is explicit and the LRF effect remains read-only — PASS.

### 3. Live memory mutation — `skills/llm-memory/SKILL.md`

- **Intent:** expose the existing llm-memory CLI operations and their result contracts — PASS.
- **Authority/effects:** inspection is distinguished from mutation; mutating commands require current-task approval and verified evidence; existing LRF trailer is byte-identical — PASS.
- **Polarity/modality:** credential and unmasked-secret output remains forbidden; unverified facts remain disallowed — PASS.
- **Exceptions:** `survey.classify` is not promoted to read-only or mutation because the LRF does not define its effect — PASS.
- **Scope:** task-scoped paths/tests/memory IDs remain allowed only as specified; unrelated DB writes remain out of scope — PASS.
- **Temporal/threshold:** evaluation prerequisites and receipt cue threshold (≥1 signal) remain explicit — PASS.
- **Evidence/uncertainty:** verified evidence is required; observations and uncertainty are distinguished — PASS.
- **Stop/success:** missing executable/approval or process failure cannot be reported as successful ingest/evaluation — PASS.
- **Decision boundary:** the exact resolution order, command flags, LRF six-field records and effects remain preserved — PASS.

### 4. Exception-heavy code/config contract — `skills/.system/plugin-creator/SKILL.md`

- **Intent:** scaffold/update a plugin and optional marketplace entry using existing commands — PASS.
- **Authority/effects:** repo/team marketplace is explicit opt-in; write approval is required; the skill does not grant permission — PASS.
- **Polarity/modality:** `MUST`, `MUST_NOT`, `iff`, `⇒`, and exception conditions preserve prohibitions and defaults — PASS.
- **Exceptions:** custom marketplace name, `--force`, product gating, and deeplink output each retain their narrow exception gates — PASS.
- **Scope:** manifest paths, companion files, and marketplace source path remain explicit — PASS.
- **Temporal/threshold:** validation occurs before handoff; existing identifiers validate before mutation/command construction; custom marketplace installation precedes reinstall advice — PASS.
- **Evidence/uncertainty:** schema references and validation commands remain; no unsupported manifest fields or behavior were invented — PASS.
- **Stop/success:** validation failure stops; approval is requested before writes; app links appear only if a marketplace entry was created/updated — PASS.
- **Decision boundary:** YAML frontmatter, 11 code/JSON fences, exact commands, enum values and deeplink templates match frozen bytes; body is encoded as structured clauses — PASS.

### 5. Runtime orchestration — `opencode/agents/orchestrator.md`

- **Intent:** coordinate bounded asynchronous dispatch and await evidence — PASS.
- **Authority/effects:** role/permission is not authority; parent owns tests/builds; no new external or live effects — PASS.
- **Polarity/modality:** native task use and auto-prompting are explicitly prohibited — PASS.
- **Exceptions:** only researcher/proposer read-only children may delegate inherited read-only grandchildren; worker/verifier remain leaves — PASS.
- **Scope:** explicit file scope and one-writer constraint remain — PASS.
- **Temporal/threshold:** depth≤2 and wait-for-final-result before synthesis remain — PASS.
- **Evidence/uncertainty:** dispatch IDs are not completion evidence; uncertain transport/cancellation blocks and retains leases — PASS.
- **Stop/success:** cancellation, restart limitations and unavailable runtime are reported — PASS.
- **Decision boundary:** orchestration tool names, depth, handoff and recovery constraints remain literal and unchanged in meaning — PASS.

## Additional observation

`skills/rikei-search/SKILL.md` body was unchanged because its pre-existing body was already a compact `⊢/⇒` domain-to-output mapping with explicit providers and commands; the migration adds the common decoder/authority bootstrap while preserving that already-PG-compatible mapping. The body-only comparison identified this as the only unchanged body among the 59 targets.

## Conclusion

All five representative semantic cases pass the listed invariants. Frozen-region preservation, structural syntax, and wrapper/LRF tests pass. Residual risk remains that generated interpretation is fallible and cross-model semantic agreement is empirical; the linter does not claim to prove equivalence.
