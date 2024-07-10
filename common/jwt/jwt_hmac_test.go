package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewJWT(t *testing.T) {
	j, err := New(Conf{})
	require.Error(t, err)
	require.Nil(t, j)
}

func TestNewJWT_default_signing_method(t *testing.T) {
	j, err := New(Conf{Secret: "test"})
	require.NoError(t, err)
	require.NotNil(t, j)
}

func TestNewJWT_signing_method(t *testing.T) {
	j, err := New(Conf{Secret: "test", Method: HS512})
	require.NoError(t, err)
	require.NotNil(t, j)
}

func TestJWT_Marshal(t *testing.T) {
	j, err := New(Conf{Secret: "test"})
	require.NoError(t, err)
	data, err := j.MarshalJWT(struct{ Name string }{"test"})
	require.NoError(t, err)
	require.NotEmpty(t, data)
	token := data
	tokenParts := strings.Split(token, ".")
	require.Equal(t, 3, len(tokenParts))
	header, body := tokenParts[0], tokenParts[1]

	jsonHeader, err := decodeSegment(header)
	require.NoError(t, err)
	require.Contains(t, string(jsonHeader), "HS256")
	jsonBody, err := decodeSegment(body)
	require.NoError(t, err)
	require.Contains(t, string(jsonBody), "test")
	require.Contains(t, string(jsonBody), "iat")
	require.Contains(t, string(jsonBody), "jti")
}

func TestJWTHMAC_Unmarshal(t *testing.T) {
	j, err := New(Conf{Secret: "test"})
	require.NoError(t, err)
	JWTMarshalUnmarshalTest(t, j)
}

func Test_ExpirationTime_ok(t *testing.T) {
	j, err := New(Conf{Secret: "test"})
	require.NoError(t, err)

	type expiringData struct {
		Data string
		Exp  int64 `json:"exp"`
	}
	d := expiringData{
		Data: "testdata",
		Exp:  time.Now().Add(1 * time.Hour).Unix(),
	}

	data, err := j.MarshalJWT(d)
	require.NoError(t, err)

	var dd expiringData
	err = j.UnmarshalJWT(data, &dd)
	require.NoError(t, err)
}

func Test_ExpirationTime_expired(t *testing.T) {
	j, err := New(Conf{Secret: "test"})
	require.NoError(t, err)

	type expiringData struct {
		Data string
		Exp  int64 `json:"exp"`
	}
	d := expiringData{
		Data: "testdata",
		Exp:  time.Now().Add(-1 * time.Hour).Unix(),
	}

	data, err := j.MarshalJWT(d)
	require.NoError(t, err)

	var dd expiringData
	err = j.UnmarshalJWT(data, dd)
	require.Error(t, err)
}
