package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type requestData struct {
	responseCode int
	duration     time.Duration
	success      bool
}

var requestDataCollection []requestData
var (
	sec int
	url string
	rps int
)
var duration time.Duration

func main() {
	initFlags()

	var wg sync.WaitGroup
	startTime := time.Now()

	for j := 0; j < sec; j++ {
		for i := 0; i < rps; i++ {
			wg.Add(1)
			go get(&wg)
		}
		time.Sleep(time.Second)
	}
	wg.Wait()
	duration = time.Now().Sub(startTime)
	printResult()
}

func printResult() {
	fmt.Println("RESULTS!")

	rcMap := make(map[int]int)
	cSuccess, cError := 0, 0
	var respTimeMs int64

	for _, requestData := range requestDataCollection {

		if requestData.responseCode > 0 {
			_, ok := rcMap[requestData.responseCode]
			if !ok {
				rcMap[requestData.responseCode] = 0
			}
			rcMap[requestData.responseCode] += 1
		}

		if requestData.success {
			cSuccess += 1
		} else {
			cError += 1
		}

		respTimeMs += requestData.duration.Milliseconds()
	}

	respTimeMsMean := respTimeMs / int64(cSuccess)

	fmt.Println("Total request count:", cError+cSuccess)
	fmt.Println("Successful request count:", cSuccess)
	fmt.Println("Unsuccessful request count:", cError)
	fmt.Println("Success Rate: %", 100*cSuccess/(cError+cSuccess))
	fmt.Println("Response time (mean):", respTimeMsMean, "ms")
	fmt.Println("Response Codes:", rcMap)
	fmt.Println("Total Duration:", duration)
}

func get(wg *sync.WaitGroup) {
	var rd requestData
	start := time.Now()
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(url)

	rd.success = err == nil
	if resp != nil {
		rd.responseCode = resp.StatusCode
		rd.duration = time.Now().Sub(start)
	}

	requestDataCollection = append(requestDataCollection, rd)
	wg.Done()
}
func initFlags() {
	flag.IntVar(&rps, "rps", 10, "Requests per second")
	flag.IntVar(&sec, "sec", 5, "Running duration")
	flag.StringVar(&url, "url", "http://localhost", "URL to test")
	flag.Parse()
}
