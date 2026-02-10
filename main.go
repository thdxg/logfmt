package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/thdxg/logfmt/pkg/cli"
	"github.com/thdxg/logfmt/pkg/formatter"
	"github.com/thdxg/logfmt/pkg/parser"
)

func main() {
	cfg := cli.ParseFlags()

	if !cli.HasStdin() {
		cli.Usage()
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		entry, err := parser.Parse(line)
		if err != nil {
			// Print raw text if unparsable
			fmt.Println(string(line))
			continue
		}
		fmt.Println(formatter.Format(cfg, entry))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "reading standard input: %v\n", err)
	}
}
