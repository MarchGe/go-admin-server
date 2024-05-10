package apis

import (
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/req"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _settingsApi = &SettingsApi{
	settingsService: service.GetSettingsService(),
}

type SettingsApi struct {
	settingsService *service.SettingsService
}

func GetSettingsApi() *SettingsApi {
	return _settingsApi
}

// Upsert godoc
//
//	@Summary	保存/更新系统配置
//	@Tags		系统配置
//	@Accept		application/json
//	@Produce	application/json
//
//	@Param		[body]	body		req.SettingsUpsertReq	true	"配置信息"
//
//	@Success	200		{object}	R.Result
//	@Router		/settings [put]
func (a *SettingsApi) Upsert(c *gin.Context) {
	m := &req.SettingsUpsertReq{}
	if err := c.ShouldBindJSON(m); err != nil {
		E.PanicErr(err)
	}
	userId := c.GetInt64(constant.SessionUserId)
	err := a.settingsService.UpsertSettings(userId, m)
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

// GetSettings godoc
//
//	@Summary	查询系统信息
//	@Tags		系统配置
//	@Produce	application/json
//
//	@Param		key	query		string	true	"键名"
//
//	@Success	200	{object}	R.Result{value=model.Settings}
//	@Router		/settings [get]
func (a *SettingsApi) GetSettings(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		R.Fail(c, "键名不能为空", http.StatusBadRequest)
		return
	}
	settings, _ := a.settingsService.FindOneByUserIdAndKey(c.GetInt64(constant.SessionUserId), key)
	R.Success(c, settings)
}
