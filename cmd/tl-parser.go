package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/igungor/tl"
)

func main() {
	flag.Usage = func() {
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()

		parser := tl.NewParser(bytes.NewReader(line))
		parser.Trace = true

		parser.ParseProgram()

		if err := parser.Err(); err != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
