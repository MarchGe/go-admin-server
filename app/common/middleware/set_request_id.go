package middleware

import (
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.NewString()
		c.Request.Header.Set(constant.RequestId, requestId)
		c.Set(constant.RequestId, requestId)
		c.Next()
	}
}
