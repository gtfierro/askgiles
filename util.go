package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	giles "github.com/gtfierro/giles/archiver"
	"net/http"
	"os"
)

func prettyPrintJSON(v interface{}) {
	if b, err := json.MarshalIndent(v, "", "  "); err != nil {
		fmt.Println("ERROR FORMATTING (%v) %v", err, v)
	} else {
		fmt.Println(string(b))
	}
}

func doQuery(query string) []giles.SmapMessage {
	var buf = bytes.NewBufferString(query)
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
	return data
}

func extractDataNumeric(msg giles.SmapMessage) []float64 {
	var data []float64
	for _, rdg := range msg.Readings {
		if rdg.IsObject() { // skip objects
			continue
		}
		if val, ok := rdg.GetValue().(uint64); ok {
			data = append(data, float64(val))
		} else if val, ok := rdg.GetValue().(float64); ok {
			data = append(data, val)
		} else {
			continue
		}
	}
	return data
}

func extractTime(msg giles.SmapMessage) []float64 {
	var times []float64
	for _, rdg := range msg.Readings {
		val := rdg.GetTime()
		times = append(times, float64(val))
	}
	return times
}
