package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Context *gin.Context
}

func (r *Response) Success(data interface{}) {
	r.Context.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	})
}

func (r *Response) Fails(message string, err error) {
	r.Context.JSON(http.StatusInternalServerError, gin.H{
		"code":    http.StatusInternalServerError,
		"message": fmt.Sprintf(message+": %v", err),
	})
}

func (r *Response) Redirect(redirectUrl string) {
	r.Context.Redirect(http.StatusFound, redirectUrl)
}
