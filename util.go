package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	giles "github.com/gtfierro/giles/archiver"
	"log"
	"net/http"
	"os"
	"strconv"
)

func prettyPrintJSON(v interface{}) {
	if b, err := json.MarshalIndent(v, "", "  "); err != nil {
		fmt.Println("ERROR FORMATTING (%v) %v", err, v)
	} else {
		fmt.Println(string(b))
	}
}

func doQuery(host, query string) []giles.SmapMessage {
	var buf = bytes.NewBufferString(query)
	resp, err := http.Post(host, MIME_TEXT, buf)
	if err != nil {
		fmt.Printf("Error running query %v: %v\n", query, err)
		os.Exit(1)
	}

	// expecting json
	var data []giles.SmapMessage
	var decoder = json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	return data
}

func doDistinctQuery(host, query string) []string {
	var buf = bytes.NewBufferString(query)
	resp, err := http.Post(host, MIME_TEXT, buf)
	if err != nil {
		fmt.Printf("Error running query %v: %v\n", query, err)
		os.Exit(1)
	}

	// expecting json
	var data []string
	var decoder = json.NewDecoder(resp.Body)
	decoder.Decode(&data)
	resp.Body.Close()
	return data
}

func extractDataNumeric(msg giles.SmapMessage) []float64 {
	var data []float64
	if len(msg.Readings) == 0 {
		fmt.Printf("No data for %s\n", msg.UUID)
		return data
	}
	for _, rdg := range msg.Readings {
		if rdg == nil || rdg.IsObject() { // skip objects
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
		if rdg == nil {
			continue
		}
		val := rdg.GetTime()
		times = append(times, float64(val))
	}
	return times
}

// encodes a list of floats as a list of strings
func toString(data []float64) []string {
	ret := make([]string, len(data))
	for i, num := range data {
		ret[i] = strconv.FormatFloat(num, 'f', -1, 64)
	}
	return ret
}
