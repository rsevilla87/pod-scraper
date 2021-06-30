package scraper

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var httpClient http.Client = http.Client{}

type PodScraper struct {
	wg     *sync.WaitGroup
	code   int
	Failed int
}

func NewScraper(wg *sync.WaitGroup, code int, timeout time.Duration) *PodScraper {
	httpClient.Timeout = timeout
	return &PodScraper{wg, code, 0}
}

func (ps *PodScraper) Scrape(target string) error {
	defer ps.wg.Done()
	var err error
	resp, err := httpClient.Get(target)
	if err != nil {
		fmt.Printf("Error scraping %v: %v\n", target, err)
		ps.Failed = 1
		return err
	}
	if resp.StatusCode != ps.code {
		fmt.Printf("Expected status code %d different than %d from %v\n", ps.code, resp.StatusCode, target)
		ps.Failed = 1
		return err
	}
	fmt.Println("Scraped", target)
	return nil
}
