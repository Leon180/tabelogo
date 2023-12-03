package main

import (
	data "logger-service/cmd/data"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type JSONResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (s *Server) WriteLog(c *gin.Context) {
	var req JSONPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	if err := s.Models.LogEntry.InsertOne(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Log entry created successfully",
	}

	c.JSON(http.StatusOK, resp)
}
