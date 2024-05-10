package recorder

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MarchGe/go-admin-server/app/admin/model"
	"github.com/MarchGe/go-admin-server/app/admin/service"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"strings"
	"time"
)

func RecordLoginLog(c *gin.Context, userId int64) {
	if !config.GetConfig().Log.LoginLog {
		return
	}
	logService := service.GetLogService()
	userAgent, clientIp := getInfoFromRequest(c)
	log := &model.LoginLog{
		UserAgent: userAgent,
		Ip:        clientIp,
	}
	u := getUserInfo(userId)
	log.UserId = u.Id
	log.Nickname = u.Nickname
	log.RealName = u.Name
	if dept := u.Dept; dept != nil {
		log.DeptName = dept.Name
	}
	log.CreateTime = time.Now()
	log.UpdateTime = time.Now()
	if err := logService.AddLoginLog(log); err != nil {
		slog.Error("-", slog.Any("err", err))
	}
}

// RecordOpLog opTarget是操作的对象；可变参数v中的第一个表示动作名称，第二个表示是否是敏感信息（日志中不会记录敏感参数），如果没有动作名称，也可以直接把第一个参数用作bool值，表示是否敏感信息
func RecordOpLog(opTarget string, v ...any) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.GetConfig().Log.OpLog {
			return
		}
		logService := service.GetLogService()
		userAgent, clientIp := getInfoFromRequest(c)
		log := &model.OpLog{
			UserAgent: userAgent,
			Ip:        clientIp,
		}
		userId := c.GetInt64(constant.SessionUserId)
		if userId != 0 {
			if user := getUserInfo(userId); user != nil {
				log.UserId = userId
				log.Nickname = user.Nickname
				log.RealName = user.Name
				if dept := user.Dept; dept != nil {
					log.DeptName = dept.Name
				}
			}
		}
		var private bool
		if len(v) == 0 {
			log.Action = parseMethod(c.Request.Method)
		} else if len(v) == 1 {
			b, ok := v[0].(bool)
			if ok {
				private = b
				log.Action = parseMethod(c.Request.Method)
			} else {
				log.Action = v[0].(string)
			}
		} else {
			log.Action = v[0].(string)
			private = v[1].(bool)
		}
		log.Path = c.Request.URL.Path
		log.Target = opTarget
		if private {
			log.Query = "[private]"
			log.Body = "[private]"
		} else {
			log.Query = c.Request.URL.Query().Encode()
			if len(log.Query) > 255 {
				log.Query = "[too long ignored]"
			}
			log.Body = getBodyParams(c)
			if len(log.Body) > 500 {
				log.Body = "[too long ignored]"
			}
		}
		log.CreateTime = time.Now()
		log.UpdateTime = time.Now()
		c.Next()
		if c.GetBool(R.SuccessKey) {
			if err := logService.AddOpLog(log); err != nil {
				slog.Error("-", slog.Any("err", err))
			}
		}
	}
}

func RecordExceptionLog(c *gin.Context, errString string) {
	if !config.GetConfig().Log.ExceptionLog {
		return
	}
	logService := service.GetLogService()
	userAgent, clientIp := getInfoFromRequest(c)
	log := &model.ExceptionLog{
		UserAgent: userAgent,
		Ip:        clientIp,
	}
	userId := c.GetInt64(constant.SessionUserId)
	if userId != 0 {
		if user := getUserInfo(userId); user != nil {
			log.UserId = userId
			log.Nickname = user.Nickname
		}
	}
	log.Path = c.Request.URL.Path
	log.Query = c.Request.URL.Query().Encode()
	if len(log.Query) > 255 {
		log.Query = "[too long ignored]"
	}
	log.Body = getBodyParams(c)
	if len(log.Body) > 500 {
		log.Body = "[too long ignored]"
	}
	log.Error = errString
	log.CreateTime = time.Now()
	log.UpdateTime = time.Now()
	if err := logService.AddExceptionLog(log); err != nil {
		slog.Error("-", slog.Any("err", err))
	}
}

func getBodyParams(c *gin.Context) string {
	requestBodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewReader(requestBodyBytes))
	return string(requestBodyBytes)
}

func parseMethod(method string) string {
	switch strings.ToLower(method) {
	case "get":
		return "查询"
	case "post":
		return "新增"
	case "put":
		return "修改"
	case "delete":
		return "删除"
	default:
		return "其他"
	}
}

func getInfoFromRequest(c *gin.Context) (userAgent, clientIp string) {
	userAgent = c.Request.UserAgent()
	clientIp = c.ClientIP()
	return
}

type cacheUserTemp struct {
	T time.Time   `json:"t"`
	U *model.User `json:"u"`
}

func getUserInfo(userId int64) *model.User {
	const cacheKey = "_curr_u_by_log"
	cache := utils.GetCache()
	userBytes, err := cache.Get(cacheKey)
	if err != nil {
		if errors.Is(err, utils.ErrCacheNotFound) {
			return getUserInfoAndSaveToCache(userId, cacheKey)
		}
		slog.Error("get current user temp info from cache error", slog.Any("err", err))
		return nil
	}
	cacheU := &cacheUserTemp{}
	if err = json.Unmarshal(userBytes, cacheU); err != nil {
		slog.Error("json unmarshal current user temp info error", slog.Any("err", err))
		return nil
	}

	if time.Now().Sub(cacheU.T) > 10*time.Second {
		return getUserInfoAndSaveToCache(userId, cacheKey)
	}
	return cacheU.U
}

func getUserInfoAndSaveToCache(userId int64, cacheKey string) *model.User {
	user := service.GetUserService().FindOneById(userId, "Dept")
	if user == nil {
		return nil
	}
	cacheU := &cacheUserTemp{
		T: time.Now(),
		U: user,
	}
	uBytes, err := json.Marshal(cacheU)
	if err != nil {
		slog.Error("json marshal current user temp info error", slog.Any("err", err))
		return nil
	}
	if err := utils.GetCache().Set(cacheKey, uBytes); err != nil {
		slog.Error("set current user temp info to cache error", slog.Any("err", err))
		return nil
	}
	return user
}
