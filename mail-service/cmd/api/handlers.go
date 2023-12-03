package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type mailResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (s *Server) SendMail(c *gin.Context) {
	var req mailMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := Message{
		From:    req.From,
		To:      req.To,
		Subject: req.Subject,
		Data:    req.Message,
	}

	err := s.mail.SendSMTPMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := mailResponse{
		Error:   false,
		Message: "sent to " + req.To,
	}

	c.JSON(http.StatusOK, resp)
}
