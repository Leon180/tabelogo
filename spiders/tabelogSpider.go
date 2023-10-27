package spiders

import (
	"fmt"

	"github.com/gocolly/colly"
)

type linkSpider struct {
	url               string
	urlParams         map[string]string
	linkFoundSelector string
	linkQueue         []string
}

func (s *linkSpider) Run() {
	c := colly.NewCollector()
	c.OnHTML(s.linkFoundSelector, func(e *colly.HTMLElement) {
		fmt.Printf("Link found: %q -> %s\n", e.Text, e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(paramsCombine(s.url, s.urlParams))
}

type TabelogSpider struct {
	linkSpider    *Spider
	contentSpider *Spider
}

func NewSpider(searchTerm string, searchArea string) *Spider {
	return &Spider{
		searchTerm: searchTerm,
		searchArea: searchArea,
	}
}

func (s *Spider) Run() {
	c := colly.NewCollector()

	// Find and visit all links on tabelog
	c.OnHTML(".list-rst__rst-name-target", func(e *colly.HTMLElement) {
		fmt.Printf("Link found: %q -> %s\n", e.Text, e.Attr("href"))
		s.linkQueue = append(s.linkQueue, e.Attr("href"))

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		fmt.Print(len(s.linkQueue))
		for i := 0; i < len(s.linkQueue); i++ {
			fmt.Println(s.linkQueue[i])
		}
	})

	c.Visit(paramsConbine(s.url, s.urlParams))
}
