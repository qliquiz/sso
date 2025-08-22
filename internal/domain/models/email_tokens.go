package models

import "time"

// EmailVerification token (UUID in an app, here is as a key)
type EmailVerification struct {
	Token     string
	UserID    UserID
	ExpiresAt time.Time
}

// PasswordReset token
type PasswordReset struct {
	Token     string
	UserID    UserID
	ExpiresAt time.Time
}
