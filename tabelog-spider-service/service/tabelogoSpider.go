package service

import (
	"context"
	"strings"
	"sync"
	"tabelog-spider/model/entitymodel"
	"tabelog-spider/model/enum"
	"tabelog-spider/utility"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/samber/lo"
)

type GetTabelogInfoHandler interface {
	GetTabelogInfo(ctx context.Context, param entitymodel.GetTabelogInfoParam) (entitymodel.TabelogElementInfoSlice, error)
}

func NewGetTabelogInfoHandler() GetTabelogInfoHandler {
	return GetTabelogInfoHandle{}
}

type GetTabelogInfoHandle struct{}

func (handle GetTabelogInfoHandle) GetTabelogInfo(
	ctx context.Context,
	param entitymodel.GetTabelogInfoParam,
) (
	entitymodel.TabelogElementInfoSlice,
	error,
) {
	// Get Tabelog Link List in the restaurant list page
	linkList, err := getTabelogLinkList(ctx, entitymodel.GetTabelogLinkListParam{
		URL: enum.TabelogBaseURL.AddSubDirectory(param.Area).AddSubDirectory(enum.TabelogSubDirectoryRstLst.ToString()).AddParams(map[enum.Param]string{
			enum.TabelogParamVS: "1",
			enum.TabelogParamSK: param.PlaceName,
			enum.TabelogParamSW: param.PlaceName,
		}),
		LinkSelector:  enum.LinkSelectorTabelogRestaurantLinkList,
		MaxLinkAmount: param.MaxLinkAmount,
	})
	if err != nil {
		return nil, err
	}

	// Get Tabelog Content
	tabelogInfo := make(entitymodel.TabelogElementInfoSlice, len(linkList))
	var wg sync.WaitGroup
	wg.Add(len(linkList))
	for i, link := range linkList {
		go func(i int, link enum.URL) {
			defer wg.Done()
			// get restaurant base info
			elements, err := getElementsByContainerCollector(ctx, entitymodel.GetTabelogContentParam{
				URL: link,
				ContainerCollector: entitymodel.ContainerCollector{
					ContainerName: enum.ContainerSelectorTabelogRestaurantInfo,
					ElementCollectorSlice: entitymodel.ElementCollectorSlice{
						{ElementName: enum.ElementSelectorTabelogRestaurantName},
						{ElementName: enum.ElementSelectorTabelogRestaurantRating},
						{ElementName: enum.ElementSelectorTabelogRestaurantRatingCount},
						{ElementName: enum.ElementSelectorTabelogRestaurantBookmarks},
						{ElementName: enum.ElementSelectorTabelogRestaurantPhone},
					},
				},
			})
			if err != nil {
				return
			}
			tabelogInfo[i] = entitymodel.TabelogElementInfo{
				Link:                  link,
				ElementCollectorSlice: elements,
			}
			// get restaurant type
			elements, err = getElementsByTwoLevelContainerCollector(ctx, entitymodel.GetTabelogTwoLevelContentParam{
				URL: link,
				TwoLevelContainerCollector: entitymodel.TwoLevelContainerCollector{
					ContainerName: enum.ContainerSelectorTabelogRestaurantType,
					CollectCondition: func(e *colly.HTMLElement) bool {
						return e.ChildText(enum.ContainerSelectorTabelogRestaurantTypeCondition.ToString()) == enum.TabelogoContainerSelectorCondition[enum.ContainerSelectorTabelogRestaurantTypeCondition]
					},
					ContainerCollector: entitymodel.ContainerCollector{
						ContainerName: enum.ContainerSelectorTabelogRestaurantTypeTwoLevel,
						ElementCollectorSlice: entitymodel.ElementCollectorSlice{
							{ElementName: enum.ElementSelectorTabelogRestaurantType},
						},
					},
				},
			})
			if err != nil {
				return
			}
			tabelogInfo[i].ElementCollectorSlice = append(tabelogInfo[i].ElementCollectorSlice, elements...)
		}(i, link)
	}
	wg.Wait()

	return tabelogInfo, nil
}

func getTabelogLinkList(
	ctx context.Context,
	param entitymodel.GetTabelogLinkListParam,
) ([]enum.URL, error) {
	links := []string{}
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	c.SetRequestTimeout(2 * time.Second)
	c.OnHTML(param.LinkSelector.ToString(), func(e *colly.HTMLElement) {
		links = append(links, e.Attr(enum.AttributeSelectorHref.ToString()))
	})
	if err := c.Visit(string(param.URL)); err != nil {
		utility.SugarLogger.Error(err)
		return nil, err
	}
	urlLinks := lo.Uniq(lo.Map(links, func(link string, _ int) enum.URL {
		return enum.URL(link)
	}))
	if param.MaxLinkAmount != nil {
		if len(urlLinks) > *param.MaxLinkAmount {
			urlLinks = urlLinks[:*param.MaxLinkAmount:*param.MaxLinkAmount]
		}
	}
	return urlLinks, nil
}

func getElementsByContainerCollector(
	ctx context.Context,
	param entitymodel.GetTabelogContentParam,
) (
	entitymodel.ElementCollectorSlice,
	error,
) {
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	c.OnHTML(param.ContainerCollector.ContainerName.ToString(), func(element *colly.HTMLElement) {
		for i := range param.ContainerCollector.ElementCollectorSlice {
			element.ForEach(param.ContainerCollector.ElementCollectorSlice[i].ElementName.ToString(), func(_ int, e *colly.HTMLElement) {
				param.ContainerCollector.ElementCollectorSlice[i].Collection = append(param.ContainerCollector.ElementCollectorSlice[i].Collection, strings.TrimSpace(e.Text))
			})
		}
	})
	if err := c.Visit(string(param.URL)); err != nil {
		utility.SugarLogger.Error(err)
		return entitymodel.ElementCollectorSlice{}, err
	}
	return param.ContainerCollector.ElementCollectorSlice, nil
}

func getElementsByTwoLevelContainerCollector(
	ctx context.Context,
	param entitymodel.GetTabelogTwoLevelContentParam,
) (
	entitymodel.ElementCollectorSlice,
	error,
) {
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	c.OnHTML(param.TwoLevelContainerCollector.ContainerName.ToString(), func(element *colly.HTMLElement) {
		if param.TwoLevelContainerCollector.CollectCondition(element) {
			element.ForEach(param.TwoLevelContainerCollector.ContainerCollector.ContainerName.ToString(), func(_ int, ele *colly.HTMLElement) {
				for i := range param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice {
					param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice[i].Collection = append(param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice[i].Collection, ele.ChildText(param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice[i].ElementName.ToString()))
				}
			})
		}
	})
	if err := c.Visit(string(param.URL)); err != nil {
		utility.SugarLogger.Error(err)
		return entitymodel.ElementCollectorSlice{}, err
	}
	return param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice, nil
}

type GetTabelogPhotoHandler interface {
	GetTabelogPhoto(ctx context.Context, link enum.URL) (entitymodel.TabelogElementInfo, error)
}

func NewGetTabelogPhotoHandler() GetTabelogPhotoHandler {
	return GetTabelogPhotoHandle{}
}

type GetTabelogPhotoHandle struct{}

func (handle GetTabelogPhotoHandle) GetTabelogPhoto(
	ctx context.Context,
	link enum.URL,

) (
	entitymodel.TabelogElementInfo,
	error,
) {
	elements, err := getTabelogPhoto(ctx, entitymodel.GetTabelogTwoLevelContentParam{
		URL: link.AddSubDirectory(enum.TabelogSubDirectoryRstPhotoLst.ToString()),
		TwoLevelContainerCollector: entitymodel.TwoLevelContainerCollector{
			ContainerName: enum.ContainerSelectorTabelogRestaurantPhotoList,
			ContainerCollector: entitymodel.ContainerCollector{
				ContainerName: enum.ContainerSelectorTabelogRestaurantPhotoListTwoLevel,
				ElementCollectorSlice: entitymodel.ElementCollectorSlice{
					{ElementName: enum.ElementSelectorTabelogRestaurantPhoto},
				},
			},
		},
	})
	if err != nil {
		return entitymodel.TabelogElementInfo{}, err
	}
	return entitymodel.TabelogElementInfo{
		Link:                  link,
		ElementCollectorSlice: elements,
	}, nil
}

func getTabelogPhoto(
	ctx context.Context,
	param entitymodel.GetTabelogTwoLevelContentParam,
) (
	entitymodel.ElementCollectorSlice,
	error,
) {
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	c.OnHTML(param.TwoLevelContainerCollector.ContainerName.ToString(), func(element *colly.HTMLElement) {
		element.ForEach(param.TwoLevelContainerCollector.ContainerCollector.ContainerName.ToString(), func(_ int, ele *colly.HTMLElement) {
			for i := range param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice {
				param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice[i].Collection = append(param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice[i].Collection, ele.ChildAttr(param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice[i].ElementName.ToString(), enum.AttributeSelectorSrc.ToString()))
			}
		})
	})
	if err := c.Visit(string(param.URL)); err != nil {
		utility.SugarLogger.Error(err)
		return entitymodel.ElementCollectorSlice{}, err
	}
	if len(param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice) == 0 {
		return entitymodel.ElementCollectorSlice{}, nil
	}
	return param.TwoLevelContainerCollector.ContainerCollector.ElementCollectorSlice, nil
}
