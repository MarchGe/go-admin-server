package middleware

import (
	"errors"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"regexp"
)

const TranslatorName = "uni-translator"

func BindingValidateTranslator() gin.HandlerFunc {
	registerExtValidationFunc()
	translator := createTranslator()
	return func(c *gin.Context) {
		c.Set(TranslatorName, translator)
		c.Next()
	}
}

func registerExtValidationFunc() {
	v := ginUtils.GetValidator().V
	// 注册"label"标签
	v.RegisterTagNameFunc(registerLabelTagName())

	// 注册"regex"校验器--正则校验
	if err := v.RegisterValidation("regex", regexValidationFunc()); err != nil {
		panic(err)
	}
}

func createTranslator() ut.Translator {
	uniTranslator := ut.New(en.New(), zh.New())
	utTranslator, found := uniTranslator.GetTranslator("zh")
	if !found {
		panic(errors.New("uniTranslator.GetTranslator(\"zh\") failed"))
	}
	v := ginUtils.GetValidator().V
	err := zh2.RegisterDefaultTranslations(v, utTranslator)
	if err != nil {
		panic(err)
	}
	registerRegexTranslator(v, utTranslator)
	return utTranslator
}

func regexValidationFunc() validator.Func {
	return func(fl validator.FieldLevel) bool {
		regex := fl.Param()
		rgp, err := regexp.Compile(regex)
		if err != nil {
			panic(errors.New("Regexp expression error on field: " + fl.FieldName()))
		}
		return rgp.MatchString(fl.Field().String())
	}
}

func registerRegexTranslator(v *validator.Validate, utTranslator ut.Translator) {
	err := v.RegisterTranslation("regex", utTranslator, func(ut ut.Translator) error {
		return ut.Add("regex", "{0}的值不符合约束", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, e := ut.T("regex", fe.Field())
		if e != nil {
			panic(e)
		}
		return t
	})
	if err != nil {
		panic(err)
	}
}

func registerLabelTagName() validator.TagNameFunc {
	return func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label != "" {
			return "【" + label + "】"
		}
		return field.Name
	}
}
