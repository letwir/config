### 2026-10-26 — subagent instruction codec audit

#### Hypothesis
The prior bounded migration may have updated central rules/examples while leaving direct role-specific codec instructions stale.

#### Tried
Inventoried 15 Harness role documents and 3 Pi profiles; searched explicit PG markers and legacy language/tokenizer/NFKC formula signatures. Confirmed stale rule bodies in Compressor and FileReader; checked router/worker/orchestrator and Pi profile pointers.

#### Rejected
Treating all frontmatter descriptions as failures: current contract says preserve frontmatter exactly, so schema metadata is an explicit preservation boundary.

#### Uncertainty
Search proves the two exact conflicts, not semantic equivalence of every role prompt. No edits were authorized by this request.

#### Attribution
Prompt=read-only compliance check. Agent=main local search; no delegated reviewer and no independent semantic proof.
