package convert

import (
	"tabelog-spider/model/entitymodel"
	"tabelog-spider/model/enum"
	"tabelog-spider/model/responsemodel"
)

type GetTabelogInfo entitymodel.TabelogElementInfoSlice

func (g GetTabelogInfo) ToResponse() responsemodel.TabelogInfoResponse {
	slice := []responsemodel.TabelogInfo{}
	for _, v := range g {
		tmp := responsemodel.TabelogInfo{
			Link: v.Link,
		}
		for _, e := range v.ElementCollectorSlice {
			if len(e.Collection) == 0 {
				continue
			}
			switch e.ElementName {
			case enum.ElementSelectorTabelogRestaurantName:
				tmp.Name = e.Collection[0]
			case enum.ElementSelectorTabelogRestaurantRating:
				tmp.Rating = e.Collection[0]
			case enum.ElementSelectorTabelogRestaurantRatingCount:
				tmp.RatingCount = e.Collection[0]
			case enum.ElementSelectorTabelogRestaurantBookmarks:
				tmp.Bookmarks = e.Collection[0]
			case enum.ElementSelectorTabelogRestaurantPhone:
				tmp.Phone = e.Collection[0]
			case enum.ElementSelectorTabelogRestaurantType:
				tmp.Type = e.Collection
			}
		}
		slice = append(slice, tmp)
	}
	return responsemodel.TabelogInfoResponse{
		TabelogInfos: slice,
	}
}

type GetTabelogPhoto entitymodel.TabelogElementInfo

func (g GetTabelogPhoto) ToResponse() responsemodel.TabelogPhoto {
	resp := responsemodel.TabelogPhoto{
		Link: g.Link,
	}
	if len(g.ElementCollectorSlice) == 0 {
		return resp
	}
	resp.Photo = g.ElementCollectorSlice[0].Collection
	return resp
}
