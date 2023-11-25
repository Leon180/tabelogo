package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

type LinkSpiderRequest struct {
	Area      string `json:"area" binding:"required"`
	PlaceName string `json:"place_name" binding:"required"`
}

type TabelogInfo struct {
	Name        string
	Rating      string
	RatingCount string
	Bookmarks   string
	Type        []string
}

func (s *Server) TabelogSpider(c *gin.Context) {
	var req LinkSpiderRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// post tabelog spider service
	lSpiderOption := LinkSpider{
		Url: `https://tabelog.com/` + req.Area + `/rstLst/`,
		UrlParams: map[string]string{
			"vs": "1",
			"sk": req.PlaceName,
			"sw": req.PlaceName,
		},
		LinkSelector: ".list-rst__rst-name-target",
	}
	lSpider := NewLinkSpider(lSpiderOption)
	err = lSpider.Collect()
	if err != nil {
		fmt.Println(err)
	}
	lCollection := lSpider.GetCollections()
	if len(lCollection) == 0 {
		fmt.Println("No links found")
		c.JSON(http.StatusBadRequest, fmt.Errorf("no links found"))
		return
	}
	if len(lCollection) > maxCollectLinks {
		lCollection = lCollection[:maxCollectLinks:maxCollectLinks]
	}

	tabelogInfoCollection := make([]TabelogInfo, len(lCollection))
	var wg sync.WaitGroup
	wg.Add(len(lCollection))
	for i, link := range lCollection {
		fmt.Println(link)
		go func(link string) {
			// table photo link: link+"/table/"
			// menu link: link+"/dtlmenu/"
			//	drink menu link: link+"/dtlmenu/drink/"
			// comments link: link+"/dtlrvwlst/"
			// rating distribute link: link+"/dtlratings/"
			defer wg.Done()
			tbcRequest := TabelogContentSpider{
				Url: link,
				ContentSelector: ContentSelector{
					ContainerSelector: "#container",
					ChildSelector: map[string]string{
						"name":        "h2.display-name",
						"rating":      ".rdheader-rating__score b.c-rating__val",
						"ratingCount": ".rdheader-rating__review-target .num",
						"bookmarks":   ".rdheader-rating__hozon-target .num",
					},
				},
			}
			tbSpider := NewtabelogContentSpider(tbcRequest)
			err = tbSpider.Collect()
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			tbc := tbSpider.GetCollections()
			fmt.Println(tbc)
			tabelogInfoCollection[i] = TabelogInfo{
				Name:        tbc["name"][0],
				Rating:      tbc["rating"][0],
				RatingCount: tbc["ratingCount"][0],
				Bookmarks:   tbc["bookmarks"][0],
			}

			// get type
			tbListOption := TabelogListContentSpider{
				Url: link,
				ListCondition: func(e *colly.HTMLElement) bool {
					return e.ChildText(".rdheader-subinfo__item-title") == "ジャンル："
				},
				ListContentSelector: ListContentSelector{
					ParentContainerSelector: ".rdheader-subinfo__item",
					ContentSelector: ContentSelector{
						ContainerSelector: ".linktree__parent",
						ChildSelector: map[string]string{
							"type": ".linktree__parent-target-text",
						},
					},
				},
			}
			tbListSpider := NewtabelogListContentSpider(tbListOption)
			err = tbListSpider.Collect()
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			tbcList := tbListSpider.GetCollections()
			typeList, exist := tbcList[0]["type"]
			if !exist {
				fmt.Println("type not found")
				c.JSON(http.StatusInternalServerError, fmt.Errorf("type not found in tabelog"))
				return
			}
			tabelogInfoCollection[i].Type = typeList
		}(link)
	}
	wg.Wait()
	c.JSON(http.StatusOK, tabelogInfoCollection)
}
