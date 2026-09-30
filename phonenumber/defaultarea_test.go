package phonenumber

import "testing"

func TestParseWithDefaultAreaFillsLocalNumbers(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Number
	}{
		{"dashed local", "555-0134", Number{AreaCode: "415", Exchange: "555", Line: "0134"}},
		{"bare local", "5550134", Number{AreaCode: "415", Exchange: "555", Line: "0134"}},
		{"local with extension", "555-0134 x9", Number{AreaCode: "415", Exchange: "555", Line: "0134", Extension: "9"}},
		{"full number keeps its own area code", "(212) 555-0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"eleven digits keeps its own area code", "1-212-555-0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseWithDefaultArea(c.in, "415")
			if err != nil {
				t.Fatalf("ParseWithDefaultArea(%q, \"415\") returned error: %v", c.in, err)
			}
			if got != c.want {
				t.Fatalf("ParseWithDefaultArea(%q, \"415\") = %+v, want %+v", c.in, got, c.want)
			}
		})
	}
}

func TestParseWithDefaultAreaEmptyDefaultRejectsLocal(t *testing.T) {
	if _, err := ParseWithDefaultArea("555-0134", ""); err == nil {
		t.Fatal("expected an error for a 7-digit number with no default area code")
	}
}

func TestParseWithDefaultAreaRejectsBadDefault(t *testing.T) {
	for _, area := range []string{"41", "4155", "abc", "4a5", "015", "115"} {
		if _, err := ParseWithDefaultArea("555-0134", area); err == nil {
			t.Errorf("default area %q: expected an error", area)
		}
	}
}

func TestParseWithDefaultAreaStillValidatesExchange(t *testing.T) {
	if _, err := ParseWithDefaultArea("155-0134", "415"); err == nil {
		t.Fatal("expected an error for an exchange starting with 1")
	}
}

func TestParseWithDefaultAreaEmptyInput(t *testing.T) {
	if _, err := ParseWithDefaultArea("  ", "415"); err != ErrEmpty {
		t.Fatalf("got %v, want ErrEmpty", err)
	}
}
