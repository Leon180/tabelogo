package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func main() {

	searchTerm := "鬼金棒"
	searchArea := "東京都"
	linkQueue := []string{}

	c := colly.NewCollector()

	// Find and visit all links on tabelog
	c.OnHTML(".list-rst__rst-name-target", func(e *colly.HTMLElement) {
		fmt.Printf("Link found: %q -> %s\n", e.Text, e.Attr("href"))
		linkQueue = append(linkQueue, e.Attr("href"))

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		fmt.Print(len(linkQueue))
		for i := 0; i < len(linkQueue); i++ {
			fmt.Println(linkQueue[i])
		}
	})

	c.Visit("https://tabelog.com/rstLst/?vs=1&sa=" + searchArea + "&sk=" + searchTerm + "&sw=" + searchTerm)
}
