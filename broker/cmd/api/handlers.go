package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// For verify the connection
func (s *Server) Broker(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func (s *Server) TransRequest(method, url string) func(*gin.Context) {
	return func(c *gin.Context) {
		fmt.Println(url)
		var err error
		body := c.Request.Body
		// submit request to tabelog spider service
		request, err := http.NewRequest(method, url, body)
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

		// return response from tabelog spider service
		var resp interface{}
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
