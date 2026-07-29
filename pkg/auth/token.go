// Package auth provides HMAC-SHA256 JWT access token generation and validation for WebRTC rooms.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ClaimGrant defines participant room permissions.
type ClaimGrant struct {
	Room           string `json:"room"`
	CanPublish     bool   `json:"can_publish"`
	CanSubscribe   bool   `json:"can_subscribe"`
	IsAgentAdmin   bool   `json:"is_agent_admin"`
}

// AccessToken holds token configuration.
type AccessToken struct {
	apiKey    string
	apiSecret string
	Identity  string
	Grant     ClaimGrant
	ValidFor  time.Duration
}

// NewAccessToken creates a new access token generator.
func NewAccessToken(apiKey, apiSecret, identity string) *AccessToken {
	return &AccessToken{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		Identity:  identity,
		ValidFor:  24 * time.Hour,
	}
}

// GenerateSignedToken builds a signed JWT token string.
func (t *AccessToken) GenerateSignedToken() (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerBytes, _ := json.Marshal(header)
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)

	now := time.Now().Unix()
	claims := map[string]interface{}{
		"iss":   t.apiKey,
		"sub":   t.Identity,
		"video": t.Grant,
		"iat":   now,
		"exp":   now + int64(t.ValidFor.Seconds()),
	}

	claimsBytes, _ := json.Marshal(claims)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)

	unsignedToken := fmt.Sprintf("%s.%s", encodedHeader, encodedClaims)
	h := hmac.New(sha256.New, []byte(t.apiSecret))
	h.Write([]byte(unsignedToken))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s.%s", unsignedToken, signature), nil
}

// VerifyToken validates a JWT token signature.
func VerifyToken(tokenStr, apiSecret string) bool {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return false
	}

	unsignedToken := fmt.Sprintf("%s.%s", parts[0], parts[1])
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(unsignedToken))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(parts[2]), []byte(expectedSignature))
}
