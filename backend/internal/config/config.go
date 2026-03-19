package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	WeChat   WeChatConfig   `mapstructure:"wechat"`
	SMS      SMSConfig      `mapstructure:"sms"`
	OSS      OSSConfig      `mapstructure:"oss"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // debug, release, test
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`    // MB
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`     // days
	Compress   bool   `mapstructure:"compress"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	AccessExpireMin int    `mapstructure:"access_expire_min"`
	RefreshExpireD  int    `mapstructure:"refresh_expire_days"`
}

type WeChatConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	MchID     string `mapstructure:"mch_id"`
	MchAPIKey string `mapstructure:"mch_api_key"`
	NotifyURL string `mapstructure:"notify_url"`
}

type SMSConfig struct {
	Provider  string `mapstructure:"provider"` // aliyun, tencent
	AccessKey string `mapstructure:"access_key"`
	Secret    string `mapstructure:"secret"`
	SignName  string `mapstructure:"sign_name"`
	Template  string `mapstructure:"template"`
}

type OSSConfig struct {
	Provider  string `mapstructure:"provider"` // aliyun, tencent
	Endpoint  string `mapstructure:"endpoint"`
	Bucket    string `mapstructure:"bucket"`
	AccessKey string `mapstructure:"access_key"`
	Secret    string `mapstructure:"secret"`
	BaseURL   string `mapstructure:"base_url"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	configDir := os.Getenv("CONFIG_DIR")
	if configDir == "" {
		execPath, _ := os.Executable()
		configDir = filepath.Join(filepath.Dir(execPath), "..", "..", "configs")
	}
	v.AddConfigPath(configDir)
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	v.AutomaticEnv()

	// defaults
	v.SetDefault("app.name", "tline-api")
	v.SetDefault("app.port", 8080)
	v.SetDefault("app.mode", "debug")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_idle_conns", 10)
	v.SetDefault("database.max_open_conns", 100)
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("log.level", "info")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 5)
	v.SetDefault("log.max_age", 30)
	v.SetDefault("jwt.access_expire_min", 120)
	v.SetDefault("jwt.refresh_expire_days", 30)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
