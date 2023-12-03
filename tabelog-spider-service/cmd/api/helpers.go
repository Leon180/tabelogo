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
		nu = nu + k + "=" + url.QueryEscape(v) + "&"
	}
	return nu
}

func RemoveDuplicateString(s []string) []string {
	m := make(map[string]byte)
	var r []string
	for _, v := range s {
		l := len(m)
		m[v] = 0
		if len(m) != l {
			r = append(r, v)
		}
	}
	return r
}
