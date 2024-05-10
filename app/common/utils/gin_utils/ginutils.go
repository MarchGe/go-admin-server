package gin_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"strconv"
)

func GetStringQuery(ctx *gin.Context, key string, defaultValue string) string {
	v := ctx.Query(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func GetInt8Query(c *gin.Context, key string, defaultValue int8) (int8, error) {
	v := c.Query(key)
	if v == "" {
		return defaultValue, nil
	}
	vv, err := strconv.ParseInt(v, 10, 8)
	if err != nil {
		return defaultValue, err
	}
	return int8(vv), nil
}
func GetInt16Query(c *gin.Context, key string, defaultValue int16) (int16, error) {
	v := c.Query(key)
	if v == "" {
		return defaultValue, nil
	}
	vv, err := strconv.ParseInt(v, 10, 16)
	if err != nil {
		return defaultValue, err
	}
	return int16(vv), nil
}
func GetInt32Query(c *gin.Context, key string, defaultValue int32) (int32, error) {
	v := c.Query(key)
	if v == "" {
		return defaultValue, nil
	}
	vv, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return defaultValue, err
	}
	return int32(vv), nil
}
func GetInt64Query(c *gin.Context, key string, defaultValue int64) (int64, error) {
	v := c.Query(key)
	if v == "" {
		return defaultValue, nil
	}
	vv, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue, err
	}
	return vv, nil
}

func GetIntQuery(c *gin.Context, key string, defaultValue int) (int, error) {
	v := c.Query(key)
	if v == "" {
		return defaultValue, nil
	}
	vv, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue, err
	}
	return vv, nil
}

type Validator struct {
	V *validator.Validate
}

func GetValidator() *Validator {
	v := binding.Validator.Engine().(*validator.Validate)
	return &Validator{
		V: v,
	}
}

func (v *Validator) Struct(s any) error {
	return v.V.Struct(s)
}
