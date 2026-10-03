# Knowledge: symbolic subagent instruction preference

- User-provided target codec: minimize token count subject to semantic/pragmatic fidelity and ambiguity approaching zero; prefer mathematical symbols for relations, preserve intent/modality/time/polarity/register, and append the exact `SIGMA/1` contract marker.
- Existing subagent handoff contracts require the visible fields `Role`, `Target`, `Acceptance`, `Scope`, `Known facts` in that order, inside the `<Γ>` wrapper.
- Open interpretation: the requested final-line marker may refer to the compressed task payload; the existing wrapper still requires `</Γ>` as the overall closing line. No persistent policy files were changed in this interaction.
