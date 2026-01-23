package model

import (
	"time"
)

type User struct {
	BaseTenantModel
	Username     string     `gorm:"type:varchar(50);not null;uniqueIndex:uniq_tenant_username" json:"username"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Phone        *string    `gorm:"type:varchar(20);index:idx_phone" json:"phone"`
	Email        *string    `gorm:"type:varchar(100)" json:"email"`
	RealName     *string    `gorm:"type:varchar(50)" json:"realName"`
	IsMain       *int8      `gorm:"type:tinyint;not null;default:0;index:idx_main" json:"isMain"`
	IsSuper      *bool      `gorm:"type:tinyint;not null;default:0" json:"isSuper"`
	Status       *int8      `gorm:"type:tinyint;not null;default:1;index:idx_status" json:"status"`
	LastLoginAt  *time.Time `json:"lastLoginAt"`
	LastLoginIP  *string    `gorm:"type:varchar(50)" json:"lastLoginIP"`

	Roles []*Role `gorm:"many2many:user_role;" json:"roles"`

	// 接收参数用
	RoleIds *[]ID `gorm:"-" json:"roleIds,omitempty"`

	//租户名称
	// gorm:"foreignKey:TenantID" 明确指定外键
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant"`
}

func (u *User) TableName() string {
	return "user"
}

// SearchFields 返回搜索字段
func (u *User) SearchFields() []string {
	return []string{"username", "realName", "phone", "email"}
}
