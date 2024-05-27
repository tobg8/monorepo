package secret

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	secretString := String("don't tell")
	assert.Equal(t, secretString.String(), "*****")
	assert.NotEqual(t, string(secretString), "*****")

	secretBytes := Bytes("don't tell")
	assert.Equal(t, secretBytes.String(), "*****")
	assert.NotEqual(t, string(secretBytes), "*****")
}

func TestJSON(t *testing.T) {
	secretString := String("don't tell")
	j, err := json.Marshal(secretString)
	require.NoError(t, err)
	assert.JSONEq(t, `"*****"`, string(j))

	secretBytes := Bytes("don't tell")
	j, err = json.Marshal(secretBytes)
	require.NoError(t, err)
	assert.JSONEq(t, `"*****"`, string(j))
}

func TestJSONWithEmptyValues(t *testing.T) {
	secretString := String("")
	j, err := json.Marshal(secretString)
	require.NoError(t, err)
	assert.JSONEq(t, `""`, string(j))

	secretBytes := Bytes("")
	j, err = json.Marshal(secretBytes)
	require.NoError(t, err)
	assert.JSONEq(t, `""`, string(j))
}

func TestEmbededJSON(t *testing.T) {
	type embed struct {
		Username string `json:"username"`
		Password String `json:"password"`
		Data     Bytes  `json:"data"`
	}

	e := embed{
		Username: "Mini me",
		Password: "azer1234",
		Data:     []byte("My ID card photo"),
	}

	j, err := json.Marshal(e)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"username": "Mini me",
		"password": "*****",
		"data": "*****"
	}`, string(j))
}
