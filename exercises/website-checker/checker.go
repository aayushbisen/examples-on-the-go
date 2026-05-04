package websitechecker

import (
	"fmt"
	"net/http"
	"sync"
)

var websiteList = []string{"https://google.com", "https://facebook.com", "https://aayushbisen.dev"}

func checkUrl(url string, wg *sync.WaitGroup, ch chan string) {
	defer wg.Done()

	d, e := http.Get(url)

	if e != nil {
		fmt.Printf("error %s", e)
		return
	}

	if d.StatusCode == 200 {
		ch <- url + " is up"
	} else {
		ch <- url + " is down"
	}
}

func CheckUrls() {
	var wg sync.WaitGroup
	results := make(chan string)

	for _, url := range websiteList {
		wg.Add(1)
		go checkUrl(url, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		fmt.Println(res)
	}

}
