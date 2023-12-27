package main

import (
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

type Spider interface {
	Collect() error
	GetCollections() any
}

// link spider will collect all the links in the page which match the selector
// LinkCollection:
//
//	{
//		["link1", "link2", ...],
//	}
type LinkSpider struct {
	Url            string
	UrlParams      map[string]string // for search
	CombinedUrl    string
	LinkSelector   string
	LinkCollection []string
}

func NewLinkSpider(ls LinkSpider) *LinkSpider {
	linkQueue := ls.LinkCollection
	combinedUrl := ParamsCombine(ls.Url, ls.UrlParams)

	if linkQueue == nil {
		linkQueue = []string{}
	}
	return &LinkSpider{
		Url:            ls.Url,
		UrlParams:      ls.UrlParams,
		CombinedUrl:    combinedUrl,
		LinkSelector:   ls.LinkSelector,
		LinkCollection: linkQueue,
	}
}

func (s *LinkSpider) Collect() error {
	linkQueue := []string{}
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	c.SetRequestTimeout(2 * time.Second)
	c.OnHTML(s.LinkSelector, func(e *colly.HTMLElement) {
		linkQueue = append(linkQueue, e.Attr("href"))
	})
	c.OnScraped(func(r *colly.Response) {
		s.LinkCollection = append(s.LinkCollection, linkQueue...)
	})
	return c.Visit(s.CombinedUrl)
}

func (s *LinkSpider) GetCollections() []string {
	return s.LinkCollection
}

// selector for content includes container selector and child selector, the selector will select the container first and then select the child in the container with child selector.
type ContentSelector struct {
	ContainerSelector string
	ChildSelector     map[string]string
}

// content spider will collect all the content in the page which match the selector
// ContentCollection:
//
//	{
//		"content1": []data{},
//		"content2": []data{},
//		...
//	}
type TabelogContentSpider struct {
	Url               string
	ContentSelector   ContentSelector
	ContentCollection map[string][]string
}

func NewtabelogContentSpider(tcs TabelogContentSpider) *TabelogContentSpider {
	return &TabelogContentSpider{
		Url:               tcs.Url,
		ContentSelector:   tcs.ContentSelector,
		ContentCollection: make(map[string][]string),
	}
}

func (s *TabelogContentSpider) Collect() error {
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	c.OnHTML(s.ContentSelector.ContainerSelector, func(element *colly.HTMLElement) {
		data := map[string][]string{}
		for key, value := range s.ContentSelector.ChildSelector {
			tmp := []string{}
			element.ForEach(value, func(_ int, e *colly.HTMLElement) {
				tmp = append(tmp, strings.TrimSpace(e.Text))
			})
			data[key] = tmp
		}
		s.ContentCollection = data
	})
	return c.Visit(s.Url)
}

func (s *TabelogContentSpider) GetCollections() map[string][]string {
	return s.ContentCollection
}

// selector for list content includes parent container selector and content selector, the selector will select the parent container first and then select the content in the parent container with content selector.
// since for each list item there might be more than 1 content we need, use content selector to select the content in each list item.
type ListContentSelector struct {
	ParentContainerSelector string
	ContentSelector         ContentSelector
}

// list condition is used to filter the list we want, since there might be more than 1 list using same class name in the page.
// ListContentCollection:
//
//	[]data
type TabelogListContentSpider struct {
	Url                   string
	ListContentSelector   ListContentSelector
	ListCondition         func(*colly.HTMLElement) bool
	CollectElement        *colly.HTMLElement
	ListContentCollection []map[string][]string
}

func NewtabelogListContentSpider(tlcs TabelogListContentSpider) *TabelogListContentSpider {
	if tlcs.ListCondition == nil {
		tlcs.ListCondition = func(e *colly.HTMLElement) bool { return true }
	}
	return &TabelogListContentSpider{
		Url:                   tlcs.Url,
		ListContentSelector:   tlcs.ListContentSelector,
		ListCondition:         tlcs.ListCondition,
		ListContentCollection: []map[string][]string{},
	}
}

func (s *TabelogListContentSpider) Collect() error {
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	// Make List Content Collection
	// container1 for ul
	c.OnHTML(s.ListContentSelector.ParentContainerSelector, func(element *colly.HTMLElement) {
		// container2 for li
		if s.ListCondition(element) {
			element.ForEach(s.ListContentSelector.ContentSelector.ContainerSelector, func(i int, element *colly.HTMLElement) {
				// for all must collect in the li
				data := map[string][]string{}
				for key, value := range s.ListContentSelector.ContentSelector.ChildSelector {
					if key == "img" {
						data[key] = append(data[key], element.ChildAttr(value, "src"))
					} else {
						data[key] = append(data[key], element.ChildText(value))
					}
				}
				s.ListContentCollection = append(s.ListContentCollection, data)
			})
		}
	})
	return c.Visit(s.Url)
}

func (s *TabelogListContentSpider) GetCollections() []map[string][]string {
	return s.ListContentCollection
}
