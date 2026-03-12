package bootstrap

import (
	"fmt"
	internalmodel "reverse-study-server/internal/model"
	"reverse-study-server/volcengine/model"
	"strings"
	"time"

	gosqlmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormDB global gorm数据库对象
var GormDB *gorm.DB

// InitGorm 初始化 Gorm 数据库连接。
func InitGorm() error {
	driver := strings.ToLower(strings.TrimSpace(GlobalConfig.DB.Driver))
	var dialector gorm.Dialector

	switch driver {
	case "mysql":
		mysqlCfg := GlobalConfig.Mysql

		// 使用 go-sql-driver/mysql 的配置对象来生成 DSN，而不是手工拼接字符串。
		//
		// 这样可以自动处理：
		// 1. loc 参数中的特殊字符转义（例如 Asia/Shanghai 里的 `/`）；
		// 2. 用户名、密码中的特殊字符；
		// 3. 参数格式化的一致性。
		//
		// 手工拼接在 loc=Asia/Shanghai 这类场景下很容易触发：
		// "invalid DSN: did you forget to escape a param value?"
		driverCfg := gosqlmysql.Config{
			User:                 mysqlCfg.Username,
			Passwd:               mysqlCfg.Password,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%d", mysqlCfg.Path, mysqlCfg.Port),
			DBName:               mysqlCfg.Database,
			Params:               map[string]string{"charset": mysqlCfg.Charset},
			ParseTime:            strings.EqualFold(strings.TrimSpace(mysqlCfg.ParseTime), "true"),
			Loc:                  time.Local,
			AllowNativePasswords: true,
		}

		if loc := strings.TrimSpace(mysqlCfg.Loc); loc != "" && !strings.EqualFold(loc, "local") {
			loadedLoc, err := time.LoadLocation(loc)
			if err != nil {
				return fmt.Errorf("load mysql location %q failed: %w", mysqlCfg.Loc, err)
			}
			driverCfg.Loc = loadedLoc
		}

		dsn := driverCfg.FormatDSN()
		dialector = mysql.Open(dsn)
	default:
		return fmt.Errorf("unsupported db.driver: %s", GlobalConfig.DB.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return fmt.Errorf("open gorm failed: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping db failed: %w", err)
	}

	GormDB = db

	if err := GormDB.AutoMigrate(
		// 基础模型
		&model.ModelAPIModel{},

		&internalmodel.Prompt{},
		&internalmodel.ReverseProgram{},
	); err != nil {
		return fmt.Errorf("gorm auto migrate failed: %w", err)
	}

	if SugaredLogger != nil {
		SugaredLogger.Info("数据库初始化成功")
	}

	return nil
}
