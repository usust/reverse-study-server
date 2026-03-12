package bootstrap

// Config 是应用启动配置的总入口结构。
//
// 这个结构体的主要职责是承接配置文件中的顶层字段，并通过 `mapstructure`
// 标签与配置加载工具（通常是 viper）完成绑定。
//
// 当前它把配置拆成三个部分：
// 1. ZapLog：日志系统配置；
// 2. DB：数据库驱动选择配置；
// 3. Mysql：MySQL 连接明细配置。
//
// 这样拆分的好处是：
// 1. 顶层结构清晰，便于快速定位配置归属；
// 2. 各配置块可以按模块独立扩展；
// 3. 后续如果增加 Redis、JWT、HTTP Server 等配置，可以继续平铺增加子结构体。
type Config struct {
	// ZapLog 对应配置文件中的 `zap_log` 段，用于初始化日志系统。
	ZapLog ZapLogConfig `mapstructure:"zap_log"`
	// DB 对应配置文件中的 `db` 段，用于描述当前选择的数据库驱动类型。
	DB ConfigDB `mapstructure:"db"`
	// Mysql 对应配置文件中的 `mysql` 段，用于承接 MySQL 连接参数。
	Mysql ConfigMYSQL   `mapstructure:"mysql"`
	Saved ConfigStorage `mapstructure:"storage"`
}
type ConfigStorage struct {
	BaseDir string `mapstructure:"base_dir"`
}

// ConfigDB 表示数据库基础配置。
//
// 这个结构目前只保留 Driver，用来表达“当前使用哪种数据库驱动”，
// 例如 mysql、sqlite、postgres 等。
//
// 这样设计的目的，是把“数据库类型选择”和“某种数据库的详细连接参数”
// 分开管理，避免所有数据库配置项都堆到一个结构里。
type ConfigDB struct {
	// Driver 是数据库驱动名称。
	//
	// 常见取值可以是：
	// 1. mysql
	// 2. sqlite
	// 3. postgres
	//
	// 启动阶段通常会根据这个字段，决定使用哪一套数据库初始化逻辑。
	Driver string `mapstructure:"driver"`
}

// ZapLogConfig 是日志系统（goUtils/logger）的配置结构。
//
// 它用于描述日志文件落盘路径、日志级别，以及日志切割/保留策略。
// 这些字段通常会在服务启动时一次性读取，用于初始化全局日志组件。
type ZapLogConfig struct {
	// LogDir 是日志文件输出目录。
	//
	// 例如：
	// 1. ./logs
	// 2. /var/log/reverse-study-server
	LogDir string `mapstructure:"log_dir"`
	// LogLevel 是日志级别。
	//
	// 常见值有 debug、info、warn、error。
	// 启动时日志组件会根据这个字段决定输出哪些级别的日志。
	LogLevel string `mapstructure:"log_level"`
	// MaxSize 是单个日志文件的最大体积，通常单位为 MB。
	//
	// 超过该大小后，日志库一般会触发滚动切割，生成新文件。
	MaxSize int `mapstructure:"max_size"`
	// MaxBackups 是最多保留的历史日志文件数量。
	//
	// 超过该数量后，最旧的备份日志通常会被删除。
	MaxBackups int `mapstructure:"max_backups"`
	// MaxAge 是日志文件最大保留天数。
	//
	// 超过这个天数的旧日志文件会被清理。
	// 该字段控制的是“时间维度”的保留策略。
	MaxAge int `mapstructure:"max_age"` // max day number of the saved log file
	// Compress 表示是否压缩历史日志文件。
	//
	// 开启后可以减少磁盘占用，但会增加一定的压缩/解压开销。
	Compress bool `mapstructure:"compress"`
}

// ConfigMYSQL 表示配置文件中的 MySQL 连接配置。
//
// 该结构专门承接 MySQL 连接所需的参数，通常会在数据库初始化阶段被组装成
// DSN（Data Source Name）后交给数据库驱动使用。
//
// 注意：
// 1. 这个结构只是“配置承载层”，不直接负责连接数据库；
// 2. 真正的连接逻辑应在 bootstrap / model / repository 等初始化代码中实现。
type ConfigMYSQL struct {
	// Path 是 MySQL 主机地址。
	//
	// 常见值：
	// 1. 127.0.0.1
	// 2. localhost
	// 3. 内网数据库域名
	Path string `mapstructure:"path"`
	// Port 是 MySQL 服务监听端口。
	//
	// 默认通常为 3306。
	Port int `mapstructure:"port"`
	// Username 是数据库用户名。
	//
	// 该账号需要具备当前业务所需的读写权限。
	Username string `mapstructure:"username"`
	// Password 是数据库密码。
	//
	// 这是敏感信息，生产环境应避免打印到日志中。
	Password string `mapstructure:"password"`
	// Database 是应用实际连接的数据库名。
	//
	// 也就是 DSN 中的 schema / database 部分。
	Database string `mapstructure:"database"`
	// Charset 是连接使用的字符集。
	//
	// 常见值例如 utf8mb4。
	// 一般情况下建议保持与数据库实例默认字符集一致。
	Charset string `mapstructure:"charset"`
	// ParseTime 表示驱动是否将时间字段解析为时间类型。
	//
	// 在 Go MySQL 驱动中，这通常是 DSN 里的 `parseTime` 参数。
	// 常见值为 "true"。
	ParseTime string `mapstructure:"parse_time"`
	// Loc 是数据库连接使用的时区配置。
	//
	// 在 Go MySQL 驱动中，这通常对应 DSN 的 `loc` 参数。
	// 常见值例如 Local、Asia%2FShanghai。
	Loc string `mapstructure:"loc"`
}
