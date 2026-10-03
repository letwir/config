# HSEQ sequence lint

HSEQ (`harness-seq/1`) is the named dialect for `.seq` workflow files. Put `@dialect harness-seq/1` on line 1, retain `@seq 1` or `@seq 2` as the file-format version, and set `@entry` to a declared section label. `harness-lint -path <file-or-directory>` scans `.seq` alongside `.lrf`; current structural checks cover the dialect marker, version directive, entry directive, duplicate section labels, and entry-to-section match.

The source implementation is in `A:\Users\letwir\repo\mcp2go\src\harness-lint` (`main.go`, `seq_lint.go`); the `.harness` skill contains the installed CLI binary. The lint is static: it does not parse every expression or execute a workflow. A clean lint result does not establish runtime behavior.
