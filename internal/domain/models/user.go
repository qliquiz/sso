package models

import "time"

type User struct {
	ID            UserID
	TenantID      TenantID
	Email         Email
	EmailVerified bool
	PassHash      PasswordHash // null/zero if social-only account
	MFAEnabled    bool
	Timestamps
}

func NewUser(id UserID, tenant TenantID, email Email, hash PasswordHash, now time.Time) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	u := &User{
		ID:            id,
		TenantID:      tenant,
		Email:         email,
		EmailVerified: false,
		PassHash:      hash,
		MFAEnabled:    false,
		Timestamps:    Timestamps{CreatedAt: now, UpdatedAt: now},
	}
	return u, nil
}

var (
	ErrInvalidEmail = Err("invalid email")
)

type Err string

func (e Err) Error() string { return string(e) }
