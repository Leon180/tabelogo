package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuickSearchRequest struct {
	PlaceID      string `json:"place_id" binding:"required"`
	APIMask      string `json:"api_mask"`
	LanguageCode string `json:"language_code" binding:"required"`
}

func (s *Server) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World!",
	})
}

func (s *Server) QuickSearch(c *gin.Context) {
	var req QuickSearchRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// post google api
	request, err := http.NewRequest(
		"GET",
		`https://places.googleapis.com/v1/places/`+req.PlaceID+`?fields=`+req.APIMask+`&key=`+s.config.GoogleMapAPIKey+`&languageCode=`+req.LanguageCode,
		nil,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer response.Body.Close()

	// return response from tabelog spider service
	var resp interface{}
	json.NewDecoder(response.Body).Decode(&resp)
	c.JSON(http.StatusOK, resp)
}

type AdvanceSearchRequest struct {
	TextQuery string `json:"text_query" binding:"required"`
	// For location bias
	LowLatitude   float64 `json:"low_latitude" binding:"required"`
	LowLongitude  float64 `json:"low_longitude" binding:"required"`
	HighLatitude  float64 `json:"high_latitude" binding:"required"`
	HighLongitude float64 `json:"high_longitude" binding:"required"`
	//
	MaxResultCount int    `json:"max_result_count" binding:"required"`
	MinRating      int    `json:"min_rating" binding:"required"`
	OpenNow        bool   `json:"open_now"`
	RankPreference string `json:"rank_preference" binding:"required"`
	LanguageCode   string `json:"language_code" binding:"required"`
	APIMask        string `json:"api_mask"`
}

func (s *Server) AdvanceSearch(c *gin.Context) {
	var req AdvanceSearchRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	bBuffer := bytes.NewBuffer(
		[]byte(`
		{
			"textQuery":"` + req.TextQuery + `",
			"locationBias":{
				"rectangle": {
						"low": {
							"latitude": ` + fmt.Sprintf("%f", req.LowLatitude) + `,
							"longitude":` + fmt.Sprintf("%f", req.LowLongitude) + `
						},
						"high": {
							"latitude": ` + fmt.Sprintf("%f", req.HighLatitude) + `,
							"longitude": ` + fmt.Sprintf("%f", req.HighLongitude) + `
						}
					}
				},
			"maxResultCount":` + strconv.Itoa(req.MaxResultCount) + `,
			"minRating":` + strconv.Itoa(req.MinRating) + `,
			"openNow":` + strconv.FormatBool(req.OpenNow) + `,
			"rankPreference":"` + req.RankPreference + `",
			"languageCode":"` + req.LanguageCode + `"
		}
		`),
	)
	fmt.Println(bBuffer)
	// post google api
	request, err := http.NewRequest(
		"POST",
		"https://places.googleapis.com/v1/places:searchText",
		bBuffer,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// set header
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Goog-Api-Key", s.config.GoogleMapAPIKey)
	request.Header.Set("X-Goog-FieldMask", req.APIMask)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer response.Body.Close()

	// return response from tabelog spider service
	var resp interface{}
	json.NewDecoder(response.Body).Decode(&resp)
	c.JSON(http.StatusOK, resp)
}
