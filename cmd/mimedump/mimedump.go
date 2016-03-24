package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"

	mime "github.com/jhillyerd/go.enmime"
)

var filename = flag.String("file", "", "path to email with a text/html part")

func main() {
	flag.Parse()

	if filename == nil || strings.TrimSpace(*filename) == "" {
		log.Fatal("--file is required")
	}

	fh, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}

	msg, err := mail.ReadMessage(fh)
	if err != nil {
		log.Fatal(err)
	}

	m, err := mime.ParseMIMEBody(msg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(os.Stderr, "HTML=%d Text=%d\n", len(m.HTML), len(m.Text))
	if len(m.HTML) <= 0 {
		log.Fatalf("No HTML part found in %s\n", *filename)
	}

	w := bufio.NewWriter(os.Stdout)
	n, err := w.WriteString(m.HTML)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "Wrote %d bytes to standard output\n", n)
	w.Flush()
}
