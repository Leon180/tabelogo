package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type tabelogSpiderRequest struct {
	SearchName string `json:"searchName"`
	SearchArea string `json:"searchArea"`
}

func (s *Server) Broker(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func (s *Server) TabelogSpider(ctx *gin.Context) {
	var req tabelogSpiderRequest
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// submit request to tabelog spider service
	request, err := http.NewRequest("POST", tabelogSpiderServiceURL, bytes.NewBuffer([]byte(`{"searchName":"`+req.SearchName+`","searchArea":"`+req.SearchArea+`"}`)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, "status code is not 200")
		return
	}

	// return response from tabelog spider service
	var resp interface{}
	json.NewDecoder(response.Body).Decode(&resp)
	ctx.JSON(http.StatusOK, resp)
}
