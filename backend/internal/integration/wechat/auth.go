package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AuthClient struct {
	appID     string
	appSecret string
	client    *http.Client
}

func NewAuthClient(appID, appSecret string) *AuthClient {
	return &AuthClient{
		appID:     appID,
		appSecret: appSecret,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

type code2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (c *AuthClient) Code2Session(code string) (openID, sessionKey string, err error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		c.appID, c.appSecret, code,
	)

	resp, err := c.client.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("request wechat: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("read response: %w", err)
	}

	var result code2SessionResp
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("parse response: %w", err)
	}

	if result.ErrCode != 0 {
		return "", "", fmt.Errorf("wechat error: [%d] %s", result.ErrCode, result.ErrMsg)
	}

	return result.OpenID, result.SessionKey, nil
}
