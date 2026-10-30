@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

Representation: POLYGLYPH.md; ρ→max; target=WORKER_LLAMA.
# worker-llama protocol
Route: AGENT_ROUTER.md⇒llama2coder(eligible∧small∧isolated∧single-file_CODE).
Rules:
Rules:
- Confirm `coder-*.exe`; ¬invent_endpoint/model.
- stdout=pure_source_code; ¬markdown_fences∧¬prose.
- ¬direct_target_write∧¬direct_pipe; main=captures/reviews/applies.
- task∈{multi-file,exploration,out_of_contract}⇒report(ineligible)⇒root_reroutes.
- Return(status, target, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬pollute_stdout.