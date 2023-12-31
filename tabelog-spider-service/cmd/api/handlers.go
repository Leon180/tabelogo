package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
)

type LinkSpiderRequest struct {
	GoogleID  string `json:"google_id" binding:"required"`
	Area      string `json:"area" binding:"required"`
	PlaceName string `json:"place_name" binding:"required"`
}

type TabelogInfo struct {
	Link        string
	Name        string
	Rating      string
	RatingCount string
	Bookmarks   string
	Phone       string
	Type        []string
}

// table photo link: link+"/table/"
// menu link: link+"/dtlmenu/"
//	drink menu link: link+"/dtlmenu/drink/"
// comments link: link+"/dtlrvwlst/"
// rating distribute link: link+"/dtlratings/"
// photo link: link+"/dtlphotolst

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
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	lCollection := lSpider.GetCollections()
	if len(lCollection) == 0 {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no search result")))
		return
	}
	lCollection = RemoveDuplicateString(lCollection)
	if len(lCollection) > maxCollectLinks {
		lCollection = lCollection[:maxCollectLinks:maxCollectLinks]
	}

	tabelogInfoCollection := make([]TabelogInfo, len(lCollection))
	var wg sync.WaitGroup
	wg.Add(len(lCollection))
	for index, link := range lCollection {
		go func(link string, index int) {
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
						"phone":       ".rstinfo-table__tel-num",
					},
				},
			}
			tbSpider := NewtabelogContentSpider(tbcRequest)
			err = tbSpider.Collect()
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			tbc := tbSpider.GetCollections()

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
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			tbcList := tbListSpider.GetCollections()
			typeList := make([]string, len(tbcList))
			for i, tbc := range tbcList {
				typeList[i] = tbc["type"][0]
			}

			tabelogInfoCollection[index] = TabelogInfo{
				Link:        link,
				Name:        tbc["name"][0],
				Rating:      tbc["rating"][0],
				RatingCount: tbc["ratingCount"][0],
				Bookmarks:   tbc["bookmarks"][0],
				Phone:       tbc["phone"][0],
				Type:        typeList,
			}
		}(link, index)
	}
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{
		"gooogle_id": req.GoogleID,
		"result":     tabelogInfoCollection,
	})
}

type TabelogPhotoRequest struct {
	GoogleID string `json:"google_id" binding:"required"`
	Link     string `json:"link" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type TabelogPhoto struct {
	Link  string
	Name  string
	Photo []string
}

func (s *Server) TabelogPhotoSpider(c *gin.Context) {
	var req TabelogPhotoRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// get photo
	tbListOption := TabelogListContentSpider{
		Url: req.Link + "dtlphotolst",
		ListCondition: func(e *colly.HTMLElement) bool {
			return true
		},
		ListContentSelector: ListContentSelector{
			ParentContainerSelector: ".rstdtl-photo-list",
			ContentSelector: ContentSelector{
				ContainerSelector: ".rstdtl-photo-list__item",
				ChildSelector: map[string]string{
					"img": ".rstdtl-photo-list__img",
				},
			},
		},
	}
	tbListSpider := NewtabelogListContentSpider(tbListOption)
	err = tbListSpider.Collect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	tbcList := tbListSpider.GetCollections()
	if len(tbcList) == 0 {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("no photo")))
		return
	}
	imgList := make([]string, len(tbcList))
	for i, tbc := range tbcList {
		imgList[i] = tbc["img"][0]
	}

	var resp = TabelogPhoto{
		Link:  req.Link,
		Name:  req.Name,
		Photo: imgList,
	}
	c.JSON(http.StatusOK, gin.H{
		"gooogle_id": req.GoogleID,
		"result":     resp,
	})
}
