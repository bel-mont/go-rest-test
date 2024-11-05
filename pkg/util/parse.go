package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseID(c *gin.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("invalid Id format")
	}
	return id, nil
}
