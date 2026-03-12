package bootstrap

import (
	"fmt"
)

type Hooks struct {
	LoadConfig func() error
	InitDB     func() error
}

func Init() error {
	return InitWith(Hooks{
		LoadConfig: LoadConfigFile,
		InitDB:     InitGorm,
	})
}

func InitWith(h Hooks) error {
	if h.LoadConfig == nil {
		return fmt.Errorf("bootstrap hook LoadConfig is nil")
	}
	if h.InitDB == nil {
		return fmt.Errorf("bootstrap hook InitDB is nil")
	}

	if err := h.LoadConfig(); err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}
	if err := h.InitDB(); err != nil {
		return fmt.Errorf("init gorm failed: %w", err)
	}
	return nil
}
