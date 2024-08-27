package enum

type ElementSelector string

func (e ElementSelector) ToString() string {
	return string(e)
}

const (
	ElementSelectorTabelogRestaurantName        ElementSelector = "h2.display-name"
	ElementSelectorTabelogRestaurantRating      ElementSelector = ".rdheader-rating__score b.c-rating__val"
	ElementSelectorTabelogRestaurantRatingCount ElementSelector = ".rdheader-rating__review-target .num"
	ElementSelectorTabelogRestaurantBookmarks   ElementSelector = ".rdheader-rating__hozon-target .num"
	ElementSelectorTabelogRestaurantPhone       ElementSelector = ".rstinfo-table__tel-num"
	ElementSelectorTabelogRestaurantType        ElementSelector = ".linktree__parent-target-text"
	ElementSelectorTabelogRestaurantPhoto       ElementSelector = ".rstdtl-photo-list__img"
)

type LinkSelector string

func (l LinkSelector) ToString() string {
	return string(l)
}

const (
	LinkSelectorTabelogRestaurantLinkList LinkSelector = ".list-rst__rst-name-target"
)

type ContainerSelector string

func (c ContainerSelector) ToString() string {
	return string(c)
}

const (
	ContainerSelectorTabelogRestaurantInfo              ContainerSelector = "#container"
	ContainerSelectorTabelogRestaurantType              ContainerSelector = ".rdheader-subinfo__item"
	ContainerSelectorTabelogRestaurantTypeTwoLevel      ContainerSelector = ".linktree__parent"
	ContainerSelectorTabelogRestaurantTypeCondition     ContainerSelector = ".rdheader-subinfo__item-title"
	ContainerSelectorTabelogRestaurantPhotoList         ContainerSelector = ".rstdtl-photo-list"
	ContainerSelectorTabelogRestaurantPhotoListTwoLevel ContainerSelector = ".rstdtl-photo-list__item"
)

type ContainerSelectorCondition map[ContainerSelector]string

var TabelogoContainerSelectorCondition ContainerSelectorCondition = ContainerSelectorCondition{
	ContainerSelectorTabelogRestaurantTypeCondition: "ジャンル：",
}

type AttributeSelector string

func (a AttributeSelector) ToString() string {
	return string(a)
}

const (
	AttributeSelectorHref AttributeSelector = "href"
	AttributeSelectorSrc  AttributeSelector = "src"
)
