package main

import (
	"fmt"

	"github.com/Leon180/tabelogo/spiders"
	"github.com/gocolly/colly/v2"
)

func main() {
	var err error
	maxLinks := 3
	searchTerm := "鬼金棒"
	searchArea := "東京都"

	linkSpideRequest := spiders.LinkSpider{
		Url: "https://tabelog.com/rstLst/",
		UrlParams: map[string]string{
			"vs": "1",
			"sa": searchArea,
			"sk": searchTerm,
			"sw": searchTerm,
		},
		LinkSelector: ".list-rst__rst-name-target",
	}
	linkSpider := spiders.NewLinkSpider(linkSpideRequest)
	err = linkSpider.Collect()
	if err != nil {
		fmt.Println(err)
	}
	LinkCollection := linkSpider.GetCollections()
	links := LinkCollection[0].Data.([]string)
	if len(links) == 0 {
		fmt.Println("No links found")
		return
	}
	if len(links) > maxLinks {
		links = links[:maxLinks:maxLinks]
	}
	for _, link := range links {
		fmt.Println(link)
	}
	// table photo link: link+"/table/"
	// menu link: link+"/dtlmenu/"
	//	drink menu link: link+"/dtlmenu/drink/"
	// comments link: link+"/dtlrvwlst/"
	// rating distribute link: link+"/dtlratings/"

	tbcRequest := spiders.TabelogContentSpider{
		Url: links[0],
		ContentSelector: spiders.ContentSelector{
			ContainerSelector: "#container",
			ChildSelector: map[string]string{
				"name":        "h2.display-name",
				"rating":      ".rdheader-rating__score b.c-rating__val",
				"ratingCount": ".rdheader-rating__review-target .num",
				"bookmarks":   ".rdheader-rating__hozon-target .num",
			},
		},
	}
	tbSpider := spiders.NewtabelogContentSpider(tbcRequest)
	err = tbSpider.Collect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tbSpider.GetCollections())

	tbcListRequest := spiders.TabelogListContentSpider{
		Url: links[0] + "dtlratings/",
		ListCondition: func(e *colly.HTMLElement) bool {
			return e.ChildText(".ratings-contents__title") == "評価分布"
		},
		ListContentSelector: spiders.ListContentSelector{
			ParentContainerSelector: ".ratings-contents__box",
			ContentSelector: spiders.ContentSelector{
				ContainerSelector: ".ratings-contents__item",
				ChildSelector: map[string]string{
					"ratingRange":  ".ratings-contents__item-score",
					"ratingCounts": ".ratings-contents__item-num-strong",
				},
			},
		},
	}
	tbListSpider := spiders.NewtabelogListContentSpider(tbcListRequest)
	err = tbListSpider.Collect()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tbListSpider.GetCollections())
}
