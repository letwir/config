# Imported local skill source: human guide

The `mcp2go-src/` subdirectory in each mapped `.HARNESS/skills/<skill-name>/` folder preserves that package's staged source tree from sibling repository `../mcp2go/.harness/skill/src/`. Package contents and relative paths are unchanged. Nested `go.mod` files mark independent Go modules, notably Kavita `src/`, kokkou_api `src/`, llama2coder `src/`, and the hudousan_api package root. Existing skill files remain alongside (not merged into) this source subtree.

## Build and test

See [INSTALL.md](INSTALL.md) for the repeatable PowerShell validation command. Run it from the config repository root. It discovers Go modules only beneath `mcp2go-src/`, excludes the obsolete outer Kavita module, and includes the nested Kavita `src/` module. Caches and build outputs are temporary and outside the repository; the command does not install programs, configure credentials, or modify the sibling source repository.

Do not add generated executables, temporary caches, credentials, or machine-specific configuration. A passing Go check means packages compile and tests pass; it does not prove external APIs, services, databases, credentials, or optional runtime integrations are available. Python skills may require separately managed dependencies; follow their per-skill documentation.

## Source refresh boundary

An intentional refresh copies source/configuration/docs from the sibling repository's `src/` package into its matching `mcp2go-src/` subtree, preserving relative paths and module roots. Exclude caches, build outputs, and transient generated files. Review any proposed diff before replacing content, and do not overwrite existing skill files outside `mcp2go-src/`.
