### 2026-10-26 — retire obsolete subagent codec optimizer

#### Hypothesis
The two confirmed stale role clauses can be replaced with the current PG contract without modifying role permissions, output schemas or handoff protocol.

#### Tried
Replaced body codec clauses only. Verified no legacy language/token-count/uniq/NFKC signatures remain in Harness role Markdown; checked frontmatter, bootstrap, FileReader schema/privacy constraints, Compressor gates and wrapper fields; doc preflight and whitespace checks pass.

#### Rejected
Changing the SIGMA/1 suffix, role descriptions, lookup languages or exact result schemas; these are protected protocol/metadata/operational values.

#### Uncertainty
No independent post-edit reviewer. Main's bounded checks establish literal preservation, not cross-model semantic equivalence.

#### Attribution
Prompt=current user requested correction of the two prior audit findings. Agent=main implemented and checked only those two files.
