package llm

// ProviderProtocol 服务商协议枚举（对应 Java 的 ProviderProtocol）
type ProviderProtocol string

const (
	ProviderProtocolOpenAI    ProviderProtocol = "OPENAI"
	ProviderProtocolAnthropic ProviderProtocol = "ANTHROPIC"
)

// AllProviderProtocols 所有支持的协议
var AllProviderProtocols = []ProviderProtocol{
	ProviderProtocolOpenAI,
	ProviderProtocolAnthropic,
}

// ProviderProtocolFromCode 根据code获取协议
func ProviderProtocolFromCode(code string) (ProviderProtocol, bool) {
	for _, p := range AllProviderProtocols {
		if string(p) == code {
			return p, true
		}
	}
	return "", false
}

// ModelType 模型类型枚举（对应 Java 的 ModelType）
type ModelType string

const (
	ModelTypeChat      ModelType = "CHAT"      // 对话模型
	ModelTypeEmbedding ModelType = "EMBEDDING" // 嵌入模型
)

// AllModelTypes 所有模型类型
var AllModelTypes = []ModelType{
	ModelTypeChat,
	ModelTypeEmbedding,
}

// ModelTypeFromCode 根据code获取模型类型
func ModelTypeFromCode(code string) (ModelType, bool) {
	for _, t := range AllModelTypes {
		if string(t) == code {
			return t, true
		}
	}
	return "", false
}

// ProviderType 服务商类型枚举（对应 Java 的 ProviderType）
type ProviderType string

const (
	ProviderTypeAll      ProviderType = "all"      // 所有服务商
	ProviderTypeOfficial ProviderType = "official"  // 官方服务商
	ProviderTypeCustom   ProviderType = "custom"    // 用户自定义服务商
)

// ProviderTypeFromCode 根据code获取服务商类型
func ProviderTypeFromCode(code string) ProviderType {
	switch code {
	case "official":
		return ProviderTypeOfficial
	case "custom":
		return ProviderTypeCustom
	default:
		return ProviderTypeAll
	}
}

// ProviderConfig 服务商配置（对应 Java 的 ProviderConfig）
type ProviderConfig struct {
	ApiKey  string `json:"apiKey,omitempty"`
	BaseUrl string `json:"baseUrl,omitempty"`
}

// MaskSensitiveInfo 脱敏处理
func (c *ProviderConfig) MaskSensitiveInfo() {
	if c != nil && c.ApiKey != "" {
		c.ApiKey = "***********"
	}
}
