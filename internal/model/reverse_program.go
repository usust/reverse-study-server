package model

import "time"

// ReverseProgram 记录一次生成的逆向程序源码。
type ReverseProgram struct {
	// ID 是程序记录主键，使用数据库自增数字。
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// Title 是题目名称。
	Title string `gorm:"size:255;not null;default:''" json:"title"`
	// Description 是题目描述。
	Description string `gorm:"type:text" json:"description"`
	// Published 表示发布状态，0=未发布，1=已发布。
	Published int `gorm:"not null;default:0" json:"published"`

	// Score 是题目分值。
	Score int `gorm:"not null;default:100" json:"score"`
	// ProgramType 是题目类型。
	ProgramType int `gorm:"not null;default:0" json:"programType"`
	// Difficulty 是题目难度。
	Difficulty int `gorm:"not null;default:0" json:"difficulty"`
	// Tags 是题目标签，支持多个标签。
	Tags []string `gorm:"serializer:json" json:"tags"`

	// CompletedCount 是完成数量。
	CompletedCount int `gorm:"not null;default:0" json:"completedCount"`

	// BaseDir 是下载路径或保存路径标识。
	BaseDir string `gorm:"size:512;not null" json:"baseDir"`
	// SourceFileName 是源码文件名，默认 main.c。
	SourceFileName string `gorm:"size:128;not null;default:main.c" json:"sourceFileName"`
	// MetaFileName 编译的元信息
	MetaFileName string `gorm:"size:128;not null;default:main.meta" json:"metaFileName"`
	// ProgramFileName 编译后文件名
	ProgramFileName string `gorm:"size:128;not null;default:main.c" json:"programFileName"`
	// ProgramFileMD5 是程序源码文件的 MD5。
	ProgramFileMD5 string `gorm:"size:64;not null;index" json:"programFileMd5"`

	// CreatedAt 是创建时间。
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt 是更新时间。
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName 指定表名。
func (ReverseProgram) TableName() string {
	return "reverse_programs"
}
