package main

import "github.com/gin-gonic/gin"

type LinkSpiderRequest struct {
	City      string `json:"city" binding:"required"`
	Area      string `json:"area" binding:"required"`
	PlaceName string `json:"place_name" binding:"required"`
}

func (s *Server) TabelogSpider(c *gin.Context) {
}
