import importlib.util, json, pathlib, unittest

P=pathlib.Path(__file__).resolve().parents[1]/'lint-polyglyph.py'
spec=importlib.util.spec_from_file_location('lint_polyglyph',P); lint=importlib.util.module_from_spec(spec); spec.loader.exec_module(lint)

class StructuralTests(unittest.TestCase):
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
        wrapper='<Γ>\nRole: x\nTarget: x\nAcceptance: x\nScope: x\nKnown facts: x\nSIGMA/1\n</Γ>'
        reordered='<Γ>\nTarget: x\nRole: x\nAcceptance: x\nScope: x\nKnown facts: x\nSIGMA/1\n</Γ>'
        self.assertTrue(lint.validate_gamma(wrapper))
        self.assertFalse(lint.validate_gamma(reordered))
        self.assertFalse(lint.validate_gamma(wrapper.replace('Scope: x','Scope: x\nRole: y')))

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
        prefix='@meta=text/POLYGLYPH×SIGMA/1\\n'+'\\n'.join(lint.META)+'\\n'
        valid='''#SIGMA LRF/1
@lrf=1|aud=x|scope=y
R|id|*|MAY|RO_LOCAL|payload

'''
        for suffix in ('@lrf=1|aud=x|scope=y\\n', '#SIGMA LRF/1\\n', 'L|id|*|BOGUS|RO_LOCAL|payload\\n'):
            data=(prefix+valid+suffix+'body\\n').encode()
            with self.assertRaises(ValueError): lint.parse('x',data)

    def test_commonmark_fence_openers(self):
        for opener in ('``` python`', '``~'):
            data=('\n'.join([lint.MARK,*lint.META,opener])+'\n').encode()
            self.assertEqual(lint.parse('x',data)['fences'], [])

if __name__=='__main__': unittest.main()
