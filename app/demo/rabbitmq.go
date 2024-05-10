package demo

import (
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/demo/mq"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

var rabbitApi = &RabbitApi{}

type RabbitApi struct {
}

func GetRabbitApi() *RabbitApi {
	return rabbitApi
}

func (a *RabbitApi) SendMsg(c *gin.Context) {
	msg := &mq.Message{}
	if err := c.ShouldBindJSON(msg); err != nil {
		E.PanicErr(err)
	}

	err := mq.SendMessageConfirm(msg)
	if err != nil {
		slog.Error("send rabbitmq message error", slog.Any("err", err))
		R.Fail(c, constant.ServerInternalError, http.StatusInternalServerError)
		return
	}
	R.Success(c, nil)
}
