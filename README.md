# phone-number-format

A parser and pretty printer for NANP phone numbers (the +1 numbering plan
covering the US, Canada, and a handful of Caribbean territories). Feed it
whatever format a human typed a number in, get back something normalized.

The problem this is solving: phone numbers show up in scraped data, CSV
exports, pasted from emails, typed into forms, in every layout imaginable --
`212-555-0134`, `(212) 555-0134`, `+12125550134`, `212.555.0134 x42`. Most of
the time you want one of two things: strict validation (reject anything that
isn't unambiguously a real NANP number) or best-effort normalization (accept
messy input and pull a number out of it). This tool does both, and makes you
pick.

## Usage

```
$ go run . "(212) 555-0134"
(212) 555-0134

$ go run . --format e164 "(212) 555-0134"
+12125550134

$ go run . --format dot "(212) 555-0134"
212.555.0134

$ go run . "212.555.0134"
phonenumber: cannot parse "212.555.0134": does not match a recognized strict format (try +1XXXXXXXXXX, +1-XXX-XXX-XXXX, (XXX) XXX-XXXX, or XXX-XXX-XXXX)

$ go run . --lenient "212.555.0134 ext 42"
(212) 555-0134 ext. 42
```

It also reads from stdin, one number per line, when no arguments are given:

```
$ printf '212-555-0134\n+1-415-555-0100\n' | go run .
(212) 555-0134
(415) 555-0100
```

Exit status is 1 if any line failed to parse; failures go to stderr, one
line's failure doesn't stop the rest from being processed.

## Strict vs. lenient

By default, input has to already be in one of four canonical layouts:

- `+1XXXXXXXXXX`
- `+1-XXX-XXX-XXXX`
- `(XXX) XXX-XXXX`
- `XXX-XXX-XXXX`

Anything else is rejected, including things that are "obviously" the same
number with different punctuation. That's deliberate: strict mode is for
places where you control the input format (a form field with a mask, a
canonicalized database column) and want a parse failure to mean something
went genuinely wrong upstream, not "someone used dots instead of dashes."

Pass `--lenient` to relax that: punctuation and whitespace are stripped down
to digits, a leading country code digit `1` is accepted and dropped, and an
`ext` / `x` / `#` suffix is recognized as an extension. Lenient mode also
skips the reserved-N11-code check (211, 411, 611, ...) since it's a common
source of false positives on real-world data and rejecting a number outright
over it isn't worth the friction. What lenient mode does *not* do is guess at
missing digits -- a 7-digit local number with no area code still fails,
because there's no correct default to assume.

Both modes still enforce the core NANP structural rule: area codes and
exchanges can't start with 0 or 1.

## Output styles

`--format` picks the output layout: `standard` (default, `(212) 555-0134`),
`dot` (`212.555.0134`), `space` (`212 555 0134`), or `e164` (`+12125550134`).
An unrecognized value is a usage error, not a silent fallback to standard.

## Layout

- `phonenumber/parse.go` -- `Parse(raw string, lenient bool) (Number, error)`
- `phonenumber/format.go` -- `Number.Pretty()`, `Number.Dotted()`,
  `Number.Spaced()`, and `Number.E164()`
- `main.go` -- CLI wrapper

## Status

Early. NANP only -- no support yet for other countries' numbering plans.
See the roadmap in commit history for what's planned next.
