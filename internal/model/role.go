package model

type Role struct {
	BaseTenantModel
	Name        string  `gorm:"type:varchar(50);not null;index:idx_tenant_name" json:"name"`
	Description *string `gorm:"type:varchar(255)" json:"description"`

	Users []*User `gorm:"many2many:user_role;" json:"users,omitempty"`
	//关联权限
	Permissions []*Permission `gorm:"many2many:role_permission;" json:"permissions,omitempty"`

	//数据传输用，不处理数据
	PermissionIds *[]ID `gorm:"-" json:"permissionIds,omitempty"`
}

func (Role) TableName() string {
	return "role"
}
