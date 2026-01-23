package scope

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一 JSON 结构
type Response struct {
	Code    int         `json:"code"`           // 业务自定义状态码
	Data    interface{} `json:"data,omitempty"` // 数据内容
	Message string      `json:"message"`        // 消息提示
}

// PageResult 分页数据结构
type PageResult struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// 预定义业务码 (ERP 常用)
const (
	SuccessCode       = 0
	ErrorCode         = 70001 // 通用错误
	InvalidParamsCode = 70002 // 参数错误
)

// Result 统一调用入口
func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Data:    data,
		Message: msg,
	})
}

// Ok 成功返回
func Ok(c *gin.Context) {
	Result(SuccessCode, map[string]interface{}{}, "OK", c)
}

// OkWithData 成功返回带数据
func OkWithData(c *gin.Context, data interface{}) {
	Result(SuccessCode, data, "OK", c)
}

func Error(c *gin.Context) {
	Result(ErrorCode, nil, "ERROR", c)
}

func Fail(c *gin.Context, msg string) {
	Result(ErrorCode, nil, msg, c)
}

func FailWithCode(c *gin.Context, code int, msg string) {
	Result(code, nil, msg, c)
}
