package public

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) Use(g *gin.RouterGroup) {

}
