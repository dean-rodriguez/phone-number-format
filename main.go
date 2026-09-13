package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/dean-rodriguez/phone-number-format/phonenumber"
)

func main() {
	lenient := flag.Bool("lenient", false, "accept loosely formatted input instead of only canonical NANP layouts")
	e164 := flag.Bool("e164", false, "print output as +1XXXXXXXXXX instead of (XXX) XXX-XXXX")
	flag.Parse()

	exitCode := 0
	handle := func(raw string) {
		n, err := phonenumber.Parse(raw, *lenient)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
			return
		}
		if *e164 {
			fmt.Println(n.E164())
		} else {
			fmt.Println(n.Pretty())
		}
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
