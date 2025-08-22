package models

import (
	"errors"
	"strings"
	"time"
)

type (
	UserID     string
	TenantID   string
	RoleID     string
	PermID     string
	SessionID  string
	CredID     string // WebAuthn credential id
	ClientID   string // OAuth client id (uuid)
	CodeID     string // OAuth authorization code id (hash key)
	JTI        string // JWT id for access tokens
	KID        string // key id for JWKS
	TokenHash  []byte // SHA-256(refresh) / SHA-256(code) / SHA-256(backup code)
	SecretBlob []byte // opaque key material (TOTP secrets, private keys)
)

// ==== Value Objects ====

type Email string

func NewEmail(s string) (Email, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" || !strings.Contains(s, "@") {
		return "", errors.New("invalid email")
	}
	return Email(s), nil
}

type PasswordHash []byte // Argon2id/Bcrypt hash bytes (not plaintext)

func (h PasswordHash) IsZero() bool { return len(h) == 0 }

// ==== Enums ====

type MFAType string

const (
	MFATypeTOTP     MFAType = "totp"
	MFATypeWebAuthn MFAType = "webauthn"
)

type JWKUse string
type JWKAlg string

const (
	JWKUseSig JWKUse = "sig"

	JWKAlgEdDSA JWKAlg = "EdDSA" // Ed25519
	JWKAlgRS256 JWKAlg = "RS256"
)

type OAuthProvider string

const (
	ProviderGoogle OAuthProvider = "google"
	ProviderGitHub OAuthProvider = "github"
	ProviderApple  OAuthProvider = "apple"
)

type RevocationReason string

const (
	RevokedUserLogout  RevocationReason = "logout"
	RevokedAdmin       RevocationReason = "admin"
	RevokedCompromised RevocationReason = "compromised"
)

// Timestamps common timestamps helper
type Timestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Now() time.Time { return time.Now().UTC() }
