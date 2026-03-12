package model

import (
	"time"
)

// ModelAPIModel 是模型接入配置的持久化模型。
//
// 它直接对齐前端“多模型接入配置”的核心字段，并补充时间戳字段，
// 便于后续做审计、排序和配置管理。
type ModelAPIModel struct {
	// ID 是这条模型配置的主键。
	//
	// 设计上优先复用前端保存的配置 ID，便于前后端直接按同一个标识同步；
	// 如果前端未传，当前服务层会在落库前自动生成 UUID。
	//
	// gorm 约束：
	// 1. 作为主键；
	// 2. 最大长度 64；
	// 3. 不依赖数据库自增，便于跨环境迁移和前端先生成 ID。
	ID string `gorm:"primaryKey;size:64" json:"id"`

	// Name 是这条模型接入配置的显示名称。
	//
	// 典型值例如：
	// 1. “火山推理-生产”
	// 2. “本地兼容网关-测试”
	//
	// 该字段主要用于前端列表展示、人工识别和后台管理。
	Name string `gorm:"size:128;not null" json:"name"`

	// Provider 表示模型服务供应商名称。
	//
	// 它用于标记这条配置属于哪一类接入来源，例如：
	// 1. Volcengine Ark
	// 2. OpenAI Compatible
	// 3. Local Gateway
	//
	// 当前字段是普通字符串，不做枚举限制，目的是保持扩展弹性。
	Provider string `gorm:"size:128;not null" json:"provider"`

	// BaseURL 是模型服务的基础地址。
	//
	// 对于需要走代理层、内网网关或兼容 OpenAI 协议的服务，
	// 这个字段可直接存储对应入口地址。
	//
	// 当前火山 SDK 调用链尚未实际消费该字段，但先持久化下来，
	// 可以避免前端配置丢失，也为后续多供应商接入预留空间。
	BaseURL string `gorm:"size:512" json:"baseUrl"`

	// APIKey 是访问模型服务使用的凭证。
	//
	// 当前直接以字符串形式存储，方便快速联调。
	// 如果后续进入更正式环境，建议改为：
	// 1. 加密存储；
	// 2. 只保存引用或密钥别名；
	// 3. 与外部密钥管理系统集成。
	APIKey string `gorm:"size:512" json:"apiKey"`

	// Model 是本条配置默认使用的模型标识。
	//
	// 该字段更偏向“前端展示使用的模型名称”，例如：
	// 1. Doubao-Seed-2.0-Code
	// 2. Doubao-Pro
	//
	// 实际 API 调用时会优先使用 APIModel；
	// 当 APIModel 为空时，才会回退到 Model。
	Model string `gorm:"size:128;not null" json:"model"`

	// APIModel 是实际提交给模型 API 的模型 ID。
	//
	// 例如：
	// 1. doubao-seed-2-0-code-preview-260215
	// 2. doubao-pro-32k-241215
	//
	// 当该字段为空时，调用层会回退使用 Model 字段，兼容旧数据。
	APIModel string `gorm:"size:128" json:"apiModel"`

	// Enabled 表示该配置当前是否启用。
	//
	// 这是一个轻量级开关，用于在不删除配置的前提下临时停用某条接入，
	// 避免误调用、误计费，或者在故障排查时快速摘除有问题的配置。
	Enabled bool `gorm:"not null;default:false" json:"enabled"`

	// CreatedAt 是配置首次入库时间。
	//
	// 由 GORM 自动维护，可用于：
	// 1. 后台按创建时间排序；
	// 2. 审计某条配置的首次接入时间；
	// 3. 未来做“最近新增配置”筛选。
	CreatedAt time.Time `json:"createdAt"`
	// UpdatedAt 是配置最后一次更新的时间。
	//
	// 同样由 GORM 自动维护，可用于判断配置是否被最近修改、
	// 是否需要重新校验，或展示“最后编辑时间”。
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName 显式指定数据库表名。
//
// 显式声明的原因：
// 1. 避免后续结构体改名时影响表名；
// 2. 保持数据库命名稳定，方便写 SQL、排查问题和做迁移；
// 3. 与“模型接入配置”这一业务概念保持一一对应。
func (ModelAPIModel) TableName() string {
	return "model_api_configs"
}
