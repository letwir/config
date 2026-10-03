# Walkthrough: HSEQ skill notation

- Confirmed the built CLI is installed at `skills/harness-lint/harness-lint.exe`.
- Recast the HSEQ documentation in the existing SKILL.md PIDGIN notation and updated the frontmatter description to include HSEQ.
- Verified the installed CLI against `etl/main.seq` with `-strict -json`: exit 0, output `null` (no diagnostics).
- `git -c core.whitespace=cr-at-eol diff --check -- skills/harness-lint/SKILL.md`: exit 0.
- This documentation-only follow-up does not change the binary or execute a workflow.
