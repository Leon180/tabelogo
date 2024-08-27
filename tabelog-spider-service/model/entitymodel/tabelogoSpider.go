package entitymodel

import (
	"tabelog-spider/model/enum"

	"github.com/gocolly/colly"
)

type GetTabelogInfoParam struct {
	Area          string
	PlaceName     string
	MaxLinkAmount *int
}

type GetTabelogLinkListParam struct {
	URL           enum.URL
	LinkSelector  enum.LinkSelector
	MaxLinkAmount *int
}

type GetTabelogContentParam struct {
	URL                enum.URL
	ContainerCollector ContainerCollector
	CollectCondition   func(*colly.HTMLElement) bool
}

type GetTabelogTwoLevelContentParam struct {
	URL                        enum.URL
	TwoLevelContainerCollector TwoLevelContainerCollector
	CollectCondition           func(*colly.HTMLElement) bool
}

type ElementCollector struct {
	ElementName      enum.ElementSelector
	Collection       []string
	CollectCondition func(*colly.HTMLElement) bool
}

type ElementCollectorSlice []ElementCollector

func (es ElementCollectorSlice) ToMap() map[enum.ElementSelector][]string {
	m := make(map[enum.ElementSelector][]string, 0)
	for _, e := range es {
		m[e.ElementName] = e.Collection
	}
	return m
}

type ContainerCollector struct {
	ContainerName         enum.ContainerSelector
	ElementCollectorSlice ElementCollectorSlice
	CollectCondition      func(*colly.HTMLElement) bool
}

func (cs ContainerCollector) ToMap() map[enum.ElementSelector][]string {
	return cs.ElementCollectorSlice.ToMap()
}

type TwoLevelContainerCollector struct {
	ContainerName      enum.ContainerSelector
	ContainerCollector ContainerCollector
	CollectCondition   func(*colly.HTMLElement) bool
}

func (tccs TwoLevelContainerCollector) ToMap() map[enum.ElementSelector][]string {
	return tccs.ContainerCollector.ToMap()
}
