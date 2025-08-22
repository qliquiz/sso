package models

import "time"

// TOTP (один секрет на пользователя)
type MFATOTP struct {
	UserID    UserID
	Secret    SecretBlob // raw secret (not base32 text)
	CreatedAt time.Time
}

// WebAuthn credential
type WebAuthnCredential struct {
	ID              CredID
	UserID          UserID
	PublicKey       []byte
	SignCount       uint32
	AttestationType string
	Transports      []string
	CreatedAt       time.Time
}

// Одноразовые backup-коды (храним хэши)
type MFABackupCode struct {
	UserID UserID
	Hash   TokenHash
	UsedAt *time.Time
}
