package models

import "time"

// JWK for JWKS
type JWK struct {
	KID        KID
	Alg        JWKAlg
	Use        JWKUse
	PublicKey  []byte
	PrivateKey SecretBlob
	CreatedAt  time.Time
	NotBefore  time.Time
	NotAfter   *time.Time
	Active     bool
}
