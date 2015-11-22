package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	giles "github.com/gtfierro/giles/archiver"
	"net/http"
	"os"
)

var host string
var formatTime bool
var timeFormat string

var (
	MIME_TEXT = "text/plain"
	MIME_JSON = "application/json"
)

func init() {
	flag.StringVar(&host, "h", "http://localhost", "Host to direct queries to")
	flag.BoolVar(&formatTime, "t", true, "If true, formats time according to RFC3339. Can specify ANSIC, RFC822, RFC1123, RFC3339 using the -f flag. If false, displays the unix-style timestamp.")
	flag.StringVar(&timeFormat, "f", "RFC3339", "Time format. Can specify ANSIC, RFC822, RFC1123, RFC3339.")
}

func prettyPrintJSON(v interface{}) {
	if b, err := json.MarshalIndent(v, "", "  "); err != nil {
		fmt.Println("ERROR FORMATTING (%v) %v", err, v)
	} else {
		fmt.Println(string(b))
	}
}

func handleQuery(args []string) {
	var (
		query string
		buf   *bytes.Buffer
	)

	if len(args) < 1 {
		fmt.Println("No query specified")
		os.Exit(1)
	}

	query = args[0]
	buf = bytes.NewBufferString(query)

	resp, err := http.Post(host, MIME_TEXT, buf)
	if err != nil {
		fmt.Printf("Error running query %v: %v\n", query, err)
		os.Exit(1)
	}

	// expecting json
	var data []giles.SmapMessage
	var decoder = json.NewDecoder(resp.Body)
	decoder.Decode(&data)
	resp.Body.Close()

	prettyPrintJSON(data)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println(`
Need to specify a command! One of:
    query   -   Send a query to a Giles instance
`)
		os.Exit(1)
	}
	subcommand := args[0]
	switch subcommand {
	case "query":
		handleQuery(args[1:])
	default:
		fmt.Printf("Subcommand %v not found\n", subcommand[0])
		os.Exit(1)
	}
}
