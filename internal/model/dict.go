package model

import "time"

// Dict 字典表
type Dict struct {
	ID          ID        `gorm:"primaryKey;autoIncrement" json:"id"`
	Code        string    `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description *string   `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time `gorm:"index;<-:create" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Items []*DictItem `gorm:"foreignKey:DictID" json:"items"`
}

func (d *Dict) TableName() string {
	return "dict"
}

// SearchFields 支持搜索的字段
func (d *Dict) SearchFields() []string {
	return []string{"code", "name"}
}

// DictItem 字典项表
type DictItem struct {
	ID        ID        `gorm:"primaryKey;autoIncrement" json:"id"`
	DictID    ID        `gorm:"index;not null" json:"dictId"`
	Label     string    `gorm:"type:varchar(100);not null" json:"label"`
	Value     string    `gorm:"type:varchar(100);index;not null" json:"value"`
	Sort      int       `gorm:"default:0;not null" json:"sort"`
	Status    int       `gorm:"type:tinyint;default:1;index;not null" json:"status"`
	CreatedAt time.Time `gorm:"index;<-:create" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	Dict *Dict `gorm:"foreignKey:DictID" json:"-"`
}

func (d *DictItem) TableName() string {
	return "dict_item"
}

// SearchFields 支持搜索的字段
func (d *DictItem) SearchFields() []string {
	return []string{"label", "value"}
}
