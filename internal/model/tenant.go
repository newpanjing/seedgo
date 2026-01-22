package model

type Tenant struct {
	BaseModel
	Name         string `gorm:"size:255;not null" json:"name"`
	ContactName  string `gorm:"size:50" json:"contactName"`
	ContactPhone string `gorm:"size:20" json:"contactPhone"`
	ContactEmail string `gorm:"size:100" json:"contactEmail"`
	Status       int    `gorm:"default:1" json:"status"`
	Users        []User `gorm:"foreignKey:TenantID;references:ID" json:"users"`

	// 接收参数用
	Username string `gorm:"-" json:"username,omitempty"`
	Password string `gorm:"-" json:"password,omitempty"`
	RealName string `gorm:"-" json:"realName,omitempty"`
	Phone    string `gorm:"-" json:"phone,omitempty"`
}

func (Tenant) TableName() string {
	return "tenant"
}
