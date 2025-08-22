package models

import (
	"net"
	"time"
)

// Session of user (агрегат для refresh-цепочки)
type Session struct {
	ID        SessionID
	UserID    UserID
	Device    string
	IP        net.IP
	UserAgent string
	CreatedAt time.Time
	RevokedAt *time.Time
	Reason    *RevocationReason
}

// RefreshToken stores only as hash (TokenHash)
type RefreshToken struct {
	Hash              TokenHash // PK
	SessionID         SessionID
	UserID            UserID
	PreviousTokenHash *TokenHash // rotation chain
	CreatedAt         time.Time
	ExpiresAt         time.Time
	ConsumedAt        *time.Time // reuse-detector
	RevokedAt         *time.Time
}

// AccessJTI blacklist for access via JTI (точечный revoke)
type AccessJTI struct {
	JTI       JTI
	RevokedAt time.Time
	Reason    RevocationReason
}
