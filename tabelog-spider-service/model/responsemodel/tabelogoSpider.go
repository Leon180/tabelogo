package responsemodel

import "tabelog-spider/model/enum"

type TabelogInfo struct {
	Link        enum.URL `json:"link"`
	Name        string   `json:"name"`
	Rating      string   `json:"rating"`
	RatingCount string   `json:"ratingCount"`
	Bookmarks   string   `json:"bookmarks"`
	Phone       string   `json:"phone"`
	Type        []string `json:"type"`
}

type TabelogInfoResponse struct {
	TabelogInfos []TabelogInfo `json:"tabelog_infos"`
}

type TabelogPhoto struct {
	Link  enum.URL `json:"link"`
	Photo []string `json:"photo"`
}
