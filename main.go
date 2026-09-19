package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/dean-rodriguez/phone-number-format/phonenumber"
)

// renderers maps a --format value to the Number method that produces it.
var renderers = map[string]func(phonenumber.Number) string{
	"standard": phonenumber.Number.Pretty,
	"dot":      phonenumber.Number.Dotted,
	"space":    phonenumber.Number.Spaced,
	"e164":     phonenumber.Number.E164,
}

func main() {
	lenient := flag.Bool("lenient", false, "accept loosely formatted input instead of only canonical NANP layouts")
	format := flag.String("format", "standard", "output style: standard, dot, space, or e164")
	flag.Parse()

	render, ok := renderers[*format]
	if !ok {
		fmt.Fprintf(os.Stderr, "phonenumber: unknown --format %q (want standard, dot, space, or e164)\n", *format)
		os.Exit(1)
	}

	exitCode := 0
	handle := func(raw string) {
		n, err := phonenumber.Parse(raw, *lenient)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
			return
		}
		fmt.Println(render(n))
	}

	if args := flag.Args(); len(args) > 0 {
		for _, raw := range args {
			handle(raw)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			handle(line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "phonenumber: reading stdin:", err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
