# Knowledge: agy failure to difficulty-selected GPT fallback

- The prior authoritative contracts limited fallback to agy pre-launch failure or CODE timeout and fixed continuation to GPT-6 Luna.
- The revised contract treats any operational agy failure, including launch, timeout, generation, backend, or malformed-result failure, as eligible for one serial GPT-family fallback.
- Authorization, policy, approval, credential, and permission blocks remain terminal and cannot be bypassed by fallback.
- GPT fallback selection is difficulty-aware and resolved against the current Codex registry plus evaluation evidence: easy uses Luna, medium uses Luna or Sol, and hard uses Sol. Astra remains main-agent-only.
- Before fallback, the agy process must be stopped or confirmed exited and any partial diff must be inspected. The fallback keeps the same bounded task and file scope and must not overlap agy.

