# Walkthrough: highest difficulty Astra fallback

1. An operational agy failure enters the serial GPT fallback route after process exit and partial-diff inspection.
2. The router classifies task difficulty and consults current model availability plus evaluation evidence.
3. Easy selects Luna; medium selects Luna or Sol; hard may select Sol or Astra.
4. Astra may be invoked as a subagent by an eligible route, but it is always a leaf and cannot invoke any subagent itself; this fallback route selects it only at highest difficulty.
5. Existing permission and policy stop gates still terminate without fallback.
