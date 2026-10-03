# Walkthrough: agy operational failure fallback

1. CODE begins with the router-selected direct backend.
2. If agy completes, normal verification continues without fallback.
3. If agy reports an operational failure, stop or confirm exit and inspect the partial workspace diff.
4. Do not fallback for authorization, policy, approval, credential, or permission blocks.
5. Otherwise select one current GPT-family implementation model by task difficulty and evaluation evidence: Luna for easy, Luna or Sol for medium, Sol for hard.
6. Continue serially within the same acceptance criteria and exact file scope; never overlap the failed agy process.
7. The main agent inspects the final diff and runs deterministic checks.

