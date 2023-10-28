package spiders

import (
	"strings"
	"time"

	"github.com/Leon180/tabelogo/helpers"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

type Collection struct {
	Name string
	Data interface{}
}

type Spider interface {
	Collect() error
	GetCollections() []Collection
}

type LinkSpider struct {
	Url            string
	UrlParams      map[string]string
	LinkSelector   string
	LinkCollection []Collection
}

func NewLinkSpider(ls LinkSpider) Spider {
	linkQueue := ls.LinkCollection
	if linkQueue == nil {
		linkQueue = []Collection{}
	}
	return &LinkSpider{
		Url:            ls.Url,
		UrlParams:      ls.UrlParams,
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
		s.LinkCollection = append(s.LinkCollection, Collection{
			Name: "links",
			Data: linkQueue,
		})
	})
	return c.Visit(helpers.ParamsCombine(s.Url, s.UrlParams))
}

func (s *LinkSpider) GetCollections() []Collection {
	return s.LinkCollection
}

type ContentSelector struct {
	ContainerSelector string
	ChildSelector     map[string]string
}

type TabelogContentSpider struct {
	Url               string
	ContentSelector   ContentSelector
	ContentCollection []Collection
}

func NewtabelogContentSpider(tcs TabelogContentSpider) Spider {
	return &TabelogContentSpider{
		Url:               tcs.Url,
		ContentSelector:   tcs.ContentSelector,
		ContentCollection: []Collection{},
	}
}

func (s *TabelogContentSpider) Collect() error {
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	// Make Content Collection
	c.OnHTML(s.ContentSelector.ContainerSelector, func(element *colly.HTMLElement) {
		data := map[string][]string{}
		for key, value := range s.ContentSelector.ChildSelector {
			tmp := []string{}
			element.ForEach(value, func(_ int, e *colly.HTMLElement) {
				tmp = append(tmp, strings.TrimSpace(e.Text))
			})
			data[key] = tmp
		}
		s.ContentCollection = append(s.ContentCollection, Collection{
			Name: "WebContent",
			Data: data,
		})
	})
	return c.Visit(s.Url)
}

func (s *TabelogContentSpider) GetCollections() []Collection {
	return s.ContentCollection
}

// list
type ListContentSelector struct {
	ParentContainerSelector string
	ContentSelector         ContentSelector
}

type TabelogListContentSpider struct {
	Url                   string
	ListContentSelector   ListContentSelector
	ListCondition         func(*colly.HTMLElement) bool
	ListContentCollection []Collection
}

func NewtabelogListContentSpider(tlcs TabelogListContentSpider) Spider {
	if tlcs.ListCondition == nil {
		tlcs.ListCondition = func(_ *colly.HTMLElement) bool { return true }
	}
	return &TabelogListContentSpider{
		Url:                   tlcs.Url,
		ListContentSelector:   tlcs.ListContentSelector,
		ListCondition:         tlcs.ListCondition,
		ListContentCollection: []Collection{},
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
					data[key] = append(data[key], element.ChildText(value))
				}
				s.ListContentCollection = append(s.ListContentCollection, Collection{
					Name: "ListContent",
					Data: data,
				})
			})
		}
	})
	return c.Visit(s.Url)
}

func (s *TabelogListContentSpider) GetCollections() []Collection {
	return s.ListContentCollection
}
