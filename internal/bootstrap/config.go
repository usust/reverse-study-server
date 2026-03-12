package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const ConfigFileEnv = "DATA_SECURITY_LAB_CONFIG_FILE"

var GlobalConfig Config
var ZapLogger = zap.NewNop()
var SugaredLogger = ZapLogger.Sugar()

func buildViper(configFile string) *viper.Viper {
	v := viper.New()
	if configFile != "" {
		v.SetConfigFile(configFile)
		return v
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./internal/config")
	v.AddConfigPath(filepath.Join("..", "internal", "config"))
	return v
}

// LoadConfigFile 加载配置文件
func LoadConfigFile() error {
	configFile := strings.TrimSpace(os.Getenv(ConfigFileEnv))
	v := buildViper(configFile)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read config file failed (env %s=%q): %w", ConfigFileEnv, configFile, err)
	}
	fmt.Printf("load config file : %s\n", v.ConfigFileUsed())

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshal config file failed: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	GlobalConfig = cfg

	if err := InitLogger(); err != nil {
		return fmt.Errorf("init logger failed: %w", err)
	}

	return nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.ZapLog.LogDir) == "" {
		return fmt.Errorf("zap_log.log_dir is required")
	}

	driver := strings.ToLower(strings.TrimSpace(c.DB.Driver))
	if driver == "" {
		return fmt.Errorf("db.driver is required")
	}

	switch driver {
	case "mysql":
		if strings.TrimSpace(c.Mysql.Path) == "" {
			return fmt.Errorf("mysql.path is required")
		}
		if c.Mysql.Port <= 0 || c.Mysql.Port > 65535 {
			return fmt.Errorf("mysql.port must be in range 1-65535")
		}
		if strings.TrimSpace(c.Mysql.Username) == "" {
			return fmt.Errorf("mysql.username is required")
		}
		if strings.TrimSpace(c.Mysql.Database) == "" {
			return fmt.Errorf("mysql.database is required")
		}

		parseTime := strings.ToLower(strings.TrimSpace(c.Mysql.ParseTime))
		if parseTime != "true" && parseTime != "false" {
			return fmt.Errorf("mysql.parse_time must be true/false")
		}
	default:
		return fmt.Errorf("unsupported db.driver: %s", c.DB.Driver)
	}

	return nil
}
