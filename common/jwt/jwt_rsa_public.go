package jwt

import (
	"errors"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// ErrUnknownKey is the error when the JWT is signed with an unknown key.
var ErrUnknownKey = errors.New("JWT signed with an unknown key")

// RSAPublic store the JWKS.
type RSAPublic struct {
	JWKS map[string]jose.JSONWebKey
}

// NewRSAPublic create a new RSAPublic.
func NewRSAPublic(jwks []jose.JSONWebKey) RSAPublic {
	out := RSAPublic{}
	out.JWKS = make(map[string]jose.JSONWebKey)
	for _, jwk := range jwks {
		out.JWKS[jwk.KeyID] = jwk
	}
	return out
}

// UnmarshalJWT unmarshal the JWT.
func (r *RSAPublic) UnmarshalJWT(data string, v interface{}) error {
	parsed, err := jwt.ParseSigned(data)
	if err != nil {
		return err
	}
	var jwk *jose.JSONWebKey
	for _, header := range parsed.Headers {
		j, ok := r.JWKS[header.KeyID]
		if ok {
			jwk = &j
		}
	}
	if jwk == nil {
		return ErrUnknownKey
	}
	err = parsed.Claims(jwk.Key, v)
	if err != nil {
		return err
	}
	return nil
}
