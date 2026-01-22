package model

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// ID 通用ID类型
type ID int64

// ToID 创建一个方法将任意类型转换为id
func ToID(value any) ID {
	//转为字符串，再转为int64
	idStr := fmt.Sprintf("%v", value)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return ID(id)
}

func (id *ID) String() string {
	return fmt.Sprintf("%d", int64(*id))
}

type Entity interface {
	// GetPrimaryKey 获取主键名
	GetPrimaryKey() string
	// GetID 获取id
	GetID() ID
	// SetID 设置id
	SetID(id ID)
}

// TenantModel 租户类
type TenantModel struct {
	//修改禁止更新
	TenantID ID `gorm:"column:tenant_id;<-:create" json:"tenantId"`
}

type BaseModel struct {
	ID        ID             `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"index;<-:create" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"index" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type BaseTenantModel struct {
	BaseModel
	TenantModel
}

// FieldTenantID 租户ID 字段名
const FieldTenantID = "TenantID"
