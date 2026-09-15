package phonenumber

import (
	"errors"
	"testing"
)

func TestParseStrictAccepts(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Number
	}{
		{"e164 no dashes", "+12125550134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"e164 with dashes", "+1-212-555-0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"parens", "(212) 555-0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"bare dashes", "212-555-0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Parse(c.in, false)
			if err != nil {
				t.Fatalf("Parse(%q, false) returned error: %v", c.in, err)
			}
			if got != c.want {
				t.Fatalf("Parse(%q, false) = %+v, want %+v", c.in, got, c.want)
			}
		})
	}
}

func TestParseStrictRejects(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"dots instead of dashes", "212.555.0134"},
		{"no separators", "2125550134"},
		{"extra whitespace", "212 - 555 - 0134"},
		{"missing space after parens", "(212)555-0134"},
		{"trailing extension", "212-555-0134 x123"},
		{"area code starts with 0", "012-555-0134"},
		{"area code starts with 1", "112-555-0134"},
		{"exchange starts with 0", "212-055-0134"},
		{"exchange starts with 1", "212-155-0134"},
		{"n11 area code", "911-555-0134"},
		{"n11 exchange", "212-611-0134"},
		{"letters", "212-555-CALL"},
		{"too few digits", "212-555-013"},
		{"missing country code plus", "1-212-555-0134"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse(c.in, false)
			if err == nil {
				t.Fatalf("Parse(%q, false) succeeded, want error", c.in)
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q, false) returned %T, want *ParseError", c.in, err)
			}
		})
	}
}

func TestParseLenientAccepts(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want Number
	}{
		{"dots", "212.555.0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"bare digits", "2125550134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"leading country code", "12125550134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"spaces", "212 555 0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"mixed punctuation", "(212) 555.0134", Number{AreaCode: "212", Exchange: "555", Line: "0134"}},
		{"n11 exchange allowed", "212-611-0134", Number{AreaCode: "212", Exchange: "611", Line: "0134"}},
		{"extension via x", "212-555-0134x123", Number{AreaCode: "212", Exchange: "555", Line: "0134", Extension: "123"}},
		{"extension via ext.", "212-555-0134 ext. 45", Number{AreaCode: "212", Exchange: "555", Line: "0134", Extension: "45"}},
		{"extension via hash", "212-555-0134 #7", Number{AreaCode: "212", Exchange: "555", Line: "0134", Extension: "7"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Parse(c.in, true)
			if err != nil {
				t.Fatalf("Parse(%q, true) returned error: %v", c.in, err)
			}
			if got != c.want {
				t.Fatalf("Parse(%q, true) = %+v, want %+v", c.in, got, c.want)
			}
		})
	}
}

func TestParseLenientRejects(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"too few digits", "555-0134"},
		{"too many digits", "212-555-01345"},
		{"leading country code wrong digit", "22125550134"},
		{"area code starts with 0", "012-555-0134"},
		{"area code starts with 1", "112-555-0134"},
		{"exchange starts with 0", "212-055-0134"},
		{"exchange starts with 1", "212-155-0134"},
		{"letters with no digits", "call-me-now"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse(c.in, true)
			if err == nil {
				t.Fatalf("Parse(%q, true) succeeded, want error", c.in)
			}
		})
	}
}

func TestParseEmpty(t *testing.T) {
	for _, lenient := range []bool{false, true} {
		if _, err := Parse("", lenient); !errors.Is(err, ErrEmpty) {
			t.Fatalf("Parse(\"\", %v) returned %v, want ErrEmpty", lenient, err)
		}
		if _, err := Parse("   ", lenient); !errors.Is(err, ErrEmpty) {
			t.Fatalf("Parse(\"   \", %v) returned %v, want ErrEmpty", lenient, err)
		}
	}
}
