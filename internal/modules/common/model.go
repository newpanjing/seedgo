package common

// SelectOption 下拉框选项
type SelectOption struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"` // 使用 interface{} 以支持 int 或 string 类型的 ID
}
