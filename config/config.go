package config

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

const (
	logWriters        = "log.writers"
	loggerLevel       = "log.logger_level"
	loggerFile        = "log.logger_file"
	logFormatText     = "log.log_format_text"
	logRlollingPolicy = "log.rollingPolicy"
	logRotateDate     = "log.log_rotate_date"
	logRotateSize     = "log.log_rotate_size"
	logBackupCount    = "log.log_backup_count"
)

type Config struct {
	Name string
}

// Init 初始化配置
func Init(cfg string) error {
	c := Config{Name: cfg}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志包
	c.initLog()

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

// initConfig 初始化配置文件
func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件, 则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件, 则解析默认的配置文件
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")     // 设置配置文件格式为 YAML
	viper.AutomaticEnv()            // 读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER") // 读取环境变量的前缀为 APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	// viper 解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

// initLog 初始化日志包
func (c *Config) initLog() {
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString(logWriters),
		LoggerLevel:    viper.GetString(loggerLevel),
		LoggerFile:     viper.GetString(loggerFile),
		LogFormatText:  viper.GetBool(logFormatText),
		RollingPolicy:  viper.GetString(logRlollingPolicy),
		LogRotateDate:  viper.GetInt(logRotateDate),
		LogRotateSize:  viper.GetInt(logRotateSize),
		LogBackupCount: viper.GetInt(logBackupCount),
	}

	if err := log.InitWithConfig(&passLagerCfg); err != nil {
		return
	}
}

// watchConfig 监控配置文件变化并热加载程序
// 热更新是指: 可以不重启 APISERVER 进程，使 APISERVER 加载最新配置项的值
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
