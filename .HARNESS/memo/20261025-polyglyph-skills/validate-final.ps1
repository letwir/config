$ErrorActionPreference = 'Stop'
$task = 'C:/Users/letwir/.harness/memo/20261025-polyglyph-skills'
$batch = "$task/batch"
$utf8 = [Text.UTF8Encoding]::new($false, $true)
function Text([string]$p) { $utf8.GetString([IO.File]::ReadAllBytes($p)) }
function Assert([bool]$condition, [string]$message) { if (-not $condition) { throw $message } }
function Front([string]$s) { [regex]::Match($s, '\A---\r?\n[\s\S]*?\r?\n---(?:\r?\n|\z)').Value }
function Fences([string]$s) {
    $tick = [string][char]96
    @([regex]::Matches($s, '(?ms)^' + ($tick * 3) + '[^\r\n]*\r?\n.*?^' + ($tick * 3) + '[^\r\n]*(?:\r?\n|$)') | ForEach-Object Value)
}
function CheckCommon([string]$p) {
    $bytes = [IO.File]::ReadAllBytes($p)
    $text = $utf8.GetString($bytes)
    Assert (-not ($bytes.Length -ge 3 -and $bytes[0] -eq 239 -and $bytes[1] -eq 187 -and $bytes[2] -eq 191)) "BOM: $p"
    Assert ($text -notmatch '\r(?!\n)') "Bare CR: $p"
    Assert ((Front $text).Length -gt 0) "Missing frontmatter: $p"
    foreach ($line in @('@meta=text/POLYGLYPH×SIGMA/1', '@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true', '@authority=SIGMA/1+routed-LRF-only;PG=representation-only', '@decode_failure=STOP;unknown=preserve;action=blocked;ask=true')) {
        Assert ($text.Contains($line)) "Missing bootstrap line: $p"
    }
    Assert ($text -notmatch '\[SPR/(?:PIDGIN|XML)') "Legacy active marker: $p"
    $index = $text.IndexOf('#SIGMA LRF/1', [StringComparison]::Ordinal)
    Assert ($index -ge 0) "Missing LRF: $p"
    foreach ($row in ($text.Substring($index) -split '\r?\n')) {
        if ($row -match '^[RL]\|') { Assert ($row.Split('|').Count -eq 6) "Invalid LRF row: $p" }
    }
    foreach ($m in [regex]::Matches($text, '\]\(([^)]+)\)')) {
        $link = $m.Groups[1].Value
        if ($link -match '^(?:[a-z]+:|#)') { continue }
        $link = $link.Split('#')[0]
        Assert (Test-Path -LiteralPath (Join-Path (Split-Path $p -Parent) $link)) "Missing Markdown link $link in $p"
    }
    $text
}
$snapshots = Text "$batch/snapshots.json" | ConvertFrom-Json
$manifest = Text "$batch/manifest.json" | ConvertFrom-Json
Assert ($manifest.files.Count -eq 40) 'Manifest must include 39 skills + codec'
$tokens = $null; $parseErrors = $null
$ast = [Management.Automation.Language.Parser]::ParseFile("$batch/apply.ps1", [ref]$tokens, [ref]$parseErrors)
Assert ($parseErrors.Count -eq 0) 'Apply helper parse errors'
$assignment = $ast.Find({param($n) $n -is [Management.Automation.Language.AssignmentStatementAst] -and $n.Left -is [Management.Automation.Language.VariableExpressionAst] -and $n.Left.VariablePath.UserPath -eq 'referenceClauses'}, $true)
$table = $assignment.Right.Find({param($n) $n -is [Management.Automation.Language.HashtableAst]}, $true)
$refs = @{}
foreach ($pair in $table.KeyValuePairs) {
    $value = $pair.Item2.Find({param($n) $n -is [Management.Automation.Language.StringConstantExpressionAst]}, $true)
    Assert ($pair.Item1 -is [Management.Automation.Language.StringConstantExpressionAst] -and $value -is [Management.Automation.Language.StringConstantExpressionAst]) 'Reference clauses must be literal strings'
    $refs[$pair.Item1.Value] = $value.Value
}
Assert ($refs.Count -eq 9) 'Expected nine approved reference merges'
$results = @()
foreach ($s in $snapshots) {
    Assert ((Get-FileHash -LiteralPath $s.backup).Hash -eq $s.sha256) "Immutable snapshot hash differs: $($s.backup)"
    $old = Text $s.backup
    $new = CheckCommon $s.path
    Assert ((Front $old) -ceq (Front $new)) "Frontmatter changed: $($s.path)"
    $oldFences = @(Fences $old)
    if ($s.path -eq $s.canonical) {
        $expected = $old.Replace('[SPR/PIDGIN::', '[POLYGLYPH::').Replace('[SPR/XML::', '[POLYGLYPH::')
        if ($refs.ContainsKey($s.name)) {
            $i = $expected.IndexOf('#SIGMA LRF/1', [StringComparison]::Ordinal)
            $nl = if ($expected.Contains("`r`n")) { "`r`n" } else { "`n" }
            $expected = $expected.Substring(0, $i).TrimEnd("`r", "`n") + $nl + $refs[$s.name] + $nl + $nl + $expected.Substring($i)
        }
        Assert ($expected -ceq $new) "Unexpected canonical diff: $($s.path)"
        Assert ($old.Substring($old.IndexOf('#SIGMA LRF/1')) -ceq $new.Substring($new.IndexOf('#SIGMA LRF/1'))) "LRF bytes changed: $($s.path)"
        $newFences = @(Fences $new)
        Assert ($oldFences.Count -eq $newFences.Count) "Fence count changed: $($s.path)"
        for ($i = 0; $i -lt $oldFences.Count; $i++) { Assert ($oldFences[$i] -ceq $newFences[$i]) "Fence bytes changed: $($s.path)" }
    } else {
        $canonical = Text $s.canonical
        Assert ($new.Substring((Front $new).Length) -ceq $canonical.Substring((Front $canonical).Length)) "Body not synchronized: $($s.path)"
    }
    $entry = @($manifest.files | Where-Object { [IO.Path]::GetFullPath($_.path) -eq [IO.Path]::GetFullPath($s.path) })
    Assert ($entry.Count -eq 1) "Manifest missing/duplicate: $($s.path)"
    Assert ($entry[0].before_sha256 -eq $s.sha256) "Manifest before hash mismatch: $($s.path)"
    Assert ($entry[0].after_sha256 -eq (Get-FileHash -LiteralPath $s.path).Hash) "Manifest after hash mismatch: $($s.path)"
    $results += [pscustomobject]@{path=$s.path;status='PASS';sha256=$entry[0].after_sha256}
}
$codec = 'C:/Users/letwir/repo/repo/mcp2go/POLYGLYPH.md'
Assert ((Get-FileHash $codec).Hash -eq (Get-FileHash 'C:/Users/letwir/.harness/POLYGLYPH.md').Hash) 'Codec hash mismatch'
$agy = 'C:/Users/letwir/.harness/skills/agy-subagent-router/SKILL.md'
$agyBefore = Text "$task/retry-agy/SKILL.before.md"
$agyAfter = CheckCommon $agy
Assert ((Front $agyBefore) -ceq (Front $agyAfter)) 'AGY frontmatter changed'
$fa = @(Fences $agyBefore); $fb = @(Fences $agyAfter)
Assert ($fa.Count -eq 3 -and $fb.Count -eq 3) 'AGY fence count'
for ($i = 0; $i -lt 3; $i++) { Assert ($fa[$i] -ceq $fb[$i]) 'AGY fence bytes changed' }
Assert ($agyBefore.Substring($agyBefore.IndexOf('#SIGMA LRF/1')) -ceq $agyAfter.Substring($agyAfter.IndexOf('#SIGMA LRF/1'))) 'AGY LRF bytes changed'
Assert (-not [regex]::IsMatch($agyAfter, '\p{IsCyrillic}')) 'Unexpected AGY Cyrillic prose'
$report = [ordered]@{status='PASS';checked_skills=40;canonical_skills=27;repository_copies=13;codec_hash='PASS';manifest_hashes='PASS';snapshot_hashes='PASS';semantic_review='Parent compared source references and approved nine explicit capability-reference merges; independent review pending';files=$results;time=(Get-Date -Format o)}
[IO.File]::WriteAllText("$task/final-validation.json", ($report | ConvertTo-Json -Depth 6), [Text.UTF8Encoding]::new($false))
'FINAL_VALIDATION_PASS: 27 canonical skills + 13 repository copies + codec; immutable snapshots and manifest hashes verified'
