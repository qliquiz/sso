package models

import "time"

// OAuthIdentity connection to our user
type OAuthIdentity struct {
	Provider       OAuthProvider
	ProviderUserID string
	UserID         UserID
}

// OAuthClient like OIDC provider
type OAuthClient struct {
	ID               ClientID
	TenantID         TenantID
	Name             string
	ClientID         string
	ClientSecretHash PasswordHash
	RedirectURIs     []string // в БД обычно отдельная таблица; в домене можно держать списком
	CreatedAt        time.Time
}

// OAuthAuthCode hash
type OAuthAuthCode struct {
	ID          CodeID // можно хранить как hash-base id
	ClientID    ClientID
	UserID      UserID
	RedirectURI string
	Scope       []string
	CodeHash    TokenHash
	ExpiresAt   time.Time
	CreatedAt   time.Time
	ConsumedAt  *time.Time
}

// OAuthToken refresh for OAuth (если отличен от нашего SSO-refresh — можно переиспользовать RefreshToken)
type OAuthToken struct {
	ID          string // uuid
	UserID      UserID
	ClientID    ClientID
	AccessToken string    // опционально: лучше не хранить, либо хранить hash
	RefreshHash TokenHash // сохранить только хэш
	ExpiresAt   time.Time
	CreatedAt   time.Time
}
