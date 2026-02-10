package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	flag.Parse()

	decoder := json.NewDecoder(os.Stdin)
	for {
		var entry map[string]any
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("error decoding: %v", err)
			continue
		}

		line := format(entry)
		fmt.Println(line)
	}
}
