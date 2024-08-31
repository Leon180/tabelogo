package entitymodel

import "fmt"

type QuickSearchRequest struct {
	PlaceID      string
	APIMask      string
	LanguageCode string
}

type AdvanceSearchRequest struct {
	TextQuery      string
	LowLatitude    float64
	LowLongitude   float64
	HighLatitude   float64
	HighLongitude  float64
	MaxResultCount int
	MinRating      int
	OpenNow        bool
	RankPreference string
	LanguageCode   string
	APIMask        string
}

func (r AdvanceSearchRequest) ToRequestBody() string {
	return fmt.Sprintf(`
		{
			"textQuery":"%s",
			"locationBias":{
				"rectangle": {
					"low": {
						"latitude": %f,
						"longitude": %f
					},
					"high": {
						"latitude": %f,
						"longitude": %f
					}
				}
			},
			"maxResultCount":%d,
			"minRating":%d,
			"openNow":%t,
			"rankPreference":"%s",
			"languageCode":"%s"
		}
		`, r.TextQuery, r.LowLatitude, r.LowLongitude, r.HighLatitude, r.HighLongitude, r.MaxResultCount, r.MinRating, r.OpenNow, r.RankPreference, r.LanguageCode)
}
