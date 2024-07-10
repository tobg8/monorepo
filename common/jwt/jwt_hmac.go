package jwt

// NewHMAC return a JWT implementation of token.Generator and token.Parse
// If the signing method is unknown or not defined, HS256 will be used
func NewHMAC(c Conf) (*HMAC, error) {
	b, err := newBase(c, nil)
	if err != nil {
		return nil, err
	}
	return &HMAC{base: *b}, nil
}

// HMAC is a structure which implement token.Parser and token.Generator in the JWT format.
// JWT is used as a standard to share stateless information between services.
type HMAC struct {
	base
}
