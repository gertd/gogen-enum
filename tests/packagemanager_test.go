package tests

import (
	"strings"
	"testing"
)

func TestPackageManager(t *testing.T) {

	for e, s := range packageManagerID {

		p := NewPackageManager(strings.ToUpper(s))

		if p.String() != e.String() && p != e {
			t.Errorf("case insensitive lookup failed expected %s actual %s", p.String(), e.String())
		}
	}

}
