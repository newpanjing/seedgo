package request

import "seedgo/internal/model"

// UserContext 存放在 Context 中的用户信息对象
type UserContext struct {
	ID       model.ID        `json:"id"`
	Username string        `json:"username"`
	TenantID model.ID        `json:"tenantId"`
	IsSuper  bool          `json:"isSuper"`
	Roles    []*model.Role `json:"roles"`
}
