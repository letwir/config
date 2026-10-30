$ErrorActionPreference = 'Stop'
$root = 'C:/Users/letwir/.harness/memo/20261025-polyglyph-skills/batch'
$repo = 'C:/Users/letwir/repo/repo/mcp2go'
$canonicalRoot = 'C:/Users/letwir/.harness/skills'
$snapshotPath = Join-Path $root 'snapshots.json'
$snapshots = Get-Content -Raw -LiteralPath $snapshotPath | ConvertFrom-Json
$byPath = @{}
foreach ($s in $snapshots) { $byPath[[IO.Path]::GetFullPath($s.path).ToLowerInvariant()] = $s }
$utf8 = [Text.UTF8Encoding]::new($false)
function Read-Text([string]$p) { [IO.File]::ReadAllText($p, [Text.Encoding]::UTF8) }
function Write-Text([string]$p, [string]$t) { [IO.File]::WriteAllText($p, $t, $utf8) }
function Sha([string]$p) { $h = [Security.Cryptography.SHA256]::HashData([IO.File]::ReadAllBytes($p)); [Convert]::ToHexString($h) }
function Front([string]$text) {
  if ($text -match '\A(---\r?\n[\s\S]*?\r?\n---(?:\r?\n|\z))') { return $Matches[1] }
  return ''
}
function Insert-Reference([string]$text, [string]$clause) {
  $nl = if ($text.Contains("`r`n")) { "`r`n" } else { "`n" }
  $index = $text.IndexOf('#SIGMA LRF/1', [StringComparison]::Ordinal)
  if ($index -lt 0) { $index = $text.Length }
  $prefix = $text.Substring(0, $index).TrimEnd("`r", "`n")
  $suffix = $text.Substring($index)
  return $prefix + $nl + $clause + $nl + $nl + $suffix
}
# Restore canonical files deterministically from snapshots, then add approved reference-only PG clauses before trailing LRF.
$referenceClauses = @{
  'economics-api' = 'Direct sibling CLI entry points (installed path must be resolved): estat.exe, edinet.exe, sec_edgar.exe, fred.exe, worldbank.exe, yfinance.exe, alphavantage.exe. Reference options: yfinance --range 1mo; alphavantage --function CURRENCY_EXCHANGE_RATE --from USD --to JPY. Documented keyless sources: SEC, WorldBank, yfinance. Missing-key guidance: consult available local USAGE.md only if verified present; never use stale plugin file URLs. No runtime installation claim; canonical gates unchanged.'
  'kokkousyou-api' = 'Option reference: -c/--count; -d/--dataset; -p/--prefecture; -j/--json; --limit defaults to 10; --around radius is meters. For LLM inspection of result data, request -j/--json. Documentation only; no API call or runtime verification.'
  'llama2coder' = 'Option reference (installed-supported options only; no override authority): --prompt/-p required; --language/-l default python3; --url/-u default http://localhost:8080; --model default default; --stream/-s default false; --system default empty; --temperature default 0.2; --max-tokens default -1. Prompt is wrapped in <instruction>; client removes fences; stdout is raw code without intro/outro unless explicitly requested. Canonical routing and permission gates remain controlling.'
  'llm-memory' = 'Capability model: BiTemporal(ValidTime×TxTime), progressive compaction(L0→L1→L2→L3), self-edit outcomes{ADD,UPDATE/Supersede,NOOP,DEPRECATE}, graph triples; canonical executable-resolution and approval contract controls. Reference commands: ingest -title <title> -text <text> -cat <cat>; search -tag <tag> -level 2; graph edge -src <S> -tgt <T> -rel <R>. Walkthrough/report ingest hooks are approval-gated references, never automatic arbitrary live writes. eval/evals remain as documented canonically.'
  'perser-xml' = 'Extraction reference: output keys={title,description,site_name,ogp}; all application/ld+json blocks → array; DMM/FANZA fields={maker,SKU/item#,release-date} where present. Canonical runtime and script invocation remain controlling.'
  'rikei-search' = 'Optional direct sub-binary references: mp-search.exe, aflow-search.exe, pubchem-search.exe, ieee-search.exe, ads-search.exe, cod-search.exe; resolve installed binaries before use. Key mapping: MP_API_KEY→MP; IEEE_API_KEY→IEEE; ADS_API_TOKEN→ADS. Missing-key fallback to keyless AFLOW/PubChem/COD and available local USAGE.md guidance only if verified present; never expose/read key values during migration.'
  'sec-forensics' = 'Reference options: cve-ingest may include -cat <cat> and -tags <t>; mega-check may include -transcript <path>. Possible output label VULNERABLE_WITH_MITIGATION is an observed label, not proof. Existing approval gates unchanged; no forensic scan performed.'
  'vm-monitor' = 'Reference queries (preserve literal documented expressions): CPU=`100 - (avg by (instance) (rate(node_cpu_seconds_total{mode!="idle"}[5m])) * 100)`; RAM=`node_memory_MemAvailable_bytes / (1024 * 1024 * 1024)`; disk=`node_filesystem_avail_bytes{mountpoint="/"} / (1024 * 1024 * 1024)`; load=`node_load1`; cache=`vm_cache_requests_total{type="metric_id"} > 0`; overload=`node_load1 > count by(instance)(node_cpu_seconds_total{mode=''idle''})`. -nodes/-top and -suggest/-capacity may request -format json. -suggest -format json → CapacityRank==1 evidence → deployment recommendation only, never deployment authorization. Candidate workloads={LLM_Inference,BatchPipeline,Microservices,EdgeDaemons}; no query execution.'
  'hudousan-api' = 'Legacy configuration-loader reference: KOKKOU_HUDOUSAN_API may be sourced from process environment or Windows user-registry configuration. Registry fallback must be verified in the selected installed binary before reliance. Never query or print key values during migration; canonical API validation/gates unchanged.'
}
$canonical = @($snapshots | Where-Object { [IO.Path]::GetFullPath($_.path).StartsWith([IO.Path]::GetFullPath($canonicalRoot), [StringComparison]::OrdinalIgnoreCase) })
foreach ($s in $canonical) {
  $backup = $s.backup
  $original = [IO.File]::ReadAllText($backup, [Text.Encoding]::UTF8)
  $converted = $original.Replace('[SPR/PIDGIN::','[POLYGLYPH::').Replace('[SPR/XML::','[POLYGLYPH::')
  if ($referenceClauses.ContainsKey($s.name)) { $converted = Insert-Reference $converted $referenceClauses[$s.name] }
  Write-Text $s.path $converted
}
$targets = @(
  @{ name='default-settings'; rel='LTWskills/skills/default-settings/SKILL.md' },
  @{ name='economics-api'; rel='src/economics-api/SKILL.md' },
  @{ name='hudousan-api'; rel='LTWskills/skills/hudousan-api/SKILL.md' },
  @{ name='kokkousyou-api'; rel='LTWskills/skills/kokkousyou-api/SKILL.md' },
  @{ name='llama2coder'; rel='LTWskills/skills/llama2coder/SKILL.md' },
  @{ name='llm-memory'; rel='src/llm-memory/SKILL.md' },
  @{ name='llm-memory'; rel='LTWskills/skills/llm-memory/SKILL.md' },
  @{ name='perser-xml'; rel='LTWskills/skills/perser-xml/SKILL.md' },
  @{ name='rikei-search'; rel='src/rikei-search/SKILL.md' },
  @{ name='rikei-search'; rel='LTWskills/skills/rikei-search/SKILL.md' },
  @{ name='sec-forensics'; rel='src/sec-forensics/SKILL.md' },
  @{ name='vm-monitor'; rel='src/vm-monitor/SKILL.md' },
  @{ name='vm-monitor'; rel='LTWskills/skills/vm-monitor/SKILL.md' }
)
$records = [Collections.Generic.List[object]]::new()
foreach ($s in $canonical) {
  $records.Add([pscustomobject]@{path=$s.path; action=$(if ($referenceClauses.ContainsKey($s.name)) {'snapshot-label-repair+approved-PG-reference-merge'} else {'snapshot-label-repair'}); before_sha256=$s.sha256; after_sha256=(Sha $s.path); baseline_backup=$s.backup; notes=$(if ($s.name -eq 'default-settings') {'Reused canonical PG; stale plugin locations/GitBash preference deliberately not restored because active_profile and canonical resolver govern.'} elseif ($referenceClauses.ContainsKey($s.name)) {'Added original documented reference options/capabilities as PG prose before LRF; canonical permission/resolver gates remain controlling.'} else {'Reused converted prose; label-only update restored with original snapshot line endings.'})})
}
foreach ($t in $targets) {
  $target = Join-Path $repo $t.rel
  $canonicalPath = Join-Path $canonicalRoot ($t.name + '/SKILL.md')
  $local = Read-Text $target
  $front = Front $local
  if (-not $front) { throw "Missing/invalid YAML frontmatter in $target" }
  $canonicalText = Read-Text $canonicalPath
  $body = $canonicalText.Substring((Front $canonicalText).Length)
  # Preserve each repository file's exact existing frontmatter bytes/text; synchronize converted canonical body.
  Write-Text $target ($front + $body)
  $key = [IO.Path]::GetFullPath($target).ToLowerInvariant()
  $before = if ($byPath.ContainsKey($key)) { $byPath[$key].sha256 } else { 'baseline-not-in-snapshot' }
  $records.Add([pscustomobject]@{path=$target; action='frontmatter-preserved+canonical-PG-body-sync'; before_sha256=$before; after_sha256=(Sha $target); baseline_backup=($(if ($byPath.ContainsKey($key)) { $byPath[$key].backup } else { $null })); notes=$(if ($t.name -eq 'llm-memory') {'Preserved BiTemporal/compaction/self-editing/graph capability in canonical PG; canonical resolver and approval gates intentionally supersede legacy plugin executable path; eval/evals covered by canonical current contract.'} else {'Preserved target YAML frontmatter; reused canonical converted body and contract.'}) })
}
$codecSource = 'C:/Users/letwir/.harness/POLYGLYPH.md'
$codecTarget = Join-Path $repo 'POLYGLYPH.md'
[IO.File]::Copy($codecSource, $codecTarget, $true)
$records.Add([pscustomobject]@{path=$codecTarget; action='byte-identical-codec-copy'; before_sha256='new-file'; after_sha256=(Sha $codecTarget); baseline_backup=$null; notes='Copied canonical codec byte-for-byte; required by nearest_ancestor bootstrap.'})
$manifest = [pscustomobject]@{files=@($records); localOnlyCapabilityNotes=@('llm-memory: retained BiTemporal valid/transaction time, progressive compaction, self-edit outcomes, graph triples, ingest/search/graph command forms, and approval-gated walkthrough/report hooks; canonical resolver and approval gates supersede obsolete plugin-specific path.','llm-memory eval/evals are represented by canonical current contract.','default-settings: stale plugin locations and GitBash preference intentionally superseded by active_profile and canonical resolver.','Other local reference details are preserved in canonical PG prose as documented options/examples only; installed support is unverified and no skill operations were run.'); counts=@{canonicalSkills=$canonical.Count;repoSkillCopies=$targets.Count;codecCopies=1}}
Write-Text (Join-Path $root 'manifest.json') ($manifest | ConvertTo-Json -Depth 8)
