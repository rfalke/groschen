package main

import (
	"net/http"
	"io/ioutil"
	"time"
	"errors"
	. "groschen"
	"gopkg.in/alecthomas/kingpin.v1"
	"fmt"
)

var (
	outputDir = kingpin.Flag("output-dir", "set the output directory.").Default(".").Short('o').String()
	seedUrl   = kingpin.Arg("url", "the start url").Required().String()
)

type MyResponse struct {
	Nested http.Response
	Body   []byte;
}

func downloadWithRetry(url string, whichTry int, totalTries int , prefix string, log LogFunc) (*MyResponse, error) {
	if (whichTry > totalTries) {
		return nil, errors.New("too many retries")
	}
	log(LogStart, prefix, "downloading '%s' %d/%d ...", url, whichTry, totalTries)
	start := time.Now()
	resp, err := request(url, 12)
	duration := time.Now().Sub(start)
	if err == nil && resp.Nested.StatusCode == 200 {
		Bytes := len(resp.Body)
		log(LogEnd, prefix, "got %s in %.1f seconds (%s)", FormatBytes(Bytes), duration.Seconds(), FormatSpeed(len(resp.Body), duration))
		return resp, nil
	} else if err == nil && (resp.Nested.StatusCode == 403 || resp.Nested.StatusCode == 404) {
		log(LogEnd, prefix, "got %d after %.1f seconds", resp.Nested.StatusCode, duration.Seconds())
		return nil, errors.New("permission problem")
	} else if err == nil {
		log(LogEnd, prefix, "got %d after %.1f seconds", resp.Nested.StatusCode, duration.Seconds())
		return downloadWithRetry(url, whichTry+1, totalTries, prefix, log)
	} else {
		log(LogEnd, prefix, "got some other error after %.1f seconds", duration.Seconds())
		return downloadWithRetry(url, whichTry+1, totalTries, prefix, log)
	}
}

func request(url string, timeoutInSec int) (*MyResponse, error) {
	var durationTimeout time.Duration = time.Second * time.Duration(timeoutInSec)
	client := &http.Client{Timeout: durationTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	result := new(MyResponse)
	result.Nested = *resp
	result.Body = body
	return result, nil
}

func handleOneUrl(prefix string, url string, log LogFunc, allNewUrls map[string]bool) {
	TotalTries := 5
	resp, err := downloadWithRetry(url, 1, TotalTries, prefix, log)
	if err != nil {
		log(LogOther, prefix, "    *** failed to download '%s'", url)
	} else {
		fname := WriteResponseToFile(*outputDir, resp.Body, url)
		log(LogOther, prefix, "    saved to %s", fname)
		newUrls := ExtractLinks(url, string(resp.Body))
		AddListToMap(newUrls, allNewUrls)
	}
}

func doOneBatchSequential(todos map[string]bool) map[string]bool{
	var newUrls = make(map[string]bool, 0)
	var counter=0
	for url, _ := range todos {
		handleOneUrl(fmt.Sprintf("  %d/%d", counter, len(todos)), url, SeqLog, newUrls)
		counter++
	}
	return newUrls
}

func driver(todos map[string]bool, done map[string]bool) {
	var batchNumber = -1
	for {
		batchNumber++
		var filteredTodo = NewFromFilter(todos, func(value string)bool {_, ok := done[value]; return !ok})
		if (len(filteredTodo) == 0) {
			break;
		}
		SeqLog(LogOther, "", "driver: start batch %d with %d urls (%d urls already done)",batchNumber, len(filteredTodo), len(done))
		var newUrls = doOneBatchSequential(filteredTodo)
		AddMapToMap(filteredTodo, done)
		todos = newUrls
	}
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	var todo = make(map[string]bool, 0)
	todo[*seedUrl] = true
	driver(todo, make(map[string]bool, 0))
}
