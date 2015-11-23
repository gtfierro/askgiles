package main

import (
	"flag"
	"os"
)

var (
	MIME_TEXT = "text/plain"
	MIME_JSON = "application/json"
)

func handleQuery(host, query string) {
	data := doQuery(host, query)
	prettyPrintJSON(data)
}

func handleVis(host, query, viztype string) {
	data := doQuery(host, query)

	switch viztype {
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
	var host string
	var formatTime bool
	var timeFormat string
	var query string

	queryCommand := flag.NewFlagSet("query", flag.ExitOnError)
	queryCommand.StringVar(&host, "h", "http://0.0.0.0:8079/api/query", "Host to direct queries to")
	queryCommand.BoolVar(&formatTime, "t", true, "If true, formats time according to RFC3339. Can specify ANSIC, RFC822, RFC1123, RFC3339 using the -f flag. If false, displays the unix-style timestamp.")
	queryCommand.StringVar(&timeFormat, "f", "RFC3339", "Time format. Can specify ANSIC, RFC822, RFC1123, RFC3339.")

	vizCommand := flag.NewFlagSet("viz", flag.ExitOnError)
	vizCommand.StringVar(&host, "h", "http://0.0.0.0:8079/api/query", "Host to direct queries to")

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "query":
		queryCommand.Parse(os.Args[2:])
		query = os.Args[len(os.Args)-1]
		handleQuery(host, query)
	case "viz":
		query = os.Args[len(os.Args)-1]
		vizCommand.Parse(os.Args[3:])
		handleVis(host, query, os.Args[2])
	}

	if queryCommand.Parsed() {
	}
}
