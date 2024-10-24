package utils

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type ValidationPair struct {
	Field interface{}
	Tag   string
	Err   string
}

func ValidatePairs(v *validator.Validate, pairs []ValidationPair) error {
	for _, p := range pairs {
		if err := v.Var(p.Field, p.Tag); err != nil {
			for _, errs := range err.(validator.ValidationErrors) {
				return fmt.Errorf("%s %s", p.Err, errs.Value())
			}
		}
	}
	return nil
}

// MyValidator 自定义参数检验结构体
type MyValidator struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

// 自定义校验函数
func requiredAllowZero(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().Uint()
	return fieldValue == 0 || fieldValue != 0
}

// IniValidator 初始化自定义校验, 使用中文报错信息
func IniValidator() (m *MyValidator) {
	var _m MyValidator
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")

	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})

	// 注册自定义校验函数
	validate.RegisterValidation("requiredAllowZero", requiredAllowZero)

	zh_translations.RegisterDefaultTranslations(validate, trans)

	_m.Validate = validate
	_m.Trans = trans

	return &_m
}

// Check 参数校验
func (m *MyValidator) Check(value interface{}) error {
	if value == nil {
		return errors.New("验证的值不能为nil")
	}

	err := m.Validate.Struct(value)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 当转换失败时，也尽可能提供详细的错误信息
			return errors.New("验证异常: " + err.Error())
		}

		// 收集所有错误信息
		errBuf := &bytes.Buffer{}
		for _, err := range errs {
			errBuf.WriteString(err.Translate(m.Trans))
			errBuf.WriteString(", ")
		}
		// 移除最后一个逗号和空格
		if errBuf.Len() > 0 {
			errBuf.Truncate(errBuf.Len() - 2)
		}

		return errors.New(errBuf.String())
	}

	return nil
}
