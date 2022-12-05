package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseCode struct {
	Code   int    `json:"code"`
	Status string `json:"name"`
	Desc   string `json:"description"`
}
type ResponseResult struct {
	ResponseCode ResponseCode `json:"responseCode"`
	Result       interface{}  `json:"result"`
}
type ResponseWithResult struct {
	ResponseWithResult interface{} `json:"responseWithResult"`
}

func responseWithJson(ctx *gin.Context, code int, status string, desc string, result interface{}) {
	response := ResponseWithResult{
		ResponseResult{
			ResponseCode: ResponseCode{
				Code:   code,
				Status: status,
				Desc:   desc,
			},
			Result: result,
		},
	}

	ctx.JSON(http.StatusOK, response)
}
