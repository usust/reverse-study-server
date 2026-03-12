package system

const currentVersion = "v0.1"

// VersionInfo 是版本接口的响应结构。
//
// 当前只返回版本号，后续如果需要追加构建时间、Git 提交哈希等信息，
// 可以直接扩展这个结构，而不需要改路由签名。
type VersionInfo struct {
	Version string `json:"version"`
}

// GetVersionInfo 返回当前服务的版本信息。
//
// 目前版本号先固定为 v0.1，便于先完成接口联调。
// 后续如果改为从配置、构建参数或发布流水线注入，优先修改这里。
func GetVersionInfo() VersionInfo {
	return VersionInfo{
		Version: currentVersion,
	}
}
