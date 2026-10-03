# Level 2 effect mismatch repair

### 2026-09-24

#### Hypothesis

The diagnostic comes from mutation language placed in a read-only `RO_LOCAL` rule.

#### Tried

Inspected the current SIGMA, documentation LRF, and harness-lint diagnostic. Split `doc.secret-found` into a read-only STOP/report rule and a `VCS_WRITE` rule requiring current-task authorization for history rewrite or tracked-file removal.

#### Rejected

Changing the entire rule to `LW_SCOPE` would overstate the effect of the ordinary stop/report action.

#### Uncertainty

No remaining Level 1 or Level 2 diagnostics were observed for this file.

#### Attribution

PromptDefect: none identified. AgentDefect: prior rule combined two effect classes; corrected in this task.

#### Search

Checked the exact lint diagnostic and all `doc.secret-found` rule text in the target file.

#### Correction

Each action now sits under its matching effect classification.

#### Emotion

Neutral; correction is localized.

#### Thoughts

Higher-level validation confirms the effect mismatch is resolved. Existing unrelated working-tree changes were not touched.
