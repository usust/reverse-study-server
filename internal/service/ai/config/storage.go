package config

import (
	"context"
	"fmt"
	"reverse-study-server/internal/bootstrap"
	"strings"
)

// StorageConfig 是通用存储配置响应结构。
type StorageConfig struct {
	BaseDir string `json:"base_dir"`
}

// GetStorageConfig 获取当前内存中的存储配置。
func GetStorageConfig(_ context.Context) (StorageConfig, error) {
	return StorageConfig{
		BaseDir: strings.TrimSpace(bootstrap.GlobalConfig.Saved.BaseDir),
	}, nil
}

// UpdateStorageConfig 更新当前内存中的存储配置。
func UpdateStorageConfig(_ context.Context, baseDir string) (StorageConfig, error) {
	trimmed := strings.TrimSpace(baseDir)
	if trimmed == "" {
		return StorageConfig{}, fmt.Errorf("baseDir is required")
	}

	bootstrap.GlobalConfig.Saved.BaseDir = trimmed
	return StorageConfig{BaseDir: trimmed}, nil
}
