package tenant

import (
	"seedgo/internal/model"
	"seedgo/internal/shared"
)

type Handler struct {
	shared.BaseHandler[model.Tenant]
	logic *TenantLogic
}

func NewHandler() *Handler {
	logic := NewTenantLogic()
	ctr := &Handler{
		logic: logic,
	}
	ctr.BaseHandler = *shared.NewBaseHandler(logic, nil, ctr)
	return ctr
}
