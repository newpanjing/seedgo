package tenant

import (
	"github.com/newpanjing/seedgo/internal/model"
	"github.com/newpanjing/seedgo/internal/shared"
)

type TenantCtrl struct {
	shared.BaseCtrl[model.Tenant]
	logic *TenantLogic
}

func NewTenantCtrl() *TenantCtrl {
	logic := NewTenantLogic()
	ctr := &TenantCtrl{
		logic: logic,
	}
	ctr.BaseCtrl = *shared.NewBaseCtrl(logic, nil, ctr)
	return ctr
}
