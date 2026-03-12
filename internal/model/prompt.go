package model

import "time"

// Prompt 用于保存各功能模块使用的 AI 提示词模板。
type Prompt struct {
	// ID 使用数据库自增主键。
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// Name 是提示词名称，便于后台识别。
	Name string `gorm:"size:128;not null" json:"name"`

	// Content 是提示词正文内容。
	Content string `gorm:"type:longtext;not null" json:"content"`

	// CreatedAt 是创建时间。
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt 是更新时间。
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Prompt) TableName() string {
	return "prompts"
}
