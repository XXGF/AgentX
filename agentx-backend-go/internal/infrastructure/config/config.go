package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置（对应 Java 的 application.yml）
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Mail     MailConfig     `mapstructure:"mail"`
	HA       HAConfig       `mapstructure:"high-availability"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        int    `mapstructure:"port"`
	ContextPath string `mapstructure:"context-path"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Name         string `mapstructure:"name"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxIdleConns int    `mapstructure:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns"`
}

// DSN 生成 PostgreSQL 连接字符串
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Shanghai",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int64  `mapstructure:"expiration"` // 过期时间（秒）
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// RabbitMQConfig RabbitMQ 配置
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// MailConfig 邮件配置
type MailConfig struct {
	SMTP SMTPConfig `mapstructure:"smtp"`
}

// SMTPConfig SMTP 配置
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// HAConfig 高可用配置
type HAConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	GatewayURL string `mapstructure:"gateway-url"`
	APIKey     string `mapstructure:"api-key"`
}

// Load 加载配置文件
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// 支持环境变量覆盖（对应 Java 的 ${ENV_VAR:default} 语法）
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// 设置默认值
	v.SetDefault("server.port", 8088)
	v.SetDefault("server.context-path", "/api")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "agentx")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max-idle-conns", 10)
	v.SetDefault("database.max-open-conns", 100)
	v.SetDefault("jwt.expiration", 86400)
	v.SetDefault("logging.level", "debug")

	// 环境变量映射（对应 Java 的 ${SERVER_PORT:8088} 等）
	_ = v.BindEnv("server.port", "SERVER_PORT")
	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.name", "DB_NAME")
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")
	_ = v.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	_ = v.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	_ = v.BindEnv("rabbitmq.username", "RABBITMQ_USERNAME")
	_ = v.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")

	// 读取配置文件（如果存在）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
		// 配置文件不存在时使用默认值和环境变量
		fmt.Fprintln(os.Stderr, "警告: 未找到配置文件，使用默认值和环境变量")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &cfg, nil
}
