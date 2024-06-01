package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func ParseValidateError(serviceErr error) error {
	//断言错误
	var err validator.ValidationErrors

	if errors.As(serviceErr, &err) {
		//错误信息
		var errMsg string
		//遍历错误
		for _, e := range err {
			//field
			field := e.Field()
			//tag
			tag := e.Tag()
			//param
			param := e.Param()

			if tag == "required" {
				errMsg = fmt.Sprintf("%s不能为空", field)
			} else if tag == "min" {
				errMsg = fmt.Sprintf("%s的长度不能小于%s", field, param)
			} else if tag == "max" {
				errMsg = fmt.Sprintf("%s的长度不能大于%s", field, param)
			} else if tag == "email" {
				errMsg = fmt.Sprintf("%s的值不是有效的邮箱格式", field)
			} else if tag == "alphanumunderscore" {
				errMsg = fmt.Sprintf("%s的值只能包括字母、数字和下划线", field)
			} else {
				errMsg = fmt.Sprintf("%s验证失败", field)
			}
		}
		return errors.New(errMsg)
	}
	return serviceErr
}
