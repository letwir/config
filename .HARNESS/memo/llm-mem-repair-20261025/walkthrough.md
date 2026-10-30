# Walkthrough — llm-mem-repair-20261025

1. Locate likely `llm-memory` source directory by filenames.
2. Check for `POLYGLYPH.md` in the target directory and its ancestors.
3. When no codec is found, stop source-dependent work and ask for the authorized codec path or an explicit resolution.
4. Do not edit/build the CLI or perform live database writes while blocked.
