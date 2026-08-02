package response

import (
	"encoding/json"
	"net/http"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`           // 状态码
	Message string      `json:"message"`        // 消息
	Data    interface{} `json:"data,omitempty"` // 数据
}

// 常用状态码常量
const (
	CodeSuccess       = 0   // 成功
	CodeBadRequest    = 400 // 参数错误
	CodeUnauthorized  = 401 // 未授权
	CodeNotFound      = 404 // 未找到
	CodeInternalError = 500 // 内部错误
)

// 预定义消息映射（可选）
var codeMsgMap = map[int]string{
	CodeSuccess:       "success",               // 成功
	CodeBadRequest:    "bad request",           // 请求错误
	CodeUnauthorized:  "unauthorized",          // 未授权
	CodeNotFound:      "not found",             // 未找到
	CodeInternalError: "internal server error", // 内部错误
}

func writeJSON(w http.ResponseWriter, statusCode int, resp Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8") // 设置响应头
	w.WriteHeader(statusCode)                                         // 设置响应状态码
	json.NewEncoder(w).Encode(resp)                                   // 编码响应
}

func Success(w http.ResponseWriter, data interface{}) {
	writeJSON(w, http.StatusOK, Response{ // 写入响应
		Code:    CodeSuccess,             // 成功
		Message: codeMsgMap[CodeSuccess], // 成功消息
		Data:    data,                    // 数据
	})
}

func Fail(w http.ResponseWriter, httpStatus int, code int, msg ...string) {
	message := codeMsgMap[code] // 获取状态码对应的消息
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0] // 如果msg有值，则使用msg的值
	}
	writeJSON(w, httpStatus, Response{ // 写入响应
		Code:    code,
		Message: message,
	})
}

type PageData struct {
	List      interface{} `json:"list"`      // 列表
	Total     int64       `json:"total"`     // 总条数
	PageIndex int         `json:"pageIndex"` // 页码
	PageSize  int         `json:"pageSize"`  // 每页条数
}
