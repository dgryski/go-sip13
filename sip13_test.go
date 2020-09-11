package sip13

import "testing"

var want = []uint64{
	0xabac0158050fc4dc,
	0xc9f49bf37d57ca93,
	0x82cb9b024dc7d44d,
	0x8bf80ab8e7ddf7fb,
	0xcf75576088d38328,
	0xdef9d52f49533b67,
	0xc50d2b50c59f22a7,
	0xd3927d989bb11140,
	0x369095118d299a8e,
	0x25a48eb36c063de4,
	0x79de85ee92ff097f,
	0x70c118c1f94dc352,
	0x78a384b157b4d9a2,
	0x306f760c1229ffa7,
	0x605aa111c0f95d34,
	0xd320d86d2a519956,
	0xcc4fdd1a7d908b66,
	0x9cf2689063dbd80c,
	0x8ffc389cb473e63e,
	0xf21f9de58d297d1c,
	0xc0dc2f46a6cce040,
	0xb992abfe2b45f844,
	0x7ffe7b9ba320872e,
	0x525a0e7fdae6c123,
	0xf464aeb267349c8c,
	0x45cd5928705b0979,
	0x3a3e35e3ca9913a5,
	0xa91dc74e4ade3b35,
	0xfb0bed02ef6cd00d,
	0x88d93cb44ab1e1f4,
	0x540f11d643c5e663,
	0x2370dd1f8c21d1bc,
	0x81157b6c16a7b60d,
	0x4d54b9e57a8ff9bf,
	0x759f12781f2a753e,
	0xcea1a3bebf186b91,
	0x2cf508d3ada26206,
	0xb6101c2da3c33057,
	0xb3f47496ae3a36a1,
	0x626b57547b108392,
	0xc1d2363299e41531,
	0x667cc1923f1ad944,
	0x65704ffec8138825,
	0x24f280d1c28949a6,
	0xc2ca1cedfaf8876b,
	0xc2164bfc9f042196,
	0xa16e9c9368b1d623,
	0x49fb169c8b5114fd,
	0x9f3143f8df074c46,
	0xc6fdaf2412cc86b3,
	0x7eaf49d10a52098f,
	0x1cf313559d292f9a,
	0xc44a30dda2f41f12,
	0x36fae98943a71ed0,
	0x318fb34c73f0bce6,
	0xa27abf3670a7e980,
	0xb4bcc0db243c6d75,
	0x23f8d852fdb71513,
	0x8f035f4da67d8a08,
	0xd89cd0e5b7e8f148,
	0xf6f4e6bcf7a644ee,
	0xaec59ad80f1837f2,
	0xc3b2f6154b6694e0,
	0x9d199062b7bbb3a8,
}

var direct = []struct {
	h uint64
	s string
}{
	{0xac20658cb45c5a1e, "Hello World"},
	{0x25f46378bb4419d6, "Привет Всем!"},
	{0x63edad07d05747c2, "^A\v\xbf)\x85\xb2I\xb4\xf8\xa6\x14"},
	{0x7156a4faa386411d, "^A\v\xbf)\x85\xb2I\xb4\xf8\xa6\x00"},
}

func TestSip13(t *testing.T) {

	var k0 uint64 = 0x0706050403020100
	var k1 uint64 = 0x0f0e0d0c0b0a0908

	var p [64]byte

	for i := 0; i < 64; i++ {
		p[i] = byte(i)
		got := Sum64(k0, k1, p[:i])

		if got != want[i] {
			t.Errorf("Sum64([%d])=%08x, want %08x\n", i, got, want[i])
		}

		got = Sum64Str(k0, k1, string(p[:i]))
		if got != want[i] {
			t.Errorf("Sum64Str([%d])=%08x, want %08x\n", i, got, want[i])
		}
	}

	for _, d := range direct {
		got := Sum64Str(k0, k1, d.s)
		if got != d.h {
			t.Errorf("Sum64Str(%q)=%08x, want %08x\n", d.s, got, d.h)
		}

		got = Sum64(k0, k1, []byte(d.s))
		if got != d.h {
			t.Errorf("Sum64(%q)=%08x, want %08x\n", d.s, got, d.h)
		}
	}
}

var sum uint64

func BenchmarkSip13(b *testing.B) {
	var k0 uint64 = 0x0706050403020100
	var k1 uint64 = 0x0f0e0d0c0b0a0908
	var pb [64]byte
	for i := 0; i < 64; i++ {
		pb[i] = byte(i)
	}

	for k := 0; k < b.N; k++ {
		for i := 0; i < 64; i++ {
			sum += Sum64(k0, k1, pb[:i])
		}
	}
	b.SetBytes(int64((64 + 65) / 2))
}

func BenchmarkSip13Str(b *testing.B) {
	var k0 uint64 = 0x0706050403020100
	var k1 uint64 = 0x0f0e0d0c0b0a0908
	var pb [64]byte
	for i := 0; i < 64; i++ {
		pb[i] = byte(i)
	}
	p := string(pb[:])

	for k := 0; k < b.N; k++ {
		for i := 0; i < 64; i++ {
			sum += Sum64Str(k0, k1, p[:i])
		}
	}
	b.SetBytes(int64((64 + 65) / 2))
}
