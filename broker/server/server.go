package server

import (
	"broker/utility"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Server struct {
	router   *gin.Engine
	rabbitMQ *amqp.Connection
}

func NewServer(
	router *gin.Engine,
	rabbitConn *amqp.Connection,
) *Server {
	return &Server{
		router:   router,
		rabbitMQ: rabbitConn,
	}
}

// For verify the connection
func (s *Server) Broker(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}

func (s *Server) HTTPTransRequest(method, url string) func(*gin.Context) {
	return func(c *gin.Context) {

		// create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
		defer cancel()

		// submit request to tabelog spider service
		request, err := http.NewRequestWithContext(ctx, method, url, c.Request.Body)
		if err != nil {
			utility.CommonErrorResponse(c, err, http.StatusBadRequest)
			return
		}

		// Copy only specific headers
		allowedHeaders := []string{"Content-Type", "Accept", "Authorization", "User-Agent"}
		for _, header := range allowedHeaders {
			if val := c.Request.Header.Get(header); val != "" {
				request.Header.Set(header, val)
			}
		}

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		response, err := client.Do(request)
		if err != nil {
			status := http.StatusServiceUnavailable
			switch {
			case os.IsTimeout(err):
				status = http.StatusGatewayTimeout
			case ctx.Err() == context.Canceled:
				status = http.StatusRequestTimeout
			}
			utility.CommonErrorResponse(c, err, status)
			return
		}
		defer response.Body.Close()

		// Copy relevant response headers back to client
		for _, header := range allowedHeaders {
			if val := response.Header.Get(header); val != "" {
				c.Header(header, val)
			}
		}

		// return response from tabelog spider service
		var resp interface{}
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			utility.CommonErrorResponse(c, err, http.StatusUnprocessableEntity)
			return
		}

		// return response
		c.JSON(response.StatusCode, resp)
	}
}
