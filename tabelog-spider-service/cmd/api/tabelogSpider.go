package main

// import (
// 	"strings"
// 	"time"

// 	"github.com/Leon180/tabelogo/helpers"
// 	"github.com/gocolly/colly/v2"
// 	"github.com/gocolly/colly/v2/extensions"
// )

// type Collection struct {
// 	Data any
// }

// type Spider interface {
// 	Collect() error
// 	GetCollections() map[string]Collection
// }

// // link spider will collect all the links in the page which match the selector
// type LinkSpider struct {
// 	Url            string
// 	UrlParams      map[string]string // for search
// 	LinkSelector   string
// 	LinkCollection map[string]Collection
// }

// func NewLinkSpider(ls LinkSpider) Spider {
// 	linkQueue := ls.LinkCollection
// 	if linkQueue == nil {
// 		linkQueue = make(map[string]Collection)
// 	}
// 	return &LinkSpider{
// 		Url:            ls.Url,
// 		UrlParams:      ls.UrlParams,
// 		LinkSelector:   ls.LinkSelector,
// 		LinkCollection: linkQueue,
// 	}
// }

// func (s *LinkSpider) Collect() error {
// 	linkQueue := []string{}
// 	c := colly.NewCollector(
// 		func(collector *colly.Collector) {
// 			extensions.RandomUserAgent(collector)
// 		},
// 	)
// 	c.SetRequestTimeout(2 * time.Second)
// 	c.OnHTML(s.LinkSelector, func(e *colly.HTMLElement) {
// 		linkQueue = append(linkQueue, e.Attr("href"))
// 	})
// 	c.OnScraped(func(r *colly.Response) {
// 		s.LinkCollection = append(s.LinkCollection, Collection{
// 			Name: "links",
// 			Data: linkQueue,
// 		})
// 	})
// 	c.OnError(func(r *colly.Response, err error) {
// 		s.LinkCollection = append(s.LinkCollection, Collection{
// 			Name: "error",
// 			Data: []string{
// 				err.Error(),
// 			},
// 		})
// 	})
// 	return c.Visit(helpers.ParamsCombine(s.Url, s.UrlParams))
// }

// func (s *LinkSpider) GetCollections() []Collection {
// 	return s.LinkCollection
// }

// // selector for content includes container selector and child selector, the selector will select the container first and then select the child in the container with child selector.
// type ContentSelector struct {
// 	ContainerSelector string
// 	ChildSelector     map[string]string
// }

// type TabelogContentSpider struct {
// 	Url               string
// 	ContentSelector   ContentSelector
// 	ContentCollection []Collection
// }

// func NewtabelogContentSpider(tcs TabelogContentSpider) Spider {
// 	return &TabelogContentSpider{
// 		Url:               tcs.Url,
// 		ContentSelector:   tcs.ContentSelector,
// 		ContentCollection: []Collection{},
// 	}
// }

// func (s *TabelogContentSpider) Collect() error {
// 	c := colly.NewCollector(
// 		func(collector *colly.Collector) {
// 			extensions.RandomUserAgent(collector)
// 		},
// 	)
// 	// Make Content Collection
// 	c.OnHTML(s.ContentSelector.ContainerSelector, func(element *colly.HTMLElement) {
// 		data := map[string][]string{}
// 		for key, value := range s.ContentSelector.ChildSelector {
// 			tmp := []string{}
// 			element.ForEach(value, func(_ int, e *colly.HTMLElement) {
// 				tmp = append(tmp, strings.TrimSpace(e.Text))
// 			})
// 			data[key] = tmp
// 		}
// 		s.ContentCollection = append(s.ContentCollection, Collection{
// 			Name: "WebContent",
// 			Data: data,
// 		})
// 	})
// 	c.OnError(func(r *colly.Response, err error) {
// 		s.ContentCollection = append(s.ContentCollection, Collection{
// 			Name: "error",
// 			Data: []string{
// 				err.Error(),
// 			},
// 		})
// 	})
// 	return c.Visit(s.Url)
// }

// func (s *TabelogContentSpider) GetCollections() []Collection {
// 	return s.ContentCollection
// }

// // selector for list content includes parent container selector and content selector, the selector will select the parent container first and then select the content in the parent container with content selector.
// // since for each list item there might be more than 1 content we need, use content selector to select the content in each list item.
// type ListContentSelector struct {
// 	ParentContainerSelector string
// 	ContentSelector         ContentSelector
// }

// // list condition is used to filter the list we want, since there might be more than 1 list using same class name in the page.
// type TabelogListContentSpider struct {
// 	Url                   string
// 	ListContentSelector   ListContentSelector
// 	ListCondition         func(*colly.HTMLElement) bool
// 	ListContentCollection []Collection
// }

// func NewtabelogListContentSpider(tlcs TabelogListContentSpider) Spider {
// 	if tlcs.ListCondition == nil {
// 		tlcs.ListCondition = func(_ *colly.HTMLElement) bool { return true }
// 	}
// 	return &TabelogListContentSpider{
// 		Url:                   tlcs.Url,
// 		ListContentSelector:   tlcs.ListContentSelector,
// 		ListCondition:         tlcs.ListCondition,
// 		ListContentCollection: []Collection{},
// 	}
// }

// func (s *TabelogListContentSpider) Collect() error {
// 	c := colly.NewCollector(
// 		func(collector *colly.Collector) {
// 			extensions.RandomUserAgent(collector)
// 		},
// 	)
// 	// Make List Content Collection
// 	// container1 for ul
// 	c.OnHTML(s.ListContentSelector.ParentContainerSelector, func(element *colly.HTMLElement) {
// 		// container2 for li
// 		if s.ListCondition(element) {
// 			element.ForEach(s.ListContentSelector.ContentSelector.ContainerSelector, func(i int, element *colly.HTMLElement) {
// 				// for all must collect in the li
// 				data := map[string][]string{}
// 				for key, value := range s.ListContentSelector.ContentSelector.ChildSelector {
// 					data[key] = append(data[key], element.ChildText(value))
// 				}
// 				s.ListContentCollection = append(s.ListContentCollection, Collection{
// 					Name: "ListContent",
// 					Data: data,
// 				})
// 			})
// 		}
// 	})
// 	c.OnError(func(r *colly.Response, err error) {
// 		s.ListContentCollection = append(s.ListContentCollection, Collection{
// 			Name: "error",
// 			Data: []string{
// 				err.Error(),
// 			},
// 		})
// 	})
// 	return c.Visit(s.Url)
// }

// func (s *TabelogListContentSpider) GetCollections() []Collection {
// 	return s.ListContentCollection
// }
