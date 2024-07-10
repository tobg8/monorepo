package jwt

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/monorepo/common/configloader"
	"github.com/monorepo/common/secret"
)

const (
	// DefaultJWKNamespace is the default JWK namespace.
	DefaultJWKNamespace = "2f219b6a-af92-45b7-be42-c723762c2a35"
)

const (
	// HS256 is a signing method define in JWA RFC
	// ref: https://tools.ietf.org/html/rfc7518#section-3.2
	HS256 = "HS256"
	// HS384 is a signing method define in JWA RFC
	// ref: https://tools.ietf.org/html/rfc7518#section-3.2
	HS384 = "HS384"
	// HS512 is a signing method define in JWA RFC
	// ref: https://tools.ietf.org/html/rfc7518#section-3.2
	HS512 = "HS512"
	// RS256 is a signing method define in JWA RFC
	// ref: https://tools.ietf.org/html/rfc7518#section-3.3
	RS256 = "RS256"
	// RS384 is a signing method define in JWA RFC
	// ref: https://tools.ietf.org/html/rfc7518#section-3.3
	RS384 = "RS384"
	// RS512 is a signing method define in JWA RFC
	// ref: https://tools.ietf.org/html/rfc7518#section-3.3
	RS512 = "RS512"
)

var (
	// ErrInvalidSecret is an error which occurs when the provided secret doesn't respect the rules.
	ErrInvalidSecret = errors.New("provided secret isn't valid")
)

// Conf holds the configuration required to generate JWT tokens
type Conf struct {
	Algorithm string        `mapstructure:"algorithm"`
	Method    string        `mapstructure:"method"`
	Secret    secret.String `mapstructure:"secret"`
	Issuer    string        `mapstructure:"issuer"`
	Audience  []string      `mapstructure:"audience"`

	// This is an arbitrary UUID v4 used as a namespace
	// to get deterministic UUID v5 for JWK jki (key identifier).
	// The keys are determinitic wrt to its content given a
	// namespace.
	JWKNamespace string `mapstructure:"jwk_namespace"`
}

// Envs bind environment keys to env variables
func (*Conf) Envs(l *configloader.Loader) {
	l.BindEnv("algorithm")
	l.BindEnv("method")
}

// Defaults bind environment keys to env variables
func (*Conf) Defaults(l *configloader.Loader) {
	l.SetDefault("algorithm", "HMAC")
	l.SetDefault("method", "HS256")
	l.SetDefault("jwk_namespace", DefaultJWKNamespace)
}

// Unmarshaler define the interface to unmarshal JWT
type Unmarshaler interface {
	UnmarshalJWT(data string, v interface{}) error
}

// UnsafeUnmarshaler define the interface to unmarshal JWT
type UnsafeUnmarshaler interface {
	UnsafeUnmarshalJWT(data []byte, v interface{}) error
}

// Marshaler define the interface to marshal JWT
type Marshaler interface {
	MarshalJWT(i interface{}) (string, error)
}

// JWT define the interface which JWT library must comply to marshal and unmarshal JWT
type JWT interface {
	Marshaler
	Unmarshaler
}

// UnsafeJWT define the interface which JWT library must comply to marshal and unmarshal JWT
type UnsafeJWT interface {
	JWT
	UnsafeUnmarshaler
}

// New return a JWT implementation of token.Generator and token.Parse
// If the signing method is unknown or not defined, HS256 will be used
func New(conf Conf) (JWT, error) {
	switch conf.Algorithm {
	case "HMAC":
		return NewHMAC(conf)
	case "RSA":
		return NewRSA(conf)
	default:
		return NewHMAC(conf)
	}
}

// GetJWK returns the JWK corresponding to the configuration (if applicable)
func GetJWK(c Conf) (*jose.JSONWebKey, error) {
	switch c.Algorithm {
	case "HMAC":
		return nil, nil
	case "RSA":
		kid, err := getKid(c)
		if err != nil {
			return nil, err
		}
		key, _, err := getKeyFromConf(c)
		if err != nil {
			return nil, err
		}
		return getRSAJWK(c.Method, kid, key)
	default:
		return nil, nil
	}
}

// Return signing and validation keys from conf
func getKeyFromConf(c Conf) (*jose.SigningKey, interface{}, error) {
	if c.Secret == "" {
		return nil, nil, ErrInvalidSecret
	}
	secret := string(c.Secret)
	switch c.Method {
	case RS256, RS384, RS512:
		key, err := GetRSAKey(secret)
		if err != nil {
			return nil, nil, err
		}
		k, err := getKey(c.Method, key)
		if err != nil {
			return nil, nil, err
		}
		return k, &key.PublicKey, nil
	case HS256, HS384, HS512:
		k, err := getKey(c.Method, []byte(secret))
		if err != nil {
			return nil, nil, err
		}
		return k, k.Key, nil
	default:
		k, err := getKey(c.Method, []byte(secret))
		if err != nil {
			return nil, nil, err
		}
		return k, k.Key, nil
	}
}

// Parse key from configuration and return lib jose structure
func getKey(method string, k interface{}) (*jose.SigningKey, error) {
	switch method {
	case HS256:
		return &jose.SigningKey{Algorithm: jose.HS256, Key: k}, nil
	case HS384:
		return &jose.SigningKey{Algorithm: jose.HS384, Key: k}, nil
	case HS512:
		return &jose.SigningKey{Algorithm: jose.HS512, Key: k}, nil
	case RS256:
		return &jose.SigningKey{Algorithm: jose.RS256, Key: k}, nil
	case RS384:
		return &jose.SigningKey{Algorithm: jose.RS384, Key: k}, nil
	case RS512:
		return &jose.SigningKey{Algorithm: jose.RS512, Key: k}, nil
	default:
		return &jose.SigningKey{Algorithm: jose.HS256, Key: k}, nil
	}
}

// base is the type containing fields common to multiple marhallers
type base struct {
	Issuer   string
	Audience jwt.Audience
	key      interface{}
	signer   jose.Signer
}

func newBase(c Conf, opts *jose.SignerOptions) (*base, error) {
	sKey, vKey, err := getKeyFromConf(c)
	if err != nil {
		return nil, err
	}
	signer, err := jose.NewSigner(*sKey, opts)
	if err != nil {
		return nil, ErrInvalidSecret
	}
	return &base{
		Issuer:   c.Issuer,
		Audience: jwt.Audience(c.Audience),
		key:      vKey,
		signer:   signer,
	}, nil
}

// Claims
func (b base) Claims() (jwt.Claims, error) {
	u, err := uuid.NewRandom()
	return jwt.Claims{
		Issuer:   b.Issuer,
		ID:       u.String(),
		Audience: b.Audience,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}, err
}

// MarshalJWT implements token.Manager
func (b *base) MarshalJWT(i interface{}) (string, error) {
	builder := jwt.Signed(b.signer)
	claims, err := b.Claims()
	if err != nil {
		return "", err
	}
	builder = builder.Claims(claims)
	builder = builder.Claims(i)
	return builder.CompactSerialize()
}

// UnmarshalJWT implements token.Manager
func (b *base) UnmarshalJWT(data string, v interface{}) error {
	parsed, err := jwt.ParseSigned(data)
	if err != nil {
		return err
	}
	err = parsed.Claims(b.key, v)
	if err != nil {
		return err
	}
	return nil
}

// UnsafeUnmarshalJWT implements token.Manager
func (b *base) UnsafeUnmarshalJWT(data []byte, v interface{}) error {
	parsed, err := jwt.ParseSigned(string(data))
	if err != nil {
		return err
	}
	err = parsed.UnsafeClaimsWithoutVerification(v)
	if err != nil {
		return err
	}
	return nil
}
