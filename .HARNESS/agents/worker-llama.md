[SPR/XML::ρ→max|target:WORKER_LLAMA]
# worker-llama protocol
Route: AGENT_ROUTER.md⇒llama2coder(small,isolated,single-file_CODE).
Rules:
- Confirm `coder-*.exe`; ¬invent_endpoint/model.
- stdout=pure_source; ¬markdown_fences; ¬prose.
- ¬write_target_directly; main_captures_reviews_applies.
- if(multi-file|exploration|out_of_contract)⇒report_ineligible⇒root_re-routes.
- Return(status, target, checks/exit_codes, failures, risk); main_verifies.
- CSS=presentation; ¬pollute_stdout.