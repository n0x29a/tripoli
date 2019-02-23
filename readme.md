TRIPOLI
--------------

Library that can do thread pool executor easiest way

Example of usage:

```go
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
	
	"github.com/lootbox/tripoli"
)

var data []interface{}

func worker_fetch_cookies(url string) []*http.Cookie {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	return resp.Cookies()
}

func worker_fetch_body(url string) []byte {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func main() {
	data = append(data,
		"https://google.com",
		"https://google.com",
		"https://yandex.com",
		"https://apple.com",
		"https://www.spacex.com",
		"https://www.twitter.com/",
		"https://www.youtube.com/",
		"http://www.wikidata.org/",
	)

	start := time.Now()
	/**
	Arguments:
	 - worker function
	 - amount of concurrent workers
	 - data for processing
	*/
	r1 := tripoli.Run(worker_fetch_cookies, 9, data)
	log.Printf("🕒 %s", time.Since(start))

	for _, v := range r1 {
		log.Printf("%s", v)
	}

	start = time.Now()
	r2 := tripoli.Run(worker_fetch_body, 9, data)
	log.Printf("🕒 %s", time.Since(start))

	for _, v := range r2 {
		log.Printf("%s", v)
	}
}
```

Roadmap
-------

 - Workers with more than one arguments
 - Data streaming for workers
 - More flexible control flow
 - Tests
 