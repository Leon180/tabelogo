package main

import (
	"broker/rabbitmq/event"
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

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (s *Server) logEventViaRabbit(c *gin.Context) {
	var l LogPayload
	if err := c.ShouldBindJSON(&l); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err := s.pushToQueue(l.Name, l.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var resp JSONResponse
	resp.Error = false
	resp.Message = "logged via RabbitMQ"

	c.JSON(http.StatusOK, resp)
}

// pushToQueue pushes a message into RabbitMQ
func (s *Server) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(s.rabbitMQ)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, err := json.MarshalIndent(&payload, "", "\t")
	if err != nil {
		return err
	}

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
