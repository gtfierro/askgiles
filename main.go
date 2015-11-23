package main

import (
	"flag"
	"fmt"
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
	flag.StringVar(&host, "h", "http://0.0.0.0:8079/api/query", "Host to direct queries to")
	flag.BoolVar(&formatTime, "t", true, "If true, formats time according to RFC3339. Can specify ANSIC, RFC822, RFC1123, RFC3339 using the -f flag. If false, displays the unix-style timestamp.")
	flag.StringVar(&timeFormat, "f", "RFC3339", "Time format. Can specify ANSIC, RFC822, RFC1123, RFC3339.")
}

func handleQuery(args []string) {
	if len(args) < 1 {
		fmt.Println("No query specified")
		os.Exit(1)
	}

	data := doQuery(args[0])

	prettyPrintJSON(data)
}

func handleVis(args []string) {
	if len(args) < 1 {
		fmt.Println("No query specified")
		os.Exit(1)
	}

	data := doQuery(args[0])

	visType := args[1]
	switch visType {
	case "spark":
		visSpark(data)
	case "plot":
		visPlot(data)
	case "plottime":
		visPlotTime(data)
	default:
		visSpark(data)
	}
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
	case "vis":
		handleVis(args[1:])
	default:
		fmt.Printf("Subcommand %v not found\n", subcommand[0])
		os.Exit(1)
	}
}
