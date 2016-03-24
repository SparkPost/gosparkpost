package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	qp "mime/quotedprintable"
	"os"
	"strings"
)

var breakEvery = flag.Int("break", 0, "insert line breaks every N characters (default: 0)")
var decode = flag.Bool("decode", false, "decode incoming quoted-printable stream")
var input = flag.String("input", "", "filename to read input from (default: stdin)")
var output = flag.String("output", "", "filename to write output to (default: stdout)")

func main() {
	flag.Parse()

	var in io.Reader
	var out io.Writer

	if input == nil || strings.TrimSpace(*input) == "" {
		in = bufio.NewReader(os.Stdin)
	} else {
		fileBytes, err := ioutil.ReadFile(*input)
		if err != nil {
			log.Fatal(err)
		}
		in = bytes.NewReader(fileBytes)
	}

	if decode != nil && *decode == true {
		in = qp.NewReader(in)
	}

	inBytes, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}

	if output == nil || strings.TrimSpace(*output) == "" {
		out = os.Stdout
	} else {
		out, err = os.OpenFile(*output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	if decode == nil || *decode == false {
		out = qp.NewWriter(out)
	}

	n, err := out.Write(inBytes)
	if err != nil {
		log.Fatal(err)
	}
	if n != len(inBytes) {
		log.Fatalf("Partial write (%d != %d)\n", n, len(inBytes))
	}
}
