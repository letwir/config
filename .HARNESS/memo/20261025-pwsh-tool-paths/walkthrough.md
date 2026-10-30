# walkthrough: 20261025-pwsh-tool-paths

Request: Locate six unresolved profile tools under C:\PATH and Program Files, and consider profile-local PATH additions.

Facts from bounded executable-name discovery (hidden/no-ignore/follow, case-insensitive globs):
- C:\PATH exists and is already on this process PATH. Contains dd.exe, flac.exe, metaflac.exe and several tool directories; none of the six target executables were found there.
- Standalone 7-Zip: C:\Program Files\7-Zip\7z.exe. Its directory is not on this process PATH. AMD copies also exist; do not prefer vendor-private copies.
- Sakura: C:\Program Files (x86)\sakura\sakura.exe. Its directory is not on this process PATH.
- FFmpeg and FFprobe found only at C:\Program Files (x86)\らじれこ\libs\ffmpeg.exe and ffprobe.exe in searched roots. These are another app's bundled binaries; suitability for general use is unverified.
- No jtrim.exe or nkf.exe matches in the three searched roots. This does not prove absence elsewhere. Some inaccessible paths may have been skipped because discovery stderr was suppressed.

Recommendation, not implemented: Add only standalone 7-Zip and Sakura directories to process PATH in the profile, with deduplication; do not add all of Program Files or vendor-private app libraries. Keep unrelated files and machine/user persistent PATH unchanged. Main-only implementation exception remains unconfirmed under standing workflow constraints; request explicit permission before direct editing.

No product edits, installs, binary execution, VCS changes, or independent review. One inspection command had a PowerShell pipeline syntax error and was corrected. Earlier memory precedent search in the same task conversation returned parseable JSON but was not used as evidence. Memory ingest is covered by standing report directive; evaluation live write remains unapproved.
