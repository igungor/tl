package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
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

		tlscanner := tl.NewScanner(bytes.NewReader(line))
		for {
			token := tlscanner.Scan()
			fmt.Println(token)

			if token.Token == tl.ItemEOF {
				break
			}

			if err := tlscanner.Err(); err != nil {
				fmt.Println("err: ", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
