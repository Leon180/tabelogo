package responsemodel

import "tabelog-spider/model/enum"

type TabelogInfoResponse struct {
	Link        enum.URL `json:"link"`
	Name        string   `json:"name"`
	Rating      string   `json:"rating"`
	RatingCount string   `json:"ratingCount"`
	Bookmarks   string   `json:"bookmarks"`
	Phone       string   `json:"phone"`
	Type        []string `json:"type"`
}

type TabelogInfoResponseList struct {
	TabelogInfos []TabelogInfoResponse `json:"tabelog_infos"`
}

type TabelogPhotoResponse struct {
	Link  enum.URL `json:"link"`
	Photo []string `json:"photo"`
}
