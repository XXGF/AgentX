package user

import "time"

// UserDTO 用户数据传输对象（对应 Java 的 UserDTO）
type UserDTO struct {
	ID            string    `json:"id"`
	Nickname      string    `json:"nickname"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	GithubID      string    `json:"githubId,omitempty"`
	GithubLogin   string    `json:"githubLogin,omitempty"`
	AvatarURL     string    `json:"avatarUrl,omitempty"`
	LoginPlatform string    `json:"loginPlatform,omitempty"`
	IsAdmin       bool      `json:"isAdmin"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// UserSettingsDTO 用户设置DTO（对应 Java 的 UserSettingsDTO）
type UserSettingsDTO struct {
	ID                    string `json:"id"`
	UserID                string `json:"userId"`
	DefaultModel          string `json:"defaultModel,omitempty"`
	DefaultOcrModel       string `json:"defaultOcrModel,omitempty"`
	DefaultEmbeddingModel string `json:"defaultEmbeddingModel,omitempty"`
}
