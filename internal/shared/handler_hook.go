package shared

import "github.com/gin-gonic/gin"

type HandlerHook[T any] interface {
	BeforePage(ctx *gin.Context)
	AfterPage(ctx *gin.Context, items []T, total int64)
}

// 默认实现Hook接口
type DefaultHandlerHook[T any] struct{}

func (h *DefaultHandlerHook[T]) BeforePage(ctx *gin.Context) {
	//默认实现为空
}

func (h *DefaultHandlerHook[T]) AfterPage(ctx *gin.Context, items []T, total int64) {
	//默认实现为空
}
