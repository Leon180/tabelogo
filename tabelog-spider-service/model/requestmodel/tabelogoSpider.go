package requestmodel

import "tabelog-spider/model/enum"

type GetTabelogInfoRequest struct {
	Area            string `form:"area" binding:"required"`
	PlaceName       string `form:"place_name" binding:"required"`
	MaxResultAmount *int   `form:"max_result_amount"`
}

type GetTabelogPhotoRequest struct {
	Link enum.URL `form:"link" binding:"required"`
}
