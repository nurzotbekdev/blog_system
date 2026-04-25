package helper

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseUintParam(ctx *gin.Context, name string) (uint, error) {
	param := ctx.Param(name)
	if param == "" {
		return 0, errors.New("Request parameter has invalid format")
	}

	value, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, errors.New("Request parameter has invalid format")
	}

	return uint(value), nil
}
