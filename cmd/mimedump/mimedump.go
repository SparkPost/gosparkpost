package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	mime "github.com/jhillyerd/enmime"
)

func main() {
	var filename = flag.String("file", "", "path to email with a text/html part")

	flag.Parse()

	if filename == nil || strings.TrimSpace(*filename) == "" {
		log.Fatal("--file is required")
	}

	fh, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	e, err := mime.ReadEnvelope(fh)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "HTML=%d Text=%d\n", len(e.HTML), len(e.Text))
	if len(e.HTML) <= 0 {
		log.Fatalf("No HTML part found in %s\n", *filename)
	}

	w := bufio.NewWriter(os.Stdout)
	n, err := w.WriteString(e.HTML)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "Wrote %d bytes to standard output\n", n)
	w.Flush()
}
