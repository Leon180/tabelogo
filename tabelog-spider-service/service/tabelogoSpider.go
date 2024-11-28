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
	"golang.org/x/time/rate"
)

type GetTabelogInfoHandler interface {
	GetTabelogInfo(ctx context.Context, param entitymodel.GetTabelogInfoParam) (entitymodel.TabelogElementInfoSlice, error)
}

func NewGetTabelogInfoHandler(
	limiter *rate.Limiter,
) GetTabelogInfoHandler {
	return &GetTabelogInfoHandle{
		limiter: limiter,
	}
}

type GetTabelogInfoHandle struct {
	limiter *rate.Limiter
}

func (handle *GetTabelogInfoHandle) GetTabelogInfo(
	ctx context.Context,
	param entitymodel.GetTabelogInfoParam,
) (
	entitymodel.TabelogElementInfoSlice,
	error,
) {
	// Get Tabelog Link List in the restaurant list page
	linkList, err := handle.getTabelogLinkList(ctx, entitymodel.GetTabelogLinkListParam{
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

	return handle.scrapeRestaurantInfoConcurrently(ctx, linkList)
}

func (handle *GetTabelogInfoHandle) scrapeRestaurantInfoConcurrently(ctx context.Context, linkList []enum.URL) (entitymodel.TabelogElementInfoSlice, error) {
	// Create error channel for goroutine error handling
	errChan := make(chan error, len(linkList))
	tabelogInfo := make(entitymodel.TabelogElementInfoSlice, len(linkList))

	var wg sync.WaitGroup
	wg.Add(len(linkList))

	for i, link := range linkList {
		go func(i int, link enum.URL) {
			defer wg.Done()

			// Check context before starting work
			if ctx.Err() != nil {
				errChan <- ctx.Err()
				return
			}

			if err := handle.limiter.Wait(ctx); err != nil {
				errChan <- err
				return
			}

			info, err := handle.scrapeRestaurantInfo(ctx, link)
			if err != nil {
				errChan <- err
				return
			}

			tabelogInfo[i] = info
		}(i, link)
	}

	// Wait with context
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errChan)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case <-done:
		return tabelogInfo, nil
	}

	return tabelogInfo, nil
}

func (handle *GetTabelogInfoHandle) scrapeRestaurantInfo(ctx context.Context, link enum.URL) (entitymodel.TabelogElementInfo, error) {
	// Get base info
	baseInfo, err := handle.scrapeBaseInfo(ctx, link)
	if err != nil {
		return entitymodel.TabelogElementInfo{}, err
	}

	// Get restaurant type
	typeInfo, err := handle.scrapeRestaurantType(ctx, link)
	if err != nil {
		return entitymodel.TabelogElementInfo{}, err
	}

	return entitymodel.TabelogElementInfo{
		Link:                  link,
		ElementCollectorSlice: append(baseInfo, typeInfo...),
	}, nil
}

func (handle *GetTabelogInfoHandle) scrapeBaseInfo(ctx context.Context, link enum.URL) (entitymodel.ElementCollectorSlice, error) {
	return handle.getElementsByContainerCollector(entitymodel.GetTabelogContentParam{
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
}

func (handle *GetTabelogInfoHandle) scrapeRestaurantType(ctx context.Context, link enum.URL) (entitymodel.ElementCollectorSlice, error) {
	return handle.getElementsByTwoLevelContainerCollector(entitymodel.GetTabelogTwoLevelContentParam{
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
}

func (handle *GetTabelogInfoHandle) getTabelogLinkList(
	ctx context.Context,
	param entitymodel.GetTabelogLinkListParam,
) ([]enum.URL, error) {
	links := []string{}
	c := handle.newConfiguredCollector()
	c.OnHTML(param.LinkSelector.ToString(), func(e *colly.HTMLElement) {
		links = append(links, e.Attr(enum.AttributeSelectorHref.ToString()))
	})
	c.OnError(func(r *colly.Response, err error) {
		utility.SugarLogger.Errorw("failed to fetch links",
			"url", param.URL,
			"status_code", r.StatusCode,
			"error", err,
		)
	})
	// Visit URL with retries
	err := handle.visitWithRetries(ctx, c, string(param.URL))
	if err != nil {
		return nil, err
	}

	return handle.processLinks(links, param.MaxLinkAmount), nil
}

func (handle *GetTabelogInfoHandle) visitWithRetries(ctx context.Context, collector *colly.Collector, url string) error {
	backoff := utility.NewExponentialBackoff(2, 100*time.Millisecond, 10*time.Second)

	return utility.RetryWithBackoff(ctx, 2, backoff, func() error {
		return collector.Visit(url)
	})
}

// processLinks handles link deduplication and limiting
func (handle *GetTabelogInfoHandle) processLinks(links []string, maxLinks *int) []enum.URL {
	// Convert to URLs and deduplicate
	urlLinks := lo.Uniq(lo.Map(links, func(link string, _ int) enum.URL {
		return enum.URL(link)
	}))

	// Apply max limit if specified
	if maxLinks != nil && len(urlLinks) > *maxLinks {
		return urlLinks[:*maxLinks]
	}

	return urlLinks
}

func (handle *GetTabelogInfoHandle) newConfiguredCollector() *colly.Collector {
	collector := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector)
		},
	)
	collector.SetRequestTimeout(2 * time.Second)
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*tabelog.*",
		Parallelism: 1,
		RandomDelay: 1 * time.Second,
	})

	return collector
}

func (handle *GetTabelogInfoHandle) getElementsByContainerCollector(
	param entitymodel.GetTabelogContentParam,
) (
	entitymodel.ElementCollectorSlice,
	error,
) {
	c := handle.newConfiguredCollector()
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

func (handle *GetTabelogInfoHandle) getElementsByTwoLevelContainerCollector(
	param entitymodel.GetTabelogTwoLevelContentParam,
) (
	entitymodel.ElementCollectorSlice,
	error,
) {
	c := handle.newConfiguredCollector()
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
	elements, err := getTabelogPhoto(entitymodel.GetTabelogTwoLevelContentParam{
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
