---
name: skill-creator
description: Create or update a Codex skill with appropriately scoped instructions and any needed supporting resources.
metadata:
  short-description: Create or update a skill
---

@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

# Skill Creator

Purpose: Codexにuseful∧non-obvious guidanceを与えるskills作成; unrelated workを制約しない.

## Core Principles

**Assume Codex already capable.** Include only decisionを変える∨workを改善する情報. Remove{generic advice, repeated instructions, speculative edge cases, examples lacking material task clarification}.

**Preserve user intent and scope.** A skill should support the requested task, not replace the user's chosen product, expand the assignment, modify unrelated configuration, or imply permission for additional external actions. Do not turn a particular example, past failure, or personal preference into a universal requirement.

Approval to complete a task does not expand its scope or execution permissions. For retrying or externally mutating workflows, define a stopping condition proportional to the risk.

**Specificity: riskに合わせる.** Multiple reasonable approaches⇒modelに適切な選択の余地. Detailed steps/deterministic scripts/absolute language only if required by correctness/safety/permissions/genuinely fragile workflow.

For open-ended work, describe the outcome and relevant decision criteria. For workflows with a preferred shape, offer useful examples or configurable scripts. Reserve fixed sequences and narrow parameters for operations where deviation would cause a concrete problem. Preserve non-obvious operational invariants, distinguish actual requirements from optional recommendations or local conventions, and avoid restating policies already enforced elsewhere.

**Keep discovery cheap and precise.** Skill names and descriptions are available before a skill is loaded. Describe the actual capability and when it applies, adding exclusions only when they prevent likely misrouting. Avoid exhaustive capability lists and catchalls that attract unrelated requests.

Keep skills self-contained; refer to another skill or tool only when the requested workflow genuinely requires it and it is available in the target environment. Specialized review, hardening, or audit workflows should apply when requested or genuinely needed, not merely because ordinary work touches the same subject.

**Progressive disclosure.** `SKILL.md`=shared purpose∧essential constraints∧useful routing. Substantial mode-specific guidance/schemas/examples/procedures→supporting references; 読込 only current-task-relevant references. Simple self-contained skill⇒router/extra files不要.

## Anatomy of a Skill

Skill=folder(required `SKILL.md`⊕optional resources needed by actual workflow):

```text
skill-name/
|-- SKILL.md                 Required skill instructions
|   |-- YAML frontmatter     Required name and description
|   `-- Markdown body        Instructions loaded when the skill is used
|-- agents/                  Optional UI metadata and invocation policy
|   `-- openai.yaml
|-- scripts/                 Optional executable helpers
|-- references/              Optional documentation loaded as needed
`-- assets/                  Optional files used in generated output
```

Choose the structure that fits the actual task. Some skills are short and self-contained; others route among operating modes or delegate complex mechanics to scripts. Avoid creating directories, placeholders, examples, or ancillary documentation without a clear use.

### SKILL.md

The YAML frontmatter identifies the skill and determines when it should be considered. Include the required `name` and `description`, and preserve supported optional fields such as existing `metadata` when appropriate.

The Markdown body is loaded only when the skill is used. Put the purpose, essential workflow, real constraints, and useful links there. Keep detailed procedures and examples in supporting references when they are relevant only to particular modes.

Skill情報→3段階disclosure:

1. **Name and description:** skill選択時にavailable⇒concise∧discriminating.
2. **SKILL.md body:** skill適用時にloaded⇒instructions=current taskにrelevant.
3. **Supporting resources:** 読込∨実行 only if current task actually needs them.

The entrypoint should be as short as the task permits while retaining important constraints. A large upper bound is not a target: move conditional detail into references when doing so improves clarity or context use, rather than waiting for the file to become unwieldy.

### Scripts

Same logicの反復再記述を避ける∨deterministic execution materially improves reliability⇒executable codeは`scripts/`へ.

- **Example:** `scripts/rotate_pdf.py` for a PDF operation that would otherwise require recreating the same code.
- **Useful for:** Repeated transformations, reliable API operations, data processing, and other concrete automation.
- **Validation:** Run new or changed scripts to verify their behavior. Scripts can usually be executed without loading their full implementation into context, although an agent may need to inspect them when patching or adapting them.

### References

Context限定documentation→`references/`.

- **Examples:** `references/schema.md` for database tables, `references/policies.md` for domain rules, `references/api_docs.md` for an API, or separate writing guides for different deliverables.
- **Useful for:** Schemas, API documentation, company policies, format-specific procedures, detailed workflows, and substantial examples.
- **Routing:** Link each reference from `SKILL.md` or another relevant resource and explain when it should be read. Keep information in one place instead of duplicating it across the entrypoint and references.

Keep references focused on maintained, task-specific information that changes the agent's decisions. Avoid copied manuals, exhaustive catalogs, and generic tutorials already available from authoritative sources. Before removing existing resources, inspect their callers and purpose.

For large references, include useful search terms or a short contents section when that makes the needed material easier to find.

### Assets

Generated outputに属するfiles(model instructionsではない)→`assets/`.

- **Examples:** `assets/logo.png`, `assets/slides.pptx`, `assets/font.ttf`, or `assets/frontend-template/`.
- **Useful for:** Templates, images, fonts, icons, boilerplate projects, and other files copied or adapted into the result.
- **Context:** Do not load assets as instructions unless the task requires inspecting them.

### UI Metadata and Invocation Policy

`agents/openai.yaml` can provide UI-facing metadata such as `display_name`, `short_description`, and `default_prompt`, along with invocation policy. When creating or updating those settings, read [references/openai_yaml.md](references/openai_yaml.md) and keep the values consistent with the skill.

Automatic skill selection is allowed by default. Change that default only when the user explicitly requests an explicit-only skill:

```yaml
policy:
  allow_implicit_invocation: false
```

This keeps the skill available when explicitly invoked as `$skill-name` without adding it to the model context automatically. Preserve unrelated existing UI, policy, and dependency fields when updating `agents/openai.yaml`.

The initializer creates this file automatically. For new or interface-only metadata, generate it with:

```bash
scripts/generate_openai_yaml.py <path/to/skill-folder> --interface key=value
```

The generator replaces the entire file. If an existing file contains `policy` or `dependencies`, update only the intended fields in place instead of regenerating it.

Include optional interface fields only when the user provides or requests them.

### What Not to Include

Include files that directly support the skill's work. Avoid adding a `README.md`, installation guide, changelog, duplicated quick reference, or other auxiliary documentation unless a specific task or packaging requirement calls for it.

## Progressive Disclosure in Practice

For a skill with multiple substantial modes, keep the shared guidance and mode-selection criteria in `SKILL.md`. Link each supporting reference where its use becomes relevant. Do not load every reference by default, duplicate reference content in the entrypoint, or add a routing layer when there is nothing meaningful to route.

For example, a deployment skill can keep provider selection in `SKILL.md` and separate provider details:

```text
cloud-deploy/
|-- SKILL.md
`-- references/
    |-- aws.md
    |-- gcp.md
    `-- azure.md
```

User chooses AWS⇒read `references/aws.md`; GCP/Azure guidesは読込禁止. Same pattern→business domains/deliverable types/other genuinely distinct modesの分離.

A short skill can instead route to details only when an advanced operation needs them:

```markdown
## Documents

Handle ordinary edits directly.

- For tracked changes, read [references/redlining.md](references/redlining.md).
- For document internals, read [references/ooxml.md](references/ooxml.md).
```

Examples=options≠required structure. Organization選択→skill利用を容易化∧irrelevant material読込回避.

## Create or Update a Skill

Adapt the work to the request. Creating a complex new skill may involve understanding realistic use cases, choosing supporting resources, initializing files, writing instructions, and validating the result. A narrow update to an existing skill may require only a focused edit and validation.

Ask clarifying questions only when the missing information matters and cannot be reasonably inferred. Respect a user-specified location; otherwise create discoverable skills in `$CODEX_HOME/skills`, or `~/.codex/skills` when `CODEX_HOME` is unset.

Keep automatic skill selection enabled unless the user explicitly requests an explicit-only skill. When the intended invocation mode is genuinely unclear and matters to the requested workflow, ask whether the user wants normal automatic discovery or explicit-only invocation; otherwise preserve the default. Do not infer explicit-only invocation from sensitive operations or required approvals: keep the skill discoverable and require authorization immediately before the actual mutation. Preserve an existing skill's invocation policy unless the user asks to change it.

For a new or substantially revised skill, consider the actual requests it should handle and which reusable resources would improve those tasks:

- A repeated PDF transformation may justify a `scripts/rotate_pdf.py` helper.
- An application-building workflow may benefit from an `assets/frontend-template/` starter.
- A data-analysis skill may need a `references/schema.md` guide to avoid rediscovering table relationships.

Create those resources only when their concrete benefit justifies them. If the user has already explained the task clearly, proceed without requesting additional examples.

### Naming

- Name文字=lowercase letters∪digits∪hyphens.
- Name長<64 characters; short action-oriented names優先.
- Tool/domain namespace if discovery改善.
- Skill folder名=skill名.

### Initialize a New Skill

For a new skill, use the bundled initializer when it helps create the required files consistently:

```bash
scripts/init_skill.py <skill-name> --path <output-directory> [--resources scripts,references,assets] [--examples]
```

For example:

```bash
scripts/init_skill.py my-skill --path "${CODEX_HOME:-$HOME/.codex}/skills"
scripts/init_skill.py my-skill --path "${CODEX_HOME:-$HOME/.codex}/skills" --resources references
```

Request only the resource directories the skill needs. Use `--examples` only when concrete placeholders would help, and replace or remove them before finishing. Do not initialize an existing skill again.

The initializer creates the skill directory, a concise `SKILL.md` starter, and `agents/openai.yaml`. It creates resource directories and example files only when requested. Pass generated UI values as `--interface key=value` when needed.

### Write the Instructions

The frontmatter `description` should briefly explain what the skill does and when it applies. Include a meaningful boundary when similar requests should not activate the skill.

For example:

```yaml
description: Create or edit Word documents when formatting, tracked changes, or comments require document-specific handling.
```

Put detailed workflows, tool choices, examples, and operating modes in the body or relevant references rather than listing them all in the description. Preserve supported optional frontmatter, such as existing `metadata`, when appropriate.

Write only the instructions needed for another Codex instance to perform the task well. State the desired outcome, non-obvious context, real constraints, and relevant references or tools. Preserve the user's explicit choices and existing authorization boundaries. Avoid prescribing a fixed structure, process, or number of steps when the task does not require one.

### Validate and Iterate

Validate the completed skill with:

```bash
scripts/quick_validate.py <path/to/skill-folder>
```

Validator検査={frontmatter, naming, unfinished scaffold placeholders}; skillのgood decisionsを証明しない. Also verify{descriptions remain discriminating, instructions preserve user intent, references discoverable, added scripts actually work}.

When testing is warranted, verify observable behavior or meaningful invariants. Avoid tests that merely match generated wording, headings, or regex patterns.

Improve the skill based on real usage or demonstrated failures. Prefer a narrow correction to accumulating universal rules for every observed example.

## Independent Forward-Testing

Use an independent subagent pass when a skill is sufficiently complex or risky that realistic behavioral validation would add meaningful confidence, and when delegation is available and authorized. Ordinary creation or small edits do not automatically require subagents.

Give the evaluating agent a realistic user request, the skill, and the minimum raw artifacts needed to perform the task. Do not provide the intended answer, suspected bug, proposed fix, or prior conclusions unless the evaluation genuinely requires them.

For example:

```text
Use $skill-name at /path/to/skill-name to complete this realistic request.
```

Keep the evaluation scoped to permitted resources and side effects. Use an isolated temporary workspace for generated artifacts so they do not enter the working tree or contaminate later evaluations. Ask for approval when the proposed evaluation would require additional authorization, affect a live production system, or impose substantial time or cost. Review the actual outcome and artifacts, then make only changes supported by the observed behavior.
