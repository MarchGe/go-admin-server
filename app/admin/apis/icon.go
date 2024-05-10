package apis

import (
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/gin-gonic/gin"
)

var _iconApi = &IconApi{
	iconService: service.GetIconService(),
}

type IconApi struct {
	iconService *service.IconService
}

func GetIconApi() *IconApi {
	return _iconApi
}

// GetAll godoc
//
//	@Summary	查询全部图标
//	@Tags		图标管理
//	@Produce	application/json
//	@Success	200	{object}	R.Result{value=model.Icon}
//	@Router		/icon/all [get]
func (a *IconApi) GetAll(c *gin.Context) {
	icons, err := a.iconService.FindAll()
	if err != nil {
		E.PanicErr(err)
	}
	R.Success(c, icons)
}
