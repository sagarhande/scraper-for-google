package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var googleDomains = map[string]string{
	"com": "https://www.google.com/search?q=",
	"za":  "https://www.google.co.za/search?q=",
}

type SearchResult struct {
	ResultRank  int
	ResultURL   string
	ResultTitle string
	ResultDesc  string
}

var userAgents = []string{}

func randomUserAgent() string {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func buildGoogleUrls(searchTerm string, languageCode string, countryCode string, pages int, count int) ([]string, error) {
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	if googleBase, found := googleDomains[countryCode]; found {
		for i := 0; i < pages; i++ {
			start := i * count
			scrapeURL := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0", googleBase, searchTerm, count, languageCode, start)

			toScrape = append(toScrape, scrapeURL)
		}
	} else {
		err := fmt.Errorf("Country (%s) is currently not supported", countryCode)
		return nil, err
	}
	return toScrape, nil

}
func scrapeClientRequest(searchURL string, proxyString interface{}) (*http.Response, err) {
	baseClient := getScrapeClient(proxyString)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)

	if req.Response.StatusCode != 200 {
		err := fmt.Errorf("Scraper recieved non 200 status code")
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return res, nil

}

func GoogleScrape(searchTerm string, languageCode string, proxyString interface{}, countryCode string, pages int, count int) ([]SearchResult, err) {
	results := []SearchResult{}
	resultCounter := 0
	googlePages, err := buildGoogleUrls(searchTerm, languageCode, countryCode, pages, count)
	if err != nil {
		return nil, err
	}
	for _, page := range googlePages {
		res, err := scrapeClientRequest(page, proxyString)
		if err != nil {
			return nil, err
		}
		data, err := googleResultParsing(res, resultCounter)
		if err != nil {
			return nil, err
		}
		resultCounter += len(data)

		for _, result := range data {
			results = append(results, result)
		}
		time.Sleep(time.Duration(10) * time.Second)
	}
	return results, nil

}

func main() {
	results, err := GoogleScrape("sagar hande", "en", nil, "com", 1, 30)

	if err == nil {
		for _, res := range results {
			fmt.Println()
		}
	}
}
