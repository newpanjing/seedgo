package model

import (
	"database/sql/driver"
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

// Searchable 支持关键字搜索的接口
type Searchable interface {
	SearchFields() []string
}

// TenantModel 租户类
type TenantModel struct {
	//修改禁止更新
	TenantID ID `gorm:"column:tenant_id;<-:create" json:"tenantId"`
}

// DateTime 自定义时间类型，用于处理前端传递的时间字符串, 返回毫秒级时间戳，由前端解析成日期时间字符串，避免跨时区问题
type DateTime time.Time

const TimeFormat = "2006-01-02 15:04:05"

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	// 如果是数字，认为是毫秒时间戳
	if timestamp, err := strconv.ParseInt(string(data), 10, 64); err == nil {
		*t = DateTime(time.UnixMilli(timestamp))
		return nil
	}

	if len(data) == 2 {
		*t = DateTime(time.Time{})
		return
	}
	now, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	*t = DateTime(now)
	return
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).UnixMilli(), 10)), nil
}

func (t DateTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return nil, nil
	}
	return time.Time(t), nil
}

func (t *DateTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = DateTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t DateTime) String() string {
	return time.Time(t).Format(TimeFormat)
}

type BaseModel struct {
	ID        ID             `gorm:"primarykey" json:"id"`
	CreatedAt *DateTime      `gorm:"index;<-:create" json:"createdAt"`
	UpdatedAt *DateTime      `gorm:"index" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type BaseTenantModel struct {
	BaseModel
	TenantModel
}

// FieldTenantID 租户ID 字段名
const FieldTenantID = "TenantID"
