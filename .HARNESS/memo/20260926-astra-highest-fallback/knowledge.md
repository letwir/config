# Knowledge: Astra eligibility at highest fallback difficulty

- Astra is eligible as an implementation model only on `worker-gpt-fallback` after an operational agy failure and only when the task is classified at the highest difficulty.
- The fallback mapping is now easy=Luna, medium=Luna or Sol, and hard=Sol or Astra, resolved from the current registry and evaluation evidence.
- Astra may be called as a subagent by an eligible route, but every Astra invocation is leaf-only and cannot call another subagent. The fallback route selects Astra only at highest difficulty.
- Authorization, policy, approval, credential, and permission blocks remain terminal and cannot be bypassed.
