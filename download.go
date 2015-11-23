package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

var wg sync.WaitGroup
var count int64
var total int64

const NUM_WORKERS = 10

//doDownload(host, query, destination beginYear, endYear, dataLimit)
func doDownload(host, where, destination, timeunit string, beginYear, endYear, dataLimit int) {
	// make directory if no exist
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		fmt.Printf("%s not found. Creating...\n", destination)
		os.Mkdir(destination, os.ModeDir|os.ModePerm)
	}

	uuidQuery := fmt.Sprintf("select distinct uuid where %s;", where)
	uuids := doDistinctQuery(host, uuidQuery)

	total = int64(len(uuids))
	count = 0

	uuidQueue := make(chan string, NUM_WORKERS)

	wg.Add(len(uuids))
	for w := 0; w < NUM_WORKERS; w++ {
		go runWorker(uuidQueue, host, destination, timeunit, beginYear, endYear, dataLimit)
	}

	for _, uuid := range uuids {
		uuidQueue <- uuid
	}
	close(uuidQueue)

	// get metadata
	mdQuery := fmt.Sprintf("select * where %s;", where)
	metadata := doQuery(host, mdQuery)
	mdFilename := fmt.Sprintf("%s/metadata.json", destination)
	file, err := os.Create(mdFilename)
	if err != nil {
		panic(err)
	}
	encoder := json.NewEncoder(file)
	encoder.Encode(metadata)
	if err := file.Sync(); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
	fmt.Println("Saved metadata!")
	wg.Wait()
}

func runWorker(uuids chan string, host, destination, timeunit string, beginYear, endYear, dataLimit int) {
	for uuid := range uuids {
		dlUUID(host, uuid, destination, timeunit, beginYear, endYear, dataLimit)
		wg.Done()
	}
}

func dlUUID(host, uuid, destination, timeunit string, beginYear, endYear, dataLimit int) {
	// open file
	filename := fmt.Sprintf("%s/%s.csv", destination, uuid)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(file)

	var query string
	if timeunit == "_" {
		query = fmt.Sprintf("select data in ('1/1/%d', '1/1/%d') limit %d where uuid = '%s';", beginYear, endYear, dataLimit, uuid)
	} else {
		query = fmt.Sprintf("select data in ('1/1/%d', '1/1/%d') limit %d as %s where uuid = '%s';", beginYear, endYear, dataLimit, timeunit, uuid)
	}
	msgs := doQuery(host, query)

	for _, msg := range msgs {
		numbers := toString(extractDataNumeric(msg))
		times := toString(extractTime(msg))
		log.Printf("%d/%d: %s -- %d points\n", atomic.AddInt64(&count, 1), total, msg.UUID, len(numbers))
		for i := 0; i < len(numbers); i++ {
			writer.Write([]string{times[i], numbers[i]})
		}
	}
	writer.Flush()
	if err := file.Sync(); err != nil {
		panic(err)
	}
	if err := file.Close(); err != nil {
		panic(err)
	}
}
