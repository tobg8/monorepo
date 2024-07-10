package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2"
)

// GetRSAKey retrieves an rsa PrivateKey given a secret as string
func GetRSAKey(secret string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(secret))
	if block == nil {
		return nil, errors.New("can't decode pem private key")
	}
	var parsed interface{}
	parsed, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsed, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}
	pkey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}
	return pkey, nil
}

// Deterministic wrt key
func getKid(c Conf) (string, error) {
	ns := c.JWKNamespace
	if ns == "" {
		ns = DefaultJWKNamespace
	}
	namespace, err := uuid.Parse(ns)
	if err != nil {
		return "", err
	}
	return uuid.NewSHA1(namespace, []byte(c.Secret)).String(), nil
}

func getRSAJWK(method, kid string, key *jose.SigningKey) (*jose.JSONWebKey, error) {
	k, ok := key.Key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid key")
	}
	return &jose.JSONWebKey{
		Key:       &k.PublicKey,
		KeyID:     kid,
		Algorithm: method,
		Use:       "sig",
	}, nil
}

// NewRSA returns an RSA JWT implementation
func NewRSA(conf Conf) (*RSA, error) {
	kid, err := getKid(conf)
	if err != nil {
		return nil, err
	}
	options := jose.SignerOptions{}
	options.WithType("JWT")
	options.WithHeader("kid", kid)
	options.EmbedJWK = false
	b, err := newBase(conf, &options)
	if err != nil {
		return nil, err
	}
	jwk, err := GetJWK(conf)
	if err != nil {
		return nil, err
	}
	return &RSA{
		base: *b,
		JWKS: []jose.JSONWebKey{*jwk},
	}, nil
}

// RSA is a structure which implement JWT (un)marshalling
type RSA struct {
	base
	JWKS []jose.JSONWebKey
}
