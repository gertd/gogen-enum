package tests

import (
	"testing"
)

func TestBitMask(t *testing.T) {

	var b, v uint8
	var p = PluralizeUnknown

	for e := range pluralizeID {

		if e.Has(e) != true {
			t.Errorf("must be true")
		}

		b = uint8(e)

		p.Set(e)

		v = uint8(p)
		if v&b == 0 {
			t.Errorf("bit not set expected %d actual %d", v, b)
		}

		if !p.Has(e) {
			t.Errorf("Has check failed expected %t actual %t", true, p.Has(e))
		}

		p.Clear(e)

		v = uint8(p)
		if v&b != 0 {
			t.Errorf("bit set expected %d actual %d", v, b)
		}

		p.Toggle(e)

		v = uint8(p)
		if v&b == 0 {
			t.Errorf("bit not set expected %d actual %d", v, b)
		}

		p.Toggle(e)

		v = uint8(p)
		if v&b != 0 {
			t.Errorf("bit set expected %d actual %d", v, b)
		}

		if s, ok := pluralizeID[p]; ok && s != p.String() {
			t.Errorf("String() != pluralizeID[] expected %s actual %s", p.String(), s)
		}

		if pluralizeID[pluralizeName[p.String()]] != p.String() {
			t.Errorf("roundtripping failed expected %s actual %s", pluralizeID[pluralizeName[p.String()]], p.String())
		}

		p.Set(e)

		v = uint8(p)
		if v&b == 0 {
			t.Errorf("bit not set expected %d actual %d", v, b)
		}

		p.Set(e)
		v = uint8(p)
		if v&b == 0 {
			t.Errorf("bit not set expected %d actual %d", v, b)
		}

	}

}
