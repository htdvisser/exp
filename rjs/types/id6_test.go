package types

import (
	"testing"
)

func TestID6(t *testing.T) {
	tt := []struct {
		desc      string
		id6       ID6
		euiString string
		macString string
		id6String string
	}{
		{
			desc:      "AABB0001020342FF",
			id6:       ID6{0xaa, 0xbb, 0x00, 0x01, 0x02, 0x03, 0x42, 0xff},
			euiString: "AA-BB-00-01-02-03-42-FF",
			id6String: "aabb:1:203:42ff",
		},
		{
			desc:      "0000000000000000",
			id6:       ID6{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			euiString: "00-00-00-00-00-00-00-00",
			id6String: "::0",
		},
		{
			desc:      "58A0CBFFFE800DBD",
			id6:       ID6{0x58, 0xa0, 0xcb, 0xff, 0xfe, 0x80, 0x0d, 0xbd},
			euiString: "58-A0-CB-FF-FE-80-0D-BD",
			macString: "58-A0-CB-80-0D-BD",
			id6String: "58a0:cbff:fe80:dbd",
		},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.euiString != "" {
				var got ID6
				err := got.UnmarshalText([]byte(tc.euiString))
				if err != nil {
					t.Errorf("Unmarshaling %q failed: %v", tc.euiString, err)
				}
				if got != tc.id6 {
					t.Errorf("Unmarshaled %q did not match expected %v, but was %v", tc.euiString, tc.id6, got)
				}
			}
			if tc.macString != "" {
				var got ID6
				err := got.UnmarshalText([]byte(tc.macString))
				if err != nil {
					t.Errorf("Unmarshaling %q failed: %v", tc.macString, err)
				}
				if got != tc.id6 {
					t.Errorf("Unmarshaled %q did not match expected %v, but was %v", tc.macString, tc.id6, got)
				}
			}
			if tc.id6String != "" {
				var got ID6
				err := got.UnmarshalText([]byte(tc.id6String))
				if err != nil {
					t.Errorf("Unmarshaling %q failed: %v", tc.id6String, err)
				}
				if got != tc.id6 {
					t.Errorf("Unmarshaled %q did not match expected %v, but was %v", tc.id6String, tc.id6, got)
				}

				b, err := tc.id6.MarshalText()
				if err != nil {
					t.Errorf("Marshaling %v failed: %v", tc.id6, err)
				}
				if string(b) != tc.id6String {
					t.Errorf("Marshaled %v did not match expected %q, but was %q", tc.id6, tc.id6String, string(b))
				}
			}
		})
	}
}
