package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 应用配置（对应 Java 的 application.yml）
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	RabbitMQ  RabbitMQConfig  `mapstructure:"rabbitmq"`
	Mail      MailConfig      `mapstructure:"mail"`
	HA        HAConfig        `mapstructure:"high-availability"`
	GitHub    GitHubConfig    `mapstructure:"github"`
	Docker    DockerConfig    `mapstructure:"docker"`
	Payment   PaymentConfig   `mapstructure:"payment"`
	Upload    UploadConfig    `mapstructure:"upload"`
	Embedding EmbeddingConfig `mapstructure:"embedding"`
	Memory    MemoryConfig    `mapstructure:"memory"`
	Rerank    RerankConfig    `mapstructure:"rerank"`
	RAG       RAGConfig       `mapstructure:"rag"`
	MCP       MCPConfig       `mapstructure:"mcp"`
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
	SMTP         SMTPConfig         `mapstructure:"smtp"`
	Verification VerificationConfig `mapstructure:"verification"`
}

// SMTPConfig SMTP 配置
type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// VerificationConfig 验证码邮件配置
type VerificationConfig struct {
	Template string `mapstructure:"template"`
	Subject  string `mapstructure:"subject"`
}

// HAConfig 高可用配置
type HAConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	GatewayURL     string `mapstructure:"gateway-url"`
	APIKey         string `mapstructure:"api-key"`
	ConnectTimeout int    `mapstructure:"connect-timeout"`
	ReadTimeout    int    `mapstructure:"read-timeout"`
}

// GitHubConfig GitHub 配置（对应 Java 的 github.target）
type GitHubConfig struct {
	Target GitHubTargetConfig `mapstructure:"target"`
}

// GitHubTargetConfig GitHub 目标仓库配置
type GitHubTargetConfig struct {
	RepoName string `mapstructure:"repo-name"`
	Username string `mapstructure:"username"`
	Token    string `mapstructure:"token"`
}

// DockerConfig Docker/容器管理配置（对应 Java 的 agentx.container）
type DockerConfig struct {
	Host                   string `mapstructure:"host"`
	UserVolumeBasePath     string `mapstructure:"user-volume-base-path"`
	DefaultMcpGatewayImage string `mapstructure:"default-mcp-gateway-image"`
	MonitorInterval        int    `mapstructure:"monitor-interval"`
	StatsUpdateInterval    int    `mapstructure:"stats-update-interval"`
}

// PaymentConfig 支付配置（对应 Java 的 payment）
type PaymentConfig struct {
	Alipay AlipayConfig `mapstructure:"alipay"`
	Stripe StripeConfig `mapstructure:"stripe"`
}

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID      string `mapstructure:"app-id"`
	PrivateKey string `mapstructure:"private-key"`
	PublicKey  string `mapstructure:"public-key"`
	NotifyURL  string `mapstructure:"notify-url"`
	ReturnURL  string `mapstructure:"return-url"`
}

// StripeConfig Stripe 配置
type StripeConfig struct {
	SecretKey      string `mapstructure:"secret-key"`
	PublishableKey string `mapstructure:"publishable-key"`
	WebhookSecret  string `mapstructure:"webhook-secret"`
}

// UploadConfig 文件上传配置（对应 Java 的 dromara.x-file-storage）
type UploadConfig struct {
	Platform string   `mapstructure:"platform"` // local / s3
	Local    LocalUploadConfig `mapstructure:"local"`
	S3       S3Config `mapstructure:"s3"`
}

// LocalUploadConfig 本地上传配置
type LocalUploadConfig struct {
	BasePath string `mapstructure:"base-path"`
}

// S3Config S3/对象存储配置
type S3Config struct {
	AccessKey  string `mapstructure:"access-key"`
	SecretKey  string `mapstructure:"secret-key"`
	Region     string `mapstructure:"region"`
	Endpoint   string `mapstructure:"endpoint"`
	BucketName string `mapstructure:"bucket-name"`
	Domain     string `mapstructure:"domain"`
	BasePath   string `mapstructure:"base-path"`
}

// EmbeddingConfig Embedding 向量化配置（对应 Java 的 embedding）
type EmbeddingConfig struct {
	Name        string           `mapstructure:"name"`
	VectorStore VectorStoreConfig `mapstructure:"vector-store"`
}

// VectorStoreConfig 向量存储配置
type VectorStoreConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Database       string `mapstructure:"database"`
	Table          string `mapstructure:"table"`
	Dimension      int    `mapstructure:"dimension"`
	DropTableFirst bool   `mapstructure:"drop-table-first"`
	CreateTable    bool   `mapstructure:"create-table"`
}

// MemoryConfig Memory 向量存储配置（对应 Java 的 memory.embedding）
type MemoryConfig struct {
	Embedding MemoryEmbeddingConfig `mapstructure:"embedding"`
}

// MemoryEmbeddingConfig Memory Embedding 配置
type MemoryEmbeddingConfig struct {
	VectorStore VectorStoreConfig `mapstructure:"vector-store"`
}

// RerankConfig Rerank 重排序配置（对应 Java 的 rerank）
type RerankConfig struct {
	Name    string `mapstructure:"name"`
	APIKey  string `mapstructure:"api-key"`
	APIURL  string `mapstructure:"api-url"`
	Model   string `mapstructure:"model"`
	Timeout int    `mapstructure:"timeout"`
}

// RAGConfig RAG 配置（对应 Java 的 rag）
type RAGConfig struct {
	Markdown RAGMarkdownConfig `mapstructure:"markdown"`
	Vector   RAGVectorConfig   `mapstructure:"vector"`
}

// RAGMarkdownConfig RAG Markdown 分段配置
type RAGMarkdownConfig struct {
	SegmentSplit RAGSegmentSplitConfig `mapstructure:"segment-split"`
}

// RAGSegmentSplitConfig RAG 分段切分配置
type RAGSegmentSplitConfig struct {
	Enabled       bool `mapstructure:"enabled"`
	MaxLength     int  `mapstructure:"max-length"`
	MinLength     int  `mapstructure:"min-length"`
	BufferSize    int  `mapstructure:"buffer-size"`
	EnableOverlap bool `mapstructure:"enable-overlap"`
	OverlapSize   int  `mapstructure:"overlap-size"`
}

// RAGVectorConfig RAG 向量配置
type RAGVectorConfig struct {
	MaxLength   int `mapstructure:"max-length"`
	MinLength   int `mapstructure:"min-length"`
	OverlapSize int `mapstructure:"overlap-size"`
}

// MCPConfig MCP 网关配置（对应 Java 的 mcp.gateway）
type MCPConfig struct {
	Gateway MCPGatewayConfig `mapstructure:"gateway"`
}

// MCPGatewayConfig MCP 网关配置
type MCPGatewayConfig struct {
	ConnectTimeout int `mapstructure:"connect-timeout"`
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

	// ========== 设置默认值 ==========
	// Server
	v.SetDefault("server.port", 8088)
	v.SetDefault("server.context-path", "/api")

	// Database
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "agentx")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max-idle-conns", 10)
	v.SetDefault("database.max-open-conns", 100)

	// JWT
	v.SetDefault("jwt.expiration", 86400)

	// Logging
	v.SetDefault("logging.level", "debug")

	// Mail
	v.SetDefault("mail.smtp.host", "smtp.qq.com")
	v.SetDefault("mail.smtp.port", 587)
	v.SetDefault("mail.verification.template", "您的验证码是:%s，有效期10分钟，请勿泄露给他人。")
	v.SetDefault("mail.verification.subject", "AgentX - 邮箱验证码")

	// GitHub
	v.SetDefault("github.target.repo-name", "agent-mcp-community")
	v.SetDefault("github.target.username", "lucky-aeon")

	// Docker
	v.SetDefault("docker.host", "unix:///var/run/docker.sock")
	v.SetDefault("docker.user-volume-base-path", "/docker/users")
	v.SetDefault("docker.default-mcp-gateway-image", "ghcr.io/lucky-aeon/mcp-gateway:latest")
	v.SetDefault("docker.monitor-interval", 300000)
	v.SetDefault("docker.stats-update-interval", 120000)

	// Upload
	v.SetDefault("upload.platform", "local")
	v.SetDefault("upload.local.base-path", "./uploads")
	v.SetDefault("upload.s3.base-path", "s3/")

	// Embedding
	v.SetDefault("embedding.name", "OpenAI")
	v.SetDefault("embedding.vector-store.table", "public.vector_store")
	v.SetDefault("embedding.vector-store.dimension", 1024)
	v.SetDefault("embedding.vector-store.create-table", true)

	// Memory
	v.SetDefault("memory.embedding.vector-store.table", "public.memory_vector_store")
	v.SetDefault("memory.embedding.vector-store.dimension", 1024)
	v.SetDefault("memory.embedding.vector-store.create-table", true)

	// Rerank
	v.SetDefault("rerank.name", "OpenAI")
	v.SetDefault("rerank.model", "Pro/BAAI/bge-reranker-v2-m3")
	v.SetDefault("rerank.timeout", 30000)

	// RAG
	v.SetDefault("rag.markdown.segment-split.enabled", true)
	v.SetDefault("rag.markdown.segment-split.max-length", 1800)
	v.SetDefault("rag.markdown.segment-split.min-length", 200)
	v.SetDefault("rag.markdown.segment-split.buffer-size", 100)
	v.SetDefault("rag.markdown.segment-split.enable-overlap", false)
	v.SetDefault("rag.markdown.segment-split.overlap-size", 50)
	v.SetDefault("rag.vector.max-length", 1800)
	v.SetDefault("rag.vector.min-length", 200)
	v.SetDefault("rag.vector.overlap-size", 100)

	// MCP
	v.SetDefault("mcp.gateway.connect-timeout", 60000)

	// High Availability
	v.SetDefault("high-availability.connect-timeout", 30000)
	v.SetDefault("high-availability.read-timeout", 60000)

	// ========== 环境变量映射 ==========
	// Server
	_ = v.BindEnv("server.port", "SERVER_PORT")

	// Database
	_ = v.BindEnv("database.host", "DB_HOST")
	_ = v.BindEnv("database.port", "DB_PORT")
	_ = v.BindEnv("database.user", "DB_USER")
	_ = v.BindEnv("database.password", "DB_PASSWORD")
	_ = v.BindEnv("database.name", "DB_NAME")

	// JWT
	_ = v.BindEnv("jwt.secret", "JWT_SECRET")

	// RabbitMQ
	_ = v.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	_ = v.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	_ = v.BindEnv("rabbitmq.username", "RABBITMQ_USERNAME")
	_ = v.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")

	// Mail
	_ = v.BindEnv("mail.smtp.host", "MAIL_SMTP_HOST")
	_ = v.BindEnv("mail.smtp.port", "MAIL_SMTP_PORT")
	_ = v.BindEnv("mail.smtp.username", "MAIL_SMTP_USERNAME")
	_ = v.BindEnv("mail.smtp.password", "MAIL_SMTP_PASSWORD")
	_ = v.BindEnv("mail.verification.template", "MAIL_VERIFICATION_TEMPLATE")
	_ = v.BindEnv("mail.verification.subject", "MAIL_VERIFICATION_SUBJECT")

	// GitHub
	_ = v.BindEnv("github.target.repo-name", "GITHUB_REPO_NAME")
	_ = v.BindEnv("github.target.username", "GITHUB_USERNAME")
	_ = v.BindEnv("github.target.token", "GITHUB_TOKEN")

	// Docker
	_ = v.BindEnv("docker.host", "AGENTX_CONTAINER_DOCKER_HOST")
	_ = v.BindEnv("docker.user-volume-base-path", "AGENTX_CONTAINER_USER_VOLUME_PATH")
	_ = v.BindEnv("docker.default-mcp-gateway-image", "AGENTX_CONTAINER_DEFAULT_MCP_IMAGE")
	_ = v.BindEnv("docker.monitor-interval", "AGENTX_CONTAINER_MONITOR_INTERVAL")
	_ = v.BindEnv("docker.stats-update-interval", "AGENTX_CONTAINER_STATS_INTERVAL")

	// Payment - Alipay
	_ = v.BindEnv("payment.alipay.app-id", "ALIPAY_APP_ID")
	_ = v.BindEnv("payment.alipay.private-key", "ALIPAY_PRIVATE_KEY")
	_ = v.BindEnv("payment.alipay.public-key", "ALIPAY_PUBLIC_KEY")
	_ = v.BindEnv("payment.alipay.notify-url", "ALIPAY_NOTIFY_URL")
	_ = v.BindEnv("payment.alipay.return-url", "ALIPAY_RETURN_URL")

	// Payment - Stripe
	_ = v.BindEnv("payment.stripe.secret-key", "STRIPE_SECRET_KEY")
	_ = v.BindEnv("payment.stripe.publishable-key", "STRIPE_PUBLISHABLE_KEY")
	_ = v.BindEnv("payment.stripe.webhook-secret", "STRIPE_WEBHOOK_SECRET")

	// Upload - S3
	_ = v.BindEnv("upload.platform", "UPLOAD_PLATFORM")
	_ = v.BindEnv("upload.s3.access-key", "S3_SECRET_ID")
	_ = v.BindEnv("upload.s3.secret-key", "S3_SECRET_KEY")
	_ = v.BindEnv("upload.s3.region", "S3_REGION")
	_ = v.BindEnv("upload.s3.endpoint", "S3_ENDPOINT")
	_ = v.BindEnv("upload.s3.bucket-name", "S3_BUCKET_NAME")
	_ = v.BindEnv("upload.s3.domain", "S3_DOMAIN")

	// Embedding / Vector Store
	_ = v.BindEnv("embedding.vector-store.host", "VECTOR_DB_HOST")
	_ = v.BindEnv("embedding.vector-store.port", "VECTOR_DB_PORT")
	_ = v.BindEnv("embedding.vector-store.user", "VECTOR_DB_USER")
	_ = v.BindEnv("embedding.vector-store.password", "VECTOR_DB_PASSWORD")
	_ = v.BindEnv("embedding.vector-store.database", "VECTOR_DB_NAME")
	_ = v.BindEnv("embedding.vector-store.table", "VECTOR_DB_TABLE")
	_ = v.BindEnv("embedding.vector-store.dimension", "VECTOR_DB_DIMENSION")

	// Rerank
	_ = v.BindEnv("rerank.api-key", "SILICONFLOW_API_KEY")
	_ = v.BindEnv("rerank.api-url", "SILICONFLOW_API_URL_RERANK")
	_ = v.BindEnv("rerank.model", "SILICONFLOW_MODEL_RERANK_MODEL")

	// High Availability
	_ = v.BindEnv("high-availability.enabled", "HIGH_AVAILABILITY_ENABLED")
	_ = v.BindEnv("high-availability.gateway-url", "HIGH_AVAILABILITY_GATEWAY_URL")
	_ = v.BindEnv("high-availability.api-key", "HIGH_AVAILABILITY_API_KEY")

	// MCP
	_ = v.BindEnv("mcp.gateway.connect-timeout", "MCP_GATEWAY_CONNECT_TIMEOUT")

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

	// 向量存储配置回退到数据库配置
	if cfg.Embedding.VectorStore.Host == "" {
		cfg.Embedding.VectorStore.Host = cfg.Database.Host
	}
	if cfg.Embedding.VectorStore.Port == 0 {
		cfg.Embedding.VectorStore.Port = cfg.Database.Port
	}
	if cfg.Embedding.VectorStore.User == "" {
		cfg.Embedding.VectorStore.User = cfg.Database.User
	}
	if cfg.Embedding.VectorStore.Password == "" {
		cfg.Embedding.VectorStore.Password = cfg.Database.Password
	}
	if cfg.Embedding.VectorStore.Database == "" {
		cfg.Embedding.VectorStore.Database = cfg.Database.Name
	}

	// Memory 向量存储回退到 Embedding 向量存储配置
	if cfg.Memory.Embedding.VectorStore.Host == "" {
		cfg.Memory.Embedding.VectorStore.Host = cfg.Embedding.VectorStore.Host
	}
	if cfg.Memory.Embedding.VectorStore.Port == 0 {
		cfg.Memory.Embedding.VectorStore.Port = cfg.Embedding.VectorStore.Port
	}
	if cfg.Memory.Embedding.VectorStore.User == "" {
		cfg.Memory.Embedding.VectorStore.User = cfg.Embedding.VectorStore.User
	}
	if cfg.Memory.Embedding.VectorStore.Password == "" {
		cfg.Memory.Embedding.VectorStore.Password = cfg.Embedding.VectorStore.Password
	}
	if cfg.Memory.Embedding.VectorStore.Database == "" {
		cfg.Memory.Embedding.VectorStore.Database = cfg.Embedding.VectorStore.Database
	}

	return &cfg, nil
}
