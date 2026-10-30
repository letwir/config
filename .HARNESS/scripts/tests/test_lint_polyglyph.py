import importlib.util, json, pathlib, re, tomllib, unittest

P=pathlib.Path(__file__).resolve().parents[1]/'lint-polyglyph.py'
spec=importlib.util.spec_from_file_location('lint_polyglyph',P)
if spec is None or spec.loader is None: raise ImportError('could not load lint_polyglyph')
lint=importlib.util.module_from_spec(spec); spec.loader.exec_module(lint)

class StructuralTests(unittest.TestCase):
    def test_inventory_paths_are_posix_and_fixture_consistent(self):
        paths=lint.inventory()
        self.assertTrue(all('\\' not in path for path in paths))
        targets=json.loads((lint.FIX/'polyglyph-targets.json').read_text())['targets']
        frozen=json.loads((lint.FIX/'polyglyph-preservation.json').read_text())['files']
        self.assertEqual(len(targets),len(set(targets)))
        self.assertEqual(set(targets),paths)
        self.assertEqual(set(targets),set(frozen))

    def test_frozen_regions_match_missing_marker_targets(self):
        remaining=lint.check()
        targets=json.loads((lint.FIX/'polyglyph-targets.json').read_text())['targets']
        expected=[]
        for key in targets:
            try: lint.parse(key,(lint.ROOT/key).read_bytes())
            except ValueError as e:
                if str(e).endswith(': missing marker'):
                    expected.append(key)
                else:
                    raise
        self.assertEqual(remaining,expected)
    def test_system_skills_are_migrated_with_frozen_regions(self):
        frozen=json.loads((lint.FIX/'polyglyph-preservation.json').read_text())['files']
        remaining=lint.check()
        for name in ('imagegen','openai-docs','review-agent','skill-creator','skill-installer'):
            key=f'skills/.system/{name}/SKILL.md'
            with self.subTest(path=key):
                self.assertIn(key,lint.inventory())
                self.assertNotIn(key,remaining)
                self.assertEqual(lint.parse(key,(lint.ROOT/key).read_bytes()),frozen[key])

    def test_marker_metadata_and_frontmatter_position(self):
        def parse(body): return lint.parse('x',body.encode())
        with self.assertRaisesRegex(ValueError,'missing marker'):
            parse('body\n')
        with self.assertRaisesRegex(ValueError,'bad bootstrap'):
            parse(lint.MARK+'\n'+lint.META[0]+'\nwrong\n'+lint.META[2]+'\n')
        valid='---\ntitle: x\n---\n\n'+lint.MARK+'\n'+'\n'.join(lint.META)+'\nbody\n'
        self.assertTrue(parse(valid))
        with self.assertRaisesRegex(ValueError,'missing marker'):
            parse('---\ntitle: x\n---\nbody\n'+lint.MARK+'\n'+'\n'.join(lint.META)+'\n')
    def test_valid_marker_and_fence(self):
        b=('''@meta=text/POLYGLYPH×SIGMA/1\n'''+ '\n'.join(lint.META)+'''\nbody\n```x\n#SIGMA LRF/1\n```\n''').encode()
        self.assertEqual(len(lint.parse('x',b)['fences']),1)
    def test_rejects_bad_bytes_marker_and_unclosed_fence(self):
        for b in (b'\xef\xbb\xbf', b'bad\rbytes', ('\n'.join([lint.MARK]+lint.META+['```'])).encode()):
            with self.assertRaises((ValueError,UnicodeDecodeError)): lint.parse('x',b)
    def test_sigma_exact_unicode_and_wrapper_order(self):
        good='[SIGMA/1]\n<∑>\nContract\n</∑>'
        self.assertTrue(lint.validate_sigma(good))
        for bad in ('[SIGMA/1]\\n<∑>\\nContract\\n</∑>', 'Σ[SIGMA/1]\n<∑>\nContract\n</∑>', 'prefix\n'+good, good+'\nfooter', good+'\n</∑>', '[SIGMA/1]\n<∑>\nContract\n</Σ>', '[SIGMA/1]\n<∑>\nContract\n</∑>\n[SIGMA/1]'):
            self.assertFalse(lint.validate_sigma(bad), bad)
        wrapper='[SIGMA/1]\n<∑>\nRole: x\nTarget: x\nAcceptance: x\nScope: x\nKnown facts: x\n</∑>'
        reordered='[SIGMA/1]\n<∑>\nTarget: x\nRole: x\nAcceptance: x\nScope: x\nKnown facts: x\n</∑>'
        self.assertTrue(lint.validate_handoff(wrapper))
        self.assertFalse(lint.validate_handoff(reordered))
        self.assertFalse(lint.validate_handoff(wrapper.replace('Scope: x','Scope: x\nRole: y')))
        self.assertFalse(lint.validate_handoff(wrapper.replace('[SIGMA/1]', '[SIGMA/1] extra')))
        self.assertTrue(lint.validate_handoff(wrapper.replace('\n','\r\n')))
        self.assertTrue(lint.validate_handoff(wrapper.replace('Target: x','Target: Γ and Σ\ncontinuation')))
        for bad in (wrapper+'\n', 'prefix\n'+wrapper, wrapper+'\nfooter',
                    wrapper.replace('<∑>','<Σ>'), wrapper.replace('</∑>','</Σ>'),
                    wrapper.replace('Target: x','Target: '),
                    wrapper.replace('Target: x','Target: x\rbroken'),
                    wrapper.replace('<∑>\n','<∑>\nextra header\n'),
                    wrapper.replace('<∑>\n','<∑>\n[SIGMA/1]\n'),
                    wrapper.replace('\n</∑>','\nSIGMA/1\n</∑>'),
                    wrapper.replace('[SIGMA/1]\n<∑>','<Γ>\n∑/1: POLYGLYPH').replace('</∑>','SIGMA/1\n</Γ>'),
                    wrapper.replace('Target: x','Target: '+'🧪'*10_000)):
            self.assertFalse(lint.validate_handoff(bad),bad[:120])

    def test_migrated_lrf_records_stay_on_one_physical_line(self):
        for key in ('rules/MANUAL.lrf','rules/agy.lrf','agents/SUBAGENTS.lrf'):
            lines=(lint.ROOT/key).read_text(encoding='utf-8').splitlines()
            self.assertRegex(lines[0],r'^@lrf=1\|aud=[^|]+\|scope=[^|]+$')
            for number,line in enumerate(lines[1:],2):
                if not line.strip() or line.startswith('#'): continue
                fields=line.split('|',5)
                self.assertEqual(len(fields),6,(key,number))
                self.assertIn(fields[0],('R','L'),(key,number))
                self.assertIn(fields[3],('MUST','MUST_NOT','MAY'),(key,number))
                self.assertTrue(all(fields),(key,number))

    def test_canonical_unfenced_templates_and_legacy_preservation(self):
        pattern=re.compile(r'^\[SIGMA/1\]\n<∑>\nRole: [^\n]+\nTarget: [^\n]+\nAcceptance: [^\n]+\nScope: [^\n]+\nKnown facts: [^\n]+\n</∑>$',re.M)
        for key in ('agents/worker.md','agents/CALL.md','skills/agy-subagent-router/SKILL.md'):
            text=(lint.ROOT/key).read_text(encoding='utf-8')
            example=pattern.search(text)
            self.assertIsNotNone(example,key)
            self.assertTrue(lint.validate_handoff(example.group()),key)
            self.assertIn('legacy',text)
            self.assertIn('SIGMA/1: PIDGEN/text',text)

    def test_role_toml_uses_canonical_envelope_without_extra_header(self):
        for name in ('auditor','blackhat','critic','explorer','filereader','proposer',
                     'refactorer','researcher','verifier-blackhat','verifier'):
            config=tomllib.loads((lint.ROOT/f'agents/{name}.toml').read_text(encoding='utf-8'))
            prompt=config['developer_instructions']
            self.assertIn('[SIGMA/1] then <∑>',prompt,name)
            self.assertIn('closing line </∑>',prompt,name)
            self.assertIn('no extra header/footer or standalone terminal SIGMA/1',prompt,name)

    def test_lrf_record_validation(self):
        b=('\n'.join([lint.MARK,*lint.META,'#SIGMA LRF/1','@lrf=1|aud=x|scope=y','R|id|*|MAY|RO_LOCAL|x'])+'\n').encode()
        self.assertTrue(lint.parse('x',b)['lrf'])
        b=b.replace(b'R|id|*|MAY|RO_LOCAL|x',b'R|id|*|MAY|RO_LOCAL')
        with self.assertRaises(ValueError): lint.parse('x',b)
        for header in (b'@lrf=1|aud=x|scope=y|extra', b'@lrf=1|scope=y|aud=x'):
            bad=b.replace(b'@lrf=1|aud=x|scope=y',header).replace(b'R|id|*|MAY|RO_LOCAL',b'R|id|*|MAY|RO_LOCAL|x')
            with self.assertRaises(ValueError): lint.parse('x',bad)

    def test_lrf_eof_preserves_trailing_blank_bytes(self):
        for trailing in ('\n', '\n\n', '\n\n\n'):
            body='\n'.join([lint.MARK,*lint.META,'#SIGMA LRF/1','@lrf=1|aud=x|scope=y','R|id|*|MAY|RO_LOCAL|payload'])+trailing
            result=lint.parse('x',body.encode())
            marker_through_eof=body[body.index('#SIGMA LRF/1'):].encode()
            self.assertEqual(result['lrf'],lint.digest(marker_through_eof))

    def test_lrf_embedded_before_body_and_fence(self):
        body='''@meta=text/POLYGLYPH×SIGMA/1
'''+ '\n'.join(lint.META)+'''\nIntroduction\n#SIGMA LRF/1\n@lrf=1|aud=x|scope=y\nR|id|*|MAY|RO_LOCAL|payload\n\n## Guidance\nMore prose\n```text\nfenced\n```\n'''
        result=lint.parse('x',body.encode())
        self.assertTrue(result['lrf'])
        self.assertEqual(len(result['fences']),1)

    def test_rejects_duplicate_lrf_headers_markers_and_bad_rows(self):
        prefix='@meta=text/POLYGLYPH×SIGMA/1\n'+'\n'.join(lint.META)+'\n'
        valid='''#SIGMA LRF/1
@lrf=1|aud=x|scope=y
R|id|*|MAY|RO_LOCAL|payload

'''
        for suffix in ('@lrf=1|aud=x|scope=y\n', '#SIGMA LRF/1\n', 'L|id|*|BOGUS|RO_LOCAL|payload\n'):
            data=(prefix+valid+suffix+'body\n').encode()
            with self.assertRaises(ValueError): lint.parse('x',data)

    def test_commonmark_fence_openers(self):
        for opener in ('``` python`', '``~'):
            data=('\n'.join([lint.MARK,*lint.META,opener])+'\n').encode()
            self.assertEqual(lint.parse('x',data)['fences'], [])

if __name__=='__main__': unittest.main()
