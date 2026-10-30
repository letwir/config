$ErrorActionPreference = 'Stop'
$root = 'C:/Users/letwir/.harness/memo/20261025-polyglyph-skills/review'
[IO.Directory]::CreateDirectory($root) | Out-Null
$strict = [Text.UTF8Encoding]::new($false, $true)
$invalid = Join-Path $root 'invalid-utf8.fixture'
[IO.File]::WriteAllBytes($invalid, [byte[]](0x61,0xC3,0x28))
$rejected = $false
try { $null = $strict.GetString([IO.File]::ReadAllBytes($invalid)) } catch [Text.DecoderFallbackException] { $rejected = $true }
if (-not $rejected) { throw 'FAIL: invalid UTF-8 fixture accepted' }
$before = 'C:/Users/letwir/.harness/memo/20261025-polyglyph-skills/retry-agy/SKILL.before.md'
$after = 'C:/Users/letwir/.harness/skills/agy-subagent-router/SKILL.md'
$a = $strict.GetString([IO.File]::ReadAllBytes($before))
$b = $strict.GetString([IO.File]::ReadAllBytes($after))
$marker = '#SIGMA LRF/1'
$ai = $a.IndexOf($marker, [StringComparison]::Ordinal)
$bi = $b.IndexOf($marker, [StringComparison]::Ordinal)
if ($ai -lt 0 -or $bi -lt 0) { throw 'Fixture source lacks LRF marker' }
$la = $a.Substring($ai)
$lb = $b.Substring($bi)
$mutated = $lb.Replace('R|bin.resolve', 'R|bin.resolvX')
if ($mutated -ceq $lb) { throw 'FAIL: mutation did not alter LRF equality' }
if ($lb -cne $la) { throw 'FAIL: current LRF differs from immutable source before mutation' }
'NEGATIVE_TEST_PASS: strict UTF-8 decoder rejects invalid bytes; mutated LRF fails exact source equality'
