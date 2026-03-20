package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key-for-unit-tests"

func newTestManager() *Manager {
	return NewManager(testSecret, 30, 7) // 30 min access, 7 day refresh
}

func TestGenerateTokenPair_AccessTokenHasCorrectType(t *testing.T) {
	mgr := newTestManager()

	pair, err := mgr.GenerateTokenPair(1001, "user", 0)
	require.NoError(t, err)
	require.NotEmpty(t, pair.AccessToken)

	claims, err := mgr.ParseToken(pair.AccessToken)
	require.NoError(t, err)

	assert.Equal(t, TokenTypeAccess, claims.TokenType, "access token must have TokenType == 'access'")
	assert.Equal(t, int64(1001), claims.UserID)
	assert.Equal(t, "user", claims.Role)
	assert.Equal(t, 0, claims.MemberLevel)
}

func TestGenerateTokenPair_RefreshTokenHasCorrectType(t *testing.T) {
	mgr := newTestManager()

	pair, err := mgr.GenerateTokenPair(1001, "admin", 2)
	require.NoError(t, err)
	require.NotEmpty(t, pair.RefreshToken)

	claims, err := mgr.ParseToken(pair.RefreshToken)
	require.NoError(t, err)

	assert.Equal(t, TokenTypeRefresh, claims.TokenType, "refresh token must have TokenType == 'refresh'")
	assert.Equal(t, int64(1001), claims.UserID)
	// refresh token should not carry role/member level (security: limit refresh token scope)
}

func TestParseToken_Expired(t *testing.T) {
	// Create a manager with 0 minute access expiry to produce already-expired tokens
	mgr := NewManager(testSecret, 0, 0)

	pair, err := mgr.GenerateTokenPair(1001, "user", 0)
	require.NoError(t, err)

	// Sleep briefly to ensure token is past expiry
	time.Sleep(10 * time.Millisecond)

	_, err = mgr.ParseToken(pair.AccessToken)
	assert.Error(t, err, "expired access token should return parse error")
}

func TestAccessTokenCannotBeUsedAsRefresh(t *testing.T) {
	mgr := newTestManager()

	pair, err := mgr.GenerateTokenPair(1001, "user", 0)
	require.NoError(t, err)

	claims, err := mgr.ParseToken(pair.AccessToken)
	require.NoError(t, err)

	assert.NotEqual(t, TokenTypeRefresh, claims.TokenType,
		"access token must NOT have TokenType == 'refresh'; using access as refresh is a security hole")
}

func TestRefreshTokenCannotBeUsedAsAccess(t *testing.T) {
	mgr := newTestManager()

	pair, err := mgr.GenerateTokenPair(1001, "user", 0)
	require.NoError(t, err)

	claims, err := mgr.ParseToken(pair.RefreshToken)
	require.NoError(t, err)

	assert.NotEqual(t, TokenTypeAccess, claims.TokenType,
		"refresh token must NOT have TokenType == 'access'; using refresh as access is a security hole")
}

func TestTokenPair_ExpiresInMatchesAccessMinutes(t *testing.T) {
	mgr := NewManager(testSecret, 15, 7)

	pair, err := mgr.GenerateTokenPair(1, "user", 0)
	require.NoError(t, err)

	assert.Equal(t, int64(15*60), pair.ExpiresIn, "ExpiresIn should equal accessExpireMin * 60")
}

func TestParseToken_InvalidSigningMethod(t *testing.T) {
	mgr := newTestManager()

	// A token signed with a different secret should fail
	otherMgr := NewManager("different-secret", 30, 7)
	pair, err := otherMgr.GenerateTokenPair(1, "user", 0)
	require.NoError(t, err)

	_, err = mgr.ParseToken(pair.AccessToken)
	assert.Error(t, err, "token signed with different secret should fail verification")
}
