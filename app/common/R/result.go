package R

import (
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	RequestId string `json:"requestId"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Value     any    `json:"value"`
}

const (
	SuccessKey = "_c_success"
)

func Success(c *gin.Context, value any) {
	c.Set(SuccessKey, true)
	c.AbortWithStatusJSON(http.StatusOK, Result{
		RequestId: c.GetString(constant.RequestId),
		Code:      http.StatusOK,
		Message:   "成功",
		Value:     value,
	})
}

func Fail(c *gin.Context, message string, code ...int) {
	rCode := http.StatusBadRequest
	length := len(code)
	if length > 0 {
		rCode = code[0]
	}
	c.AbortWithStatusJSON(http.StatusOK, Result{
		RequestId: c.GetString(constant.RequestId),
		Code:      rCode,
		Message:   message,
	})
}
