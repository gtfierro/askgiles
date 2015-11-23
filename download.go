package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

const NUM_WORKERS = 10

//doDownload(host, query, destination beginYear, endYear, dataLimit)
func doDownload(host, where, destination, timeunit string, beginYear, endYear, dataLimit int) {
	uuidQuery := fmt.Sprintf("select distinct uuid where %s;", where)
	uuids := doDistinctQuery(host, uuidQuery)

	uuidQueue := make(chan string, NUM_WORKERS)

	wg.Add(len(uuids))
	for w := 0; w < NUM_WORKERS; w++ {
		go runWorker(uuidQueue, host, destination, timeunit, beginYear, endYear, dataLimit)
	}

	for _, uuid := range uuids {
		uuidQueue <- uuid
	}
	close(uuidQueue)
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

	query := fmt.Sprintf("select data in ('1/1/%d', '1/1/%d') limit %d as %s where uuid = '%s';", endYear, beginYear, dataLimit, timeunit, uuid)
	msgs := doQuery(host, query)

	for _, msg := range msgs {
		numbers := toString(extractDataNumeric(msg))
		times := toString(extractTime(msg))
		log.Printf("%s -- %d points\n", msg.UUID, len(numbers))
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
