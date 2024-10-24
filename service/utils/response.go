package utils

import (
	"strconv"
)

/**
 * <h1>通用响应结构体定义</h1>
 */
type CommonResponse struct {
	Code    int         `json:"code"`    // 响应状态：0成功 其他失败
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

/*
* <h2>响应成功</h2>
* @param data 响应数据
**/
func ResponseSuccess(data interface{}) CommonResponse {
	return CommonResponse{
		Code:    0,
		Message: "操作成功",
		Data:    data,
	}
}

/*
* <h2>响应失败</h2>
* @param code 响应异常状态
* @param message 响应异常消息
**/
func ResponseError(code int, message string) CommonResponse {
	return CommonResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

/*
* <h2>响应失败</h2>
* @param code 响应异常状态
**/
func ResponseErrorCode(code int) CommonResponse {
	return CommonResponse{
		Code:    code,
		Message: "操作失败",
		Data:    nil,
	}
}

func ResponseErrorMessage(message string) CommonResponse {

	// 异常消息长度大于4
	if len(message) > 4 {
		// 截取异常消息前四位
		var codeStr = message[0:4]
		// 前四位字符串类型转 int 类型
		code, err := strconv.Atoi(codeStr)
		// 如果没有异常，自动分割，示例：0000_XXXXXXX
		if err == nil {
			return CommonResponse{
				Code:    code,
				Message: message[5:],
				Data:    nil,
			}
		}
	}

	return CommonResponse{
		Code:    -1,
		Message: message,
		Data:    nil,
	}
}
