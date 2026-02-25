package handler

import (
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
