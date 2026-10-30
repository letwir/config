from __future__ import annotations

import json
import os
import pathlib
import re
import shutil
import subprocess

ROOT = pathlib.Path(r"C:\Users\letwir\.harness")
RESULT = ROOT / "memo" / "20261025-agy-skill-evaluation" / "agy-result.json"
MODEL_ID = "gemini-4-argon-high"

def redact(text: str) -> str:
    patterns = [
        (re.compile(r"(?i)(password|token|secret|api[_-]?key|credential|authorization)(\s*[:=]\s*)\S+"), r"\1\2[REDACTED]"),
        (re.compile(r"(?i)(postgres(?:ql)?://)[^\s]+"), r"\1[REDACTED]"),
    ]
    for pattern, replacement in patterns:
        text = pattern.sub(replacement, text)
    return text

home = pathlib.Path.home()
userprofile = pathlib.Path(os.environ.get("USERPROFILE", str(home)))
candidates = [shutil.which("agy"), shutil.which("agy.exe"), str(home / "bin" / "agy"), str(userprofile / "bin" / "agy.exe")]
exe = next((candidate for candidate in candidates if candidate and pathlib.Path(candidate).is_file()), None)
if not exe:
    raise SystemExit("AGY executable unresolved; stop")

inventory = subprocess.run([exe, "models"], cwd=ROOT, capture_output=True, timeout=45)
stdout = inventory.stdout.decode("utf-8", "replace")
rows = []
for line in stdout.splitlines():
    fields = line.split("\t", 1)
    if len(fields) == 2 and fields[0].strip() and fields[1].strip():
        rows.append({"id": fields[0].strip(), "display_name": fields[1].strip()})
if inventory.returncode != 0 or not any(row["id"] == MODEL_ID for row in rows):
    RESULT.write_text(json.dumps({"inventory_exit": inventory.returncode, "selected_model_found": False, "rows": rows}, ensure_ascii=False, indent=2), encoding="utf-8")
    raise SystemExit("Argon ID not confirmed in fresh inventory; stop")

prompt = """<Γ>
SIGMA/1: PIDGEN/text
Role: Independent evaluator / read-only skill review
Target: Assess the quality and operational correctness of `skills/agy-subagent-router/SKILL.md` for future bounded AGY CODE delegation.
Acceptance: Give concise evidence-backed findings: (1) overall verdict PASS|PASS_WITH_ADVISORIES|REJECT, (2) what is clear/useful, (3) concrete defects or ambiguities with exact section/line evidence and smallest corrections, (4) whether command examples and model inventory gate agree with the referenced current CLI facts in this task, (5) remaining uncertainties. Distinguish source-backed facts from inference. Do not edit, execute examples, run models, delegate, or claim runtime verification.
Scope: Read only `skills/agy-subagent-router/SKILL.md`, `rules/agy.lrf`, `rules/documentation.lrf`, `agents/AGENT_ROUTER.md`, `agents/worker-agy.md`, `agents/SUBAGENTS.lrf`, and `HARDWARE.md` under the mounted Harness workspace. Read-only critique only; no file writes, external/network/DB/VCS effects, credential access, or child delegation.
Known facts: This is a one-off evaluation explicitly requested by the current user, outside the skill normal bounded-CODE trigger; it grants only this read-only review task, not implementation permission. Current installed CLI help showed effort options low|medium|high|max and modes plan|accept-edits. A fresh read-only models call in this retry returned the current inventory; selected exact ID is `gemini-4-argon-high`. Do not infer model facts beyond this. The skill and its preservation record passed parent-run target checks; full project Polyglyph lint is separately blocked by `manifest/inventory mismatch`.
SIGMA/1
</Γ>"""

try:
    result = subprocess.run(
        [exe, "--print", prompt, "--model", MODEL_ID, "--effort", "high", "--mode", "plan", "--sandbox", "--add-dir", str(ROOT), "--output-format", "json"],
        cwd=ROOT,
        capture_output=True,
        timeout=600,
    )
    stdout_text = result.stdout.decode("utf-8", "replace")
    stderr_text = redact(result.stderr.decode("utf-8", "replace"))
    payload = {
        "executable": exe,
        "inventory_exit": inventory.returncode,
        "selected_model": MODEL_ID,
        "inventory_rows": rows,
        "agy_exit": result.returncode,
        "timed_out": False,
        "stdout": redact(stdout_text),
        "stderr": stderr_text,
    }
except subprocess.TimeoutExpired as exc:
    payload = {
        "executable": exe,
        "inventory_exit": inventory.returncode,
        "selected_model": MODEL_ID,
        "inventory_rows": rows,
        "agy_exit": None,
        "timed_out": True,
        "stdout": redact((exc.stdout or b"").decode("utf-8", "replace") if isinstance(exc.stdout, bytes) else str(exc.stdout or "")),
        "stderr": redact((exc.stderr or b"").decode("utf-8", "replace") if isinstance(exc.stderr, bytes) else str(exc.stderr or "")),
    }
RESULT.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
print(f"INVENTORY_EXIT={inventory.returncode}; MODEL_CONFIRMED={MODEL_ID}; AGY_EXIT={payload['agy_exit']}; TIMED_OUT={payload['timed_out']}; RESULT={RESULT}")
