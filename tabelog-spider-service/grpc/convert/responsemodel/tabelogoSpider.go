package responsemodel

import (
	"tabelog-spider/grpc/proto"
	"tabelog-spider/model/entitymodel"
	"tabelog-spider/model/enum"
)

type TabelogElementInfoSlice entitymodel.TabelogElementInfoSlice

func (entity TabelogElementInfoSlice) TabelogInfoResponseListProto() *proto.TabelogInfoResponseList {
	slice := []*proto.TabelogInfoResponse{}
	for _, v := range entity {
		tmp := &proto.TabelogInfoResponse{
			Link: string(v.Link),
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
	return &proto.TabelogInfoResponseList{
		TabelogInfos: slice,
	}
}

type GetTabelogPhoto entitymodel.TabelogElementInfo

func (g GetTabelogPhoto) ToTabelogPhotoResponseProto() *proto.TabelogPhotoResponse {
	resp := &proto.TabelogPhotoResponse{
		Link: string(g.Link),
	}
	if len(g.ElementCollectorSlice) == 0 {
		return resp
	}
	resp.Photo = g.ElementCollectorSlice[0].Collection
	return resp
}
