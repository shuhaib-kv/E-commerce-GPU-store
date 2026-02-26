package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getIntQuery(c *gin.Context, key string, defaultVal int) int {
	val, err := strconv.Atoi(c.Query(key))
	if err != nil || val < 1 {
		return defaultVal
	}
	return val
}

func respondSuccess(c *gin.Context, code int, message string, data interface{}) {
	resp := gin.H{"status": true, "message": message}
	if data != nil {
		resp["data"] = data
	}
	c.JSON(code, resp)
}

func respondError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"status": false, "message": message})
}

// respondInternalError logs the actual error and returns a generic message to the client.
func respondInternalError(c *gin.Context, err error, context string) {
	log.Printf("[ERROR] %s: %v", context, err)
	c.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": "Internal server error"})
}
