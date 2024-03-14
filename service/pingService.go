package service

import (
	"LarsWebV0/model"
	"github.com/gin-gonic/gin"
)

// @Summary ping
// @Tags ping
// @version 1.0
// @Accept application/x-json-stream
// @Router /ping [get]
func Ping(context *gin.Context) {
	response := model.Response{Context: context}
	response.Success(nil)
}
