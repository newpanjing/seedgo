package model

type Permission struct {
	ID             ID            `gorm:"primarykey" json:"id"`
	ParentID       *ID           `gorm:"index" json:"parentId"`
	Name           string        `gorm:"index" json:"name"`
	Path           string        `gorm:"index" json:"path"`
	Icon           string        `gorm:"index" json:"icon"`
	PermissionCode string        `gorm:"index" json:"permissionCode"`
	Sort           *int          `gorm:"index" json:"sort"`
	Visible        *bool         `gorm:"index;default:true" json:"visible"`
	Type           *int          `gorm:"index" json:"type"`
	PermissionUrls string        `gorm:"column:permission_urls" json:"permissionUrls"` // 多个Url用英文逗号分割
	Children       []*Permission `gorm:"-" json:"children,omitempty"`
}

func (p *Permission) TableName() string {
	return "permission"
}
