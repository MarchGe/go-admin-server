package middleware

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/middleware/recorder"
	"github.com/MarchGe/go-admin-server/config"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

func GlobalErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				cfg := config.GetConfig()
				if _, ok := err.(*E.ApplicationError); ok {
					appErr := err.(*E.ApplicationError)
					var validationErrors validator.ValidationErrors
					if appErr.Code == E.NormalErrCode {
						R.Fail(c, appErr.Err.Error(), http.StatusBadRequest)
						return
					} else if errors.As(appErr.Err, &validationErrors) {
						translator := c.MustGet(TranslatorName).(ut.Translator)
						R.Fail(c, joinErrors(validationErrors, translator), http.StatusBadRequest)
						return
					}
					args := []any{
						slog.String("requestId", c.GetString(constant.RequestId)),
						slog.Any("err", appErr.Err),
					}
					if cfg.Log.StackTrace {
						args = append(args, slog.String("error stack", string(debug.Stack())))
					}
					slog.Error(constant.ServerInternalError, args...)
					if cfg.Environment == config.DEV {
						_, _ = os.Stderr.Write(debug.Stack())
					}
					recorder.RecordExceptionLog(c, fmt.Sprintf("requestId=%s", c.GetString(constant.RequestId)))
					R.Fail(c, constant.ServerInternalError, http.StatusInternalServerError)
					return
				}
				args := []any{
					slog.String("requestId", c.GetString(constant.RequestId)),
					slog.Any("err", err),
				}
				if cfg.Log.StackTrace {
					args = append(args, slog.String("error stack", string(debug.Stack())))
				}
				slog.Error(constant.ServerInternalError, args...)
				if cfg.Environment == config.DEV {
					_, _ = os.Stderr.Write(debug.Stack())
				}
				recorder.RecordExceptionLog(c, fmt.Sprintf("requestId=%s", c.GetString(constant.RequestId)))
				R.Fail(c, constant.ServerInternalError, http.StatusInternalServerError)
				return
			}
		}()
		c.Next()
	}
}

func joinErrors(errs validator.ValidationErrors, translator ut.Translator) string {
	errStrings := make([]string, len(errs))
	for i, item := range errs {
		errStrings[i] = item.Translate(translator)
	}
	return strings.Join(errStrings, ",")
}
