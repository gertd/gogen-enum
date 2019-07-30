package tests

import "testing"

func TestArchitecture(t *testing.T) {

	{
		// default to unknown
		var a Architecture
		if a != ArchitectureUnknown {
			t.Errorf("implicit default value variable assignment failure")
		}
	}
	{
		a := ArchitectureUnknown
		if a.String() != "unknown" {
			t.Errorf("explicit value variable assignment failure")
		}
	}
	{
		a := ArchitectureARM
		b := ArchitectureARM
		if a != b {
			t.Errorf("equality check failure")
		}
	}
	{
		a := ArchitectureARM
		b := ArchitectureX64
		if a == b {
			t.Errorf("inequality check failure")
		}
	}

}
