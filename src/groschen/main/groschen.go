package main

import (
	"net/http"
	"io/ioutil"
	"time"
	"errors"
	. "groschen"
	"gopkg.in/alecthomas/kingpin.v1"
	"fmt"
	"sync"
	"runtime"
	"strings"
	"regexp"
)

var (
	outputDir           = kingpin.Flag("output-dir", "set the output directory.").Default(".").Short('o').String()
	seedUrl             = kingpin.Arg("url", "the start url").Required().String()
	parallelConnections = kingpin.Flag("connections", "number of parallel connections").Default("20").Short('c').Int()
	restrictSameDir     = kingpin.Flag("r1", "limit recursive download to the sub directory of the initial URL").Bool()
	restrictSameHost    = kingpin.Flag("r2", "limit recursive download to host of the initial URL").Bool()
	restrictNone        = kingpin.Flag("r3", "do not limit the recursive download").Bool()
	includePattern      = kingpin.Flag("include", "only visit urls which match PATTERNS. PATTERNS consists of a delimeter char (which is the comma) and a list of positive and negative patterns separated by the delimter char. The url is matched against each pattern and the first match decided (i.e. the url is accepted or rejected). Empty parts mean no match.").Short('i').String()
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

var allNewUrlsLock sync.Mutex

func handleOneUrl(prefix string, url string, log LogFunc, allNewUrls map[string]bool) {
	TotalTries := 5
	resp, err := downloadWithRetry(url, 1, TotalTries, prefix, log)
	if err != nil {
		log(LogOther, prefix, "    *** failed to download '%s'", url)
	} else {
		fname := WriteResponseToFile(*outputDir, resp.Body, url)
		log(LogOther, prefix, "    saved to %s", fname)
		newUrls := ExtractLinks(url, string(resp.Body))
		var filteredNewUrls = filterUrls(newUrls)

		allNewUrlsLock.Lock()
		defer allNewUrlsLock.Unlock()
		AddListToMap(filteredNewUrls, allNewUrls)
	}
}

func filterUrls(urls []string) []string {
	result := make([]string, 0)
	for _, link := range urls {
		if is_acceptable_by_restrictions(*seedUrl, link) && is_acceptable_by_include_pattern(link) {
			result = append(result, link)
		}
	}
	return result
}

func is_acceptable_by_include_pattern(link string) bool {
	if *includePattern != "" {
		var sep = (*includePattern)[0:1]
		var parts = strings.Split((*includePattern)[1:], regexp.QuoteMeta(sep))
		var accept = true
		for _, part := range parts {
			if part != "" {
				if regexp.MustCompile(part).FindStringIndex(link) != nil {
					return accept
				}
			}
			accept = !accept
		}
		return accept
	} else {
		return true
	}
}

func is_acceptable_by_restrictions(startlink string, link string) bool {
	if *restrictNone {
		return true
	} else if *restrictSameHost {
		return IsSameHost(link, startlink)
	} else if *restrictSameDir {
		return IsInSubdir(link, startlink)
	} else {
		// no recursion
		return false
	}
}

func doOneBatchSequential(todos map[string]bool) map[string]bool {
	var newUrls = make(map[string]bool, 0)
	var counter = 0
	for url, _ := range todos {
		handleOneUrl(fmt.Sprintf("  %d/%d", counter, len(todos)), url, SeqLog, newUrls)
		counter++
	}
	return newUrls
}

func releaseSlot(finishedChan chan bool) {
	finishedChan <- true
}

func doOneBatchParallel(todos map[string]bool, parallelDownloads int) map[string]bool {
	var finishedChan = make(chan bool)
	var openSlots = parallelDownloads
	var newUrls = make(map[string]bool, 0)
	var todosAsSlice = SliceFromMapKeys(todos)
	var nextIndex = 0
	var finished = false
	for {
		for {
			if (openSlots > 0 && !finished) {
				go func(index int, url string) {
					defer releaseSlot(finishedChan)
					handleOneUrl(fmt.Sprintf("  %d/%d", index, len(todos)), url, SeqLog, newUrls)
				}(nextIndex, todosAsSlice[nextIndex]);
				nextIndex++
				openSlots--
				finished = nextIndex >= len(todosAsSlice)
			} else {
				break
			}
		}
		<-finishedChan
		openSlots++
		if openSlots == parallelDownloads && finished {
			break
		}
	}
	return newUrls
}

func driver(todos map[string]bool, done map[string]bool) {
	var batchNumber = -1
	for {
		batchNumber++
		var filteredTodo = NewFromFilter(todos, func(value string) bool {_, ok := done[value]; return !ok})
		if (len(filteredTodo) == 0) {
			break;
		}
		SeqLog(LogOther, "", "driver: start batch %d with %d urls (%d urls already done)", batchNumber, len(filteredTodo), len(done))
		var newUrls map[string]bool
		if *parallelConnections <= 1 {
			newUrls = doOneBatchSequential(filteredTodo)
		} else {
			newUrls = doOneBatchParallel(filteredTodo, *parallelConnections)
		}
		AddMapToMap(filteredTodo, done)
		todos = newUrls
	}
}

func validateOptions() bool {
	var numRestrictions = 0

	if *restrictSameDir {
		numRestrictions++
	}
	if *restrictSameHost {
		numRestrictions++
	}
	if *restrictNone {
		numRestrictions++
	}
	if numRestrictions > 1 {
		fmt.Println("Either give --r1, --r2 or --r3.")
		return false
	}
	return true
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() + 2)
	kingpin.Version("0.0.1")
	kingpin.Parse()

	if validateOptions() {
		var todo = make(map[string]bool, 0)
		todo[*seedUrl] = true
		driver(todo, make(map[string]bool, 0))
	}
}
