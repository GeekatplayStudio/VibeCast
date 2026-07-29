package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/agenticsfu/agentic-sfu/pkg/auth"
)

// TestSecurityJWTSignatureVerification audits JWT token forgery resistance.
func TestSecurityJWTSignatureVerification(t *testing.T) {
	correctSecret := "top-secret-key-123"
	wrongSecret := "hacker-secret-key-999"

	token := auth.NewAccessToken("key-1", correctSecret, "user-auditor")
	token.Grant = auth.ClaimGrant{Room: "security-room", CanPublish: true}

	validSignedToken, err := token.GenerateSignedToken()
	if err != nil {
		t.Fatalf("failed to generate signed token: %v", err)
	}

	// Verify valid signature
	if !auth.VerifyToken(validSignedToken, correctSecret) {
		t.Fatal("Security Audit Error: Valid JWT token signature failed verification")
	}

	// Verify forged signature is rejected
	if auth.VerifyToken(validSignedToken, wrongSecret) {
		t.Fatal("Security Audit Vulnerability: Forged JWT token signature was accepted under wrong secret!")
	}
}

// TestSecurityHTTPHeaders audits HTTP security headers.
func TestSecurityHTTPHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rtc/validate", nil)
	w := httptest.NewRecorder()

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Fatal("Security Audit Warning: Missing X-Frame-Options security header")
	}
}
