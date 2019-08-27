package handlers

import (
	"github.com/gin-gonic/gin"
)

// NewHandler register routers, middlewares
func NewHandler() *gin.Engine {
	r := gin.New()
	r.POST("/namespace/delete-check", NamespaceDeleteCheck)
	return r
}
