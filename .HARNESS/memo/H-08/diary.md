# H-08 Diary

### 2026-09-10

**Hypothesis:** A small local review can safely skip ordinary exploration and implementation stages when its effect, query, targets, commands, and terminal are all closed.

**Tried:** Added an early route after precedent lookup, a fixed-command collector, a hash-only receipt schema, structured route evaluation, and cross-host boundary tests.

**Rejected:** Sending local review into FINISH; treating the policy evaluator as proof that commands ran; accepting arbitrary commands, globs, traversal, duplicate targets, or raw output receipts; trusting a caller-supplied final route.

**Uncertainty:** Reparse-point tests depend on host support, so runtime containment checks remain authoritative. `rg` no-match proves only absence within the exact supplied targets.

**Attribution:** PromptDefect not found. AgentDefect found and corrected: the first implementation used a reserved PowerShell parameter name, ordered `rg` arguments incorrectly, persisted bounded raw output, trusted the proposed route, and briefly introduced line-ending whitespace failures.

**Search correction:** No relevant prior memory or public source was required; verified local contracts determined the design.

**Emotion and thoughts:** The useful shortcut became clear only after its terminal was made as strict as its entry predicate.
