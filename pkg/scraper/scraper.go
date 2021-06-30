package scraper

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var httpClient http.Client = http.Client{
	Timeout: 10 * time.Second,
}

func Scrape(target string, code int, wg *sync.WaitGroup, failed *int) {
	defer wg.Done()
	resp, err := httpClient.Get(target)
	if err != nil {
		fmt.Printf("Error scraping %v: %v", target, err)
		*failed = 1
	}
	if resp.StatusCode != code {
		fmt.Printf("Expected status code %d different than %d", code, resp.StatusCode)
		*failed = 1
	}
}
