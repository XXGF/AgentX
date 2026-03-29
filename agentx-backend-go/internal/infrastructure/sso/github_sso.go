package sso

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	domainSso "github.com/lucky-aeon/agentx/agentx-backend-go/internal/domain/sso"
	"github.com/lucky-aeon/agentx/agentx-backend-go/internal/infrastructure/exception"
	"go.uber.org/zap"
)

// GitHubSsoService GitHub SSO服务实现（对应 Java 的 GitHubSsoService）
type GitHubSsoService struct {
	configProvider *ConfigProvider
	logger         *zap.Logger
}

// NewGitHubSsoService 创建GitHub SSO服务
func NewGitHubSsoService(configProvider *ConfigProvider, logger *zap.Logger) *GitHubSsoService {
	return &GitHubSsoService{
		configProvider: configProvider,
		logger:         logger,
	}
}

func (s *GitHubSsoService) GetProvider() domainSso.SsoProvider {
	return domainSso.SsoProviderGitHub
}

func (s *GitHubSsoService) GetLoginUrl(redirectUrl string) (string, error) {
	config := s.getEffectiveConfig()
	if config == nil {
		return "", exception.NewBusinessException("GitHub SSO配置不完整")
	}

	callbackUrl := redirectUrl
	if callbackUrl == "" {
		callbackUrl = config.RedirectURI
	}

	loginUrl := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=user:email",
		config.AuthorizeURL, config.ClientID, callbackUrl)
	return loginUrl, nil
}

func (s *GitHubSsoService) GetUserInfo(authCode string) (*domainSso.SsoUserInfo, error) {
	config := s.getEffectiveConfig()
	if config == nil {
		return nil, exception.NewBusinessException("GitHub SSO配置不完整")
	}

	// 1. 获取访问令牌
	accessToken, err := s.getAccessToken(authCode, config)
	if err != nil || accessToken == "" {
		return nil, exception.NewBusinessException("获取GitHub访问令牌失败")
	}

	// 2. 获取用户信息
	userInfo, err := s.getGitHubUserInfo(accessToken, config)
	if err != nil || userInfo == nil {
		return nil, exception.NewBusinessException("获取GitHub用户信息失败")
	}

	// 3. 如果用户邮箱为空，尝试获取用户主邮箱
	email, _ := userInfo["email"].(string)
	if email == "" {
		email = s.getPrimaryEmail(accessToken, config)
	}

	// 4. 转换为统一的SsoUserInfo
	name, _ := userInfo["name"].(string)
	login, _ := userInfo["login"].(string)
	avatarUrl, _ := userInfo["avatar_url"].(string)

	var idStr string
	if id, ok := userInfo["id"].(float64); ok {
		idStr = fmt.Sprintf("%.0f", id)
	}

	displayName := name
	if displayName == "" {
		displayName = login
	}

	return &domainSso.SsoUserInfo{
		ID:       idStr,
		Name:     displayName,
		Email:    email,
		Avatar:   avatarUrl,
		Desc:     "GitHub用户: " + login,
		Provider: domainSso.SsoProviderGitHub,
	}, nil
}

// getAccessToken 获取GitHub访问令牌
func (s *GitHubSsoService) getAccessToken(code string, config *GitHubSsoConfig) (string, error) {
	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", config.RedirectURI)

	req, err := http.NewRequest("POST", config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("获取GitHub访问令牌失败", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	token, _ := tokenResp["access_token"].(string)
	return token, nil
}

// getGitHubUserInfo 获取GitHub用户信息
func (s *GitHubSsoService) getGitHubUserInfo(accessToken string, config *GitHubSsoConfig) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", config.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("获取GitHub用户信息失败", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

// getPrimaryEmail 获取GitHub用户主邮箱
func (s *GitHubSsoService) getPrimaryEmail(accessToken string, config *GitHubSsoConfig) string {
	req, err := http.NewRequest("GET", config.UserEmailURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("获取GitHub用户邮箱失败", zap.Error(err))
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var emails []map[string]interface{}
	if err := json.Unmarshal(body, &emails); err != nil {
		return ""
	}

	for _, emailInfo := range emails {
		if primary, ok := emailInfo["primary"].(bool); ok && primary {
			if email, ok := emailInfo["email"].(string); ok {
				return email
			}
		}
	}

	return ""
}

// getEffectiveConfig 获取有效的GitHub配置
func (s *GitHubSsoService) getEffectiveConfig() *GitHubSsoConfig {
	config := s.configProvider.GetGitHubConfig()
	if config.ClientID == "" || config.ClientSecret == "" || config.RedirectURI == "" {
		return nil
	}
	return config
}
