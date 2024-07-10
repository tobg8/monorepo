package jwt

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func JWTMarshalUnmarshalTest(t *testing.T, j JWT) {
	t.Run("unmarshalNoError", func(t *testing.T) {
		unmarshalNoError(t, j)
	})
}

func unmarshalNoError(t *testing.T, j JWT) {
	type token struct{ Name string }
	data, err := j.MarshalJWT(&token{"test"})
	require.NoError(t, err)
	var d token
	err = j.UnmarshalJWT(data, &d)
	require.NoError(t, err)
}

// Decode JWT specific base64url encoding with padding stripped
// Borrowed from "github.com/golang-jwt/jwt"
func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
