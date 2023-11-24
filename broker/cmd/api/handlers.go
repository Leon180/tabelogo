package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) Broker(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) Regist(c *gin.Context) {
	var req CreateUserRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// submit request to tabelog spider service
	request, err := http.NewRequest("POST", authenticateServiceURL+"/regist", bytes.NewBuffer([]byte(`{"email":"`+req.Email+`","password":"`+req.Password+`"}`)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	fmt.Println("request: ", request)
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

type tabelogSpiderRequest struct {
	SearchName string `json:"searchName"`
	SearchArea string `json:"searchArea"`
}

func (s *Server) TabelogSpider(c *gin.Context) {
	var req tabelogSpiderRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// submit request to tabelog spider service
	request, err := http.NewRequest("POST", tabelogSpiderServiceURL, bytes.NewBuffer([]byte(`{"searchName":"`+req.SearchName+`","searchArea":"`+req.SearchArea+`"}`)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, "status code is not 200")
		return
	}

	// return response from tabelog spider service
	var resp interface{}
	json.NewDecoder(response.Body).Decode(&resp)
	c.JSON(http.StatusOK, resp)
}
