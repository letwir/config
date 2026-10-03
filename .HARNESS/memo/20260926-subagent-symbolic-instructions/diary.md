# Diary: symbolic subagent instruction format

- Request: use a compact instruction format for subagents, choosing mathematical symbols for unambiguous relations and preserving semantic/pragmatic meaning, with an exact `SIGMA/1` suffix contract.
- Inspection: compiled bootstrap instructions and reviewed the existing manual handoff contract and `agents/SUBAGENTS.lrf`. The current contract already defines a semantic-preserving symbolic codec, while handoff wrapper and field ordering remain mandatory.
- Outcome: acknowledged the requested format for subagent task payloads. No subagent was dispatched and no persistent policy or application files were changed.
- Uncertainty: the intended interaction between the payload's final `SIGMA/1` line and the enclosing wrapper's required final `</Γ>` line was not resolved.
