package phonenumber

import "strings"

// Pretty renders the standard US display format, e.g. "(212) 555-0134",
// with an "ext. NNNN" suffix when an extension was captured.
func (n Number) Pretty() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(n.AreaCode)
	b.WriteString(") ")
	b.WriteString(n.Exchange)
	b.WriteString("-")
	b.WriteString(n.Line)
	if n.Extension != "" {
		b.WriteString(" ext. ")
		b.WriteString(n.Extension)
	}
	return b.String()
}

// E164 renders the number per ITU-T E.164: "+1" followed by 10 digits.
// E.164 has no notion of an extension, so RFC 3966's ";ext=" suffix is
// used when one is present -- it's the closest thing to a standard.
func (n Number) E164() string {
	s := "+1" + n.AreaCode + n.Exchange + n.Line
	if n.Extension != "" {
		s += ";ext=" + n.Extension
	}
	return s
}
