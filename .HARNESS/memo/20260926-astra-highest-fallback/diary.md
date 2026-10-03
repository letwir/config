# Diary: allow Astra for highest difficulty fallback

- Request: allow Astra at the highest difficulty in the previously implemented agy-to-GPT fallback routing.
- Change: updated BOOTSTRAP, agy fallback, AGENT_ROUTER, SUBAGENTS, worker contracts, and the agy router skill so hard difficulty may select Sol or Astra.
- Boundary: Astra may be called as a subagent, but an Astra invocation cannot call another subagent. The fallback route selects Astra only at highest difficulty.
- Verification: targeted assertions, strict LRF checks, and `git diff --check` passed.
- Scope: no VCS write, publication, deployment, credential action, or external message was performed.
