package model

import "time"

// OperationLog 操作日志
type OperationLog struct {
	BaseTenantModel
	UserID        ID        `gorm:"index" json:"userId"`
	Username      string    `gorm:"size:64" json:"username"`
	Method        string    `gorm:"size:10" json:"method"`
	Path          string    `gorm:"size:255" json:"path"`
	Query         string    `gorm:"type:text" json:"query"`
	Body          string    `gorm:"type:text" json:"body"`
	IP            string    `gorm:"size:64" json:"ip"`
	UserAgent     string    `gorm:"size:255" json:"userAgent"`
	Status        int       `json:"status"`
	Latency       int64     `json:"latency"` // 耗时(ms)
	ErrorMessage  string    `gorm:"type:text" json:"errorMessage"`
	OperationTime *DateTime `gorm:"index;<-:create" json:"operationTime"`
}

func (o *OperationLog) SetOperationTime() {
	now := DateTime(time.Now())
	o.OperationTime = &now
}
