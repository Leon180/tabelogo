package main

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func ParamsCombine(u string, params map[string]string) string {
	nu := u + "?"
	for k, v := range params {
		if k == "sk" || k == "sw" {
			nu = nu + k + "=" + url.QueryEscape(v) + "&"
		}
	}
	return nu
}
