package jwt

import (
	"fmt"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID      int64  `json:"user_id"`
	Role        string `json:"role"`
	MemberLevel int    `json:"member_level"`
	jwtgo.RegisteredClaims
}

type Manager struct {
	secret          []byte
	accessExpireMin int
	refreshExpireD  int
}

func NewManager(secret string, accessMin, refreshDays int) *Manager {
	return &Manager{
		secret:          []byte(secret),
		accessExpireMin: accessMin,
		refreshExpireD:  refreshDays,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (m *Manager) GenerateTokenPair(userID int64, role string, memberLevel int) (*TokenPair, error) {
	now := time.Now()

	accessClaims := Claims{
		UserID:      userID,
		Role:        role,
		MemberLevel: memberLevel,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(now.Add(time.Duration(m.accessExpireMin) * time.Minute)),
			IssuedAt:  jwtgo.NewNumericDate(now),
			Issuer:    "tline",
		},
	}
	accessToken, err := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, accessClaims).SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := Claims{
		UserID: userID,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(now.Add(time.Duration(m.refreshExpireD) * 24 * time.Hour)),
			IssuedAt:  jwtgo.NewNumericDate(now),
			Issuer:    "tline",
		},
	}
	refreshToken, err := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, refreshClaims).SignedString(m.secret)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(m.accessExpireMin * 60),
	}, nil
}

func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwtgo.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtgo.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
