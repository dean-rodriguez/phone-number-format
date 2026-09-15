package phonenumber

import "testing"

func TestPretty(t *testing.T) {
	cases := []struct {
		name string
		n    Number
		want string
	}{
		{
			"no extension",
			Number{AreaCode: "212", Exchange: "555", Line: "0134"},
			"(212) 555-0134",
		},
		{
			"with extension",
			Number{AreaCode: "212", Exchange: "555", Line: "0134", Extension: "45"},
			"(212) 555-0134 ext. 45",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.n.Pretty(); got != c.want {
				t.Fatalf("Pretty() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestE164(t *testing.T) {
	cases := []struct {
		name string
		n    Number
		want string
	}{
		{
			"no extension",
			Number{AreaCode: "212", Exchange: "555", Line: "0134"},
			"+12125550134",
		},
		{
			"with extension",
			Number{AreaCode: "212", Exchange: "555", Line: "0134", Extension: "45"},
			"+12125550134;ext=45",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.n.E164(); got != c.want {
				t.Fatalf("E164() = %q, want %q", got, c.want)
			}
		})
	}
}
