package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	RegisterHandler(func(engine *gin.Engine) {
		engine.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	})
}
