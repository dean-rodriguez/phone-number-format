// Package phonenumber parses and formats North American Numbering Plan
// (NANP) phone numbers: the +1 country code territories (US, Canada, and
// a handful of Caribbean nations). It does not attempt to handle other
// numbering plans.
package phonenumber

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Number is a parsed NANP number, split into its NXX-NXX-XXXX components.
// AreaCode, Exchange, and Line are always exactly 3, 3, and 4 digits.
type Number struct {
	AreaCode  string
	Exchange  string
	Line      string
	Extension string // digits only, empty if none was present
}

// ErrEmpty is returned when Parse is given a blank or whitespace-only string.
var ErrEmpty = errors.New("phonenumber: empty input")

// ParseError explains why a given input could not be parsed.
type ParseError struct {
	Input  string
	Reason string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("phonenumber: cannot parse %q: %s", e.Input, e.Reason)
}

// Strict mode only accepts a small, unambiguous set of canonical layouts.
// Anything else -- missing punctuation, unusual separators, trailing junk --
// is rejected rather than guessed at.
var strictPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^\+1(\d{3})(\d{3})(\d{4})$`),
	regexp.MustCompile(`^\+1-(\d{3})-(\d{3})-(\d{4})$`),
	regexp.MustCompile(`^\((\d{3})\) (\d{3})-(\d{4})$`),
	regexp.MustCompile(`^(\d{3})-(\d{3})-(\d{4})$`),
}

var extensionPattern = regexp.MustCompile(`(?i)^(.*?)\s*(?:ext\.?|x|#)\s*(\d{1,6})$`)

// Parse turns raw into a Number. With lenient=false, raw must already be in
// one of the canonical layouts listed above. With lenient=true, punctuation
// and whitespace are stripped and only the digit count and NANP digit rules
// are enforced -- see the package README for the full comparison.
func Parse(raw string, lenient bool) (Number, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return Number{}, ErrEmpty
	}
	if lenient {
		return parseLenient(trimmed)
	}
	return parseStrict(trimmed)
}

func parseStrict(s string) (Number, error) {
	for _, pat := range strictPatterns {
		m := pat.FindStringSubmatch(s)
		if m == nil {
			continue
		}
		n := Number{AreaCode: m[1], Exchange: m[2], Line: m[3]}
		if err := validate(n, false); err != nil {
			return Number{}, &ParseError{Input: s, Reason: err.Error()}
		}
		return n, nil
	}
	return Number{}, &ParseError{
		Input:  s,
		Reason: "does not match a recognized strict format (try +1XXXXXXXXXX, +1-XXX-XXX-XXXX, (XXX) XXX-XXXX, or XXX-XXX-XXXX)",
	}
}

func parseLenient(s string) (Number, error) {
	body := s
	extension := ""
	if m := extensionPattern.FindStringSubmatch(s); m != nil {
		body = strings.TrimSpace(m[1])
		extension = m[2]
	}

	digits := stripNonDigits(vanityToDigits(body))
	switch {
	case len(digits) == 11 && strings.HasPrefix(digits, "1"):
		digits = digits[1:]
	case len(digits) == 10:
		// already bare
	default:
		return Number{}, &ParseError{
			Input:  s,
			Reason: fmt.Sprintf("expected 10 digits (or 11 with a leading country code 1), got %d", len(digits)),
		}
	}

	n := Number{
		AreaCode:  digits[0:3],
		Exchange:  digits[3:6],
		Line:      digits[6:10],
		Extension: extension,
	}
	if err := validate(n, true); err != nil {
		return Number{}, &ParseError{Input: s, Reason: err.Error()}
	}
	return n, nil
}

// vanityLetters maps each letter to the digit it occupies on a standard
// telephone keypad, e.g. 1-800-FLOWERS dials the same as 1-800-3569377.
var vanityLetters = map[rune]byte{
	'A': '2', 'B': '2', 'C': '2',
	'D': '3', 'E': '3', 'F': '3',
	'G': '4', 'H': '4', 'I': '4',
	'J': '5', 'K': '5', 'L': '5',
	'M': '6', 'N': '6', 'O': '6',
	'P': '7', 'Q': '7', 'R': '7', 'S': '7',
	'T': '8', 'U': '8', 'V': '8',
	'W': '9', 'X': '9', 'Y': '9', 'Z': '9',
}

// vanityToDigits replaces letters with the digit they map to on a phone
// keypad, leaving digits and punctuation untouched. It's applied before
// stripNonDigits so a number like "1-800-CALL-NOW" parses the same way
// its digit equivalent would.
func vanityToDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if d, ok := vanityLetters[unicode.ToUpper(r)]; ok {
			b.WriteByte(d)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func stripNonDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// validate applies NANP's digit rules. Area codes and exchanges can't start
// with 0 or 1 (those lead digits are reserved for operator/international
// routing). N11 codes (211, 411, 611, ...) are reserved service numbers; in
// lenient mode that check is skipped, since real-world input -- old business
// listings, numbers typed by hand -- sometimes hits a false positive there
// and rejecting it outright isn't worth the friction.
func validate(n Number, lenient bool) error {
	if n.AreaCode[0] < '2' || n.AreaCode[0] > '9' {
		return fmt.Errorf("area code %s cannot start with %c", n.AreaCode, n.AreaCode[0])
	}
	if n.Exchange[0] < '2' || n.Exchange[0] > '9' {
		return fmt.Errorf("exchange %s cannot start with %c", n.Exchange, n.Exchange[0])
	}
	if !lenient {
		if n.AreaCode[1] == '1' && n.AreaCode[2] == '1' {
			return fmt.Errorf("area code %s is a reserved N11 service code", n.AreaCode)
		}
		if n.Exchange[1] == '1' && n.Exchange[2] == '1' {
			return fmt.Errorf("exchange %s is a reserved N11 service code", n.Exchange)
		}
	}
	return nil
}
