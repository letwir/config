#!/usr/bin/env python3
"""Read-only structural POLYGLYPH lint; never semantic proof."""
import argparse, hashlib, json, pathlib, re, sys, unicodedata

ROOT = pathlib.Path(__file__).resolve().parents[1]
FIX = ROOT / 'scripts/tests/fixtures'
META = ['@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true', '@authority=SIGMA/1+routed-LRF-only;PG=representation-only', '@decode_failure=STOP;unknown=preserve;action=blocked;ask=true']
MARK = '@meta=text/POLYGLYPH×SIGMA/1'
FENCE = re.compile(r'^( {0,3})(`{3,}|~{3,})')

def fence_open(line):
    m=FENCE.match(line)
    if not m: return None
    run=m.group(2); info=line[m.end():]
    if run[0]=='`' and '`' in info: return None
    return m

def digest(data): return hashlib.sha256(data).hexdigest()
def parse(path, data, require_marker=True, validate_trailer=True):
    if data.startswith(b'\xef\xbb\xbf'): raise ValueError(f'{path}: BOM')
    text=data.decode('utf-8')
    if '\r' in text.replace('\r\n',''): raise ValueError(f'{path}: bare CR')
    lines=text.splitlines(keepends=True); bare=[x.rstrip('\r\n') for x in lines]
    i=0
    if bare and bare[0]=='---':
        try: i=bare.index('---',1)+1
        except ValueError: raise ValueError(f'{path}: unterminated frontmatter')
    fm=''.join(lines[:i]).encode()
    if require_marker:
        while i<len(lines) and bare[i]=='': i+=1
        if i>=len(lines) or bare[i]!=MARK: raise ValueError(f'{path}: missing marker')
        if bare[i+1:i+4]!=META: raise ValueError(f'{path}: bad bootstrap')
        start=i+4
    else: start=i
    fences=[]; lrf=None; opened=None; fencebytes=[]; regions=[]; n=start
    header_re=re.compile(r'@lrf=1\|aud=[^|]+\|scope=[^|]+')
    reserved=lambda x: x.startswith(('R|','L|','@lrf=','#SIGMA LRF/1'))
    while n<len(lines):
        line=bare[n]
        if opened:
            fencebytes.append(lines[n].encode())
            m=FENCE.match(line)
            if m and m.group(2)[0]==opened[0] and len(m.group(2))>=len(opened) and not line[m.end():].strip(' \t'):
                fences.append(digest(b''.join(fencebytes))); opened=None; fencebytes=[]
            n+=1; continue
        m=fence_open(line)
        if m:
            opened=m.group(2); fencebytes=[lines[n].encode()]; n+=1; continue
        if line=='#SIGMA LRF/1':
            if lrf is not None: raise ValueError(f'{path}: duplicate LRF marker')
            begin=n; n+=1
            while n<len(lines) and not bare[n].strip(): n+=1
            if n>=len(lines) or not header_re.fullmatch(bare[n]): raise ValueError(f'{path}: bad LRF header')
            n+=1
            while n<len(lines):
                row=bare[n]
                if not row.strip(): n+=1; continue
                if row.startswith(('R|','L|')):
                    f=row.split('|')
                    if len(f)!=6 or f[0] not in ('R','L') or not f[1] or not f[2] or f[3] not in ('MUST','MUST_NOT','MAY') or not f[4] or not f[5]: raise ValueError(f'{path}: bad LRF record')
                    n+=1; continue
                if row.startswith(('@lrf=','#SIGMA LRF/1')): raise ValueError(f'{path}: duplicate/malformed LRF header')
                break
            end=n
            # At EOF the block owns every trailing separator byte; before resumed prose,
            # separator blanks also remain part of the immutable region.
            if n<len(lines):
                while end>begin and not bare[end-1].strip(): end-=1
                if end<n: end=n
            lrf=digest(''.join(lines[begin:end]).encode())
            continue
        if reserved(line): raise ValueError(f'{path}: reserved LRF prefix outside block')
        n+=1
    if opened: raise ValueError(f'{path}: unclosed fence')
    return {'frontmatter':digest(fm) if fm else None,'fences':fences,'lrf':lrf if lrf else None}

def validate_sigma(value):
    if not isinstance(value,str): return False
    value=value.replace('\r\n','\n')
    return (value.startswith('[SIGMA/1]\n<∑>\n') and value.endswith('\n</∑>')
            and '\r' not in value
            and value.count('<∑>') == 1 and value.count('</∑>') == 1
            and value.count('[SIGMA/1]') == 1 and '[/SIGMA/1]' not in value)

def validate_handoff(value):
    if not validate_sigma(value) or len(value.encode('utf-16-le'))//2>20_000: return False
    body=value.replace('\r\n','\n').split('\n')[2:-1]
    fields=('Role:','Target:','Acceptance:','Scope:','Known facts:')
    reserved=re.compile(r'(?:\[SIGMA/1\]|SIGMA/1(?::.*)?|∑/1: POLYGLYPH|</?[∑Γ]>)')
    if any(reserved.fullmatch(line) for line in body): return False
    rows=[line for line in body if line.startswith(fields)]
    return (len(rows)==len(fields) and body[0]==rows[0]
            and all(row.startswith(field+' ') and row[len(field)+1:].strip()
                    for row,field in zip(rows,fields)))

def inventory():
    paths=set(['AGENTS.md','rules/README-forLLM.md'])
    paths.update(p.relative_to(ROOT).as_posix() for p in (ROOT/'agents').glob('*.md') if p.name!='README.md')
    paths.update(p.relative_to(ROOT).as_posix() for p in (ROOT/'skills').rglob('SKILL.md') if 'node_modules' not in p.parts)
    paths.update(p.relative_to(ROOT).as_posix() for p in (ROOT/'opencode/agents').glob('*.md'))
    return paths

def check(root=ROOT):
    targets=json.loads((FIX/'polyglyph-targets.json').read_text())['targets']; frozen=json.loads((FIX/'polyglyph-preservation.json').read_text())['files']
    if len(targets)!=len(set(targets)) or set(targets)!=inventory(): raise ValueError('manifest/inventory mismatch')
    if set(frozen)!=set(targets): raise ValueError('preservation fixture inventory mismatch')
    remaining=[]
    for key in targets:
        p=root/key
        if p.is_symlink() or not p.is_file(): raise ValueError(f'{key}: missing/non-regular/symlink')
        data=p.read_bytes()
        sig=parse(key,data,require_marker=False,validate_trailer=True)
        if sig!=frozen[key]: raise ValueError(f'{key}: immutable region differs from baseline')
        try:
            parse(key,data,validate_trailer=True)
        except ValueError as e:
            if str(e).endswith(': missing marker'):
                remaining.append(key)
            else:
                raise
    return remaining

def main():
    ap=argparse.ArgumentParser(description=__doc__); ap.add_argument('--root',type=pathlib.Path,default=ROOT); ap.add_argument('--targets',nargs='*'); a=ap.parse_args()
    try:
        remaining=check(a.root)
        selected=a.targets or ['AGENTS.md','rules/README-forLLM.md']
        for t in selected: parse(t,(a.root/t).read_bytes())
        print(f'STRUCTURAL ONLY, NOT SEMANTIC PROOF: PASS selected={len(selected)}; incomplete targets={len(remaining)} paths={json.dumps(remaining)}')
        return 0
    except (ValueError, OSError, KeyError, json.JSONDecodeError) as e: print(f'FAIL: {e}',file=sys.stderr); return 1
if __name__=='__main__': raise SystemExit(main())
