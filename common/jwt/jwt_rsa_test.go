package jwt

import (
	"testing"

	"github.com/monorepo/common/secret"
	"github.com/stretchr/testify/require"
	"gopkg.in/square/go-jose.v2"
)

func generateRSASecret(t *testing.T) secret.String {
	return secret.String(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAp0xNSbYBFh6h1PfiQNgjMxShQq7jAWkPDEVKtGo0XLCGcZry
P+Idzqqsl7K3TCjZpuG6X9XBIA9hlpru6Fv0wadD21i6TVyZWtCCLuGMSAu8AoPj
UvelafkhfqE4gmBBOONRX5lDxOOwnHZz55RtSBMH8vHWFpXRUYDLsLC0GghgoeUB
iKRBCSO5KIskyynC9IcCl6qLr1sH+5kLu4qypRZMMezN2q2eHiyHhpT3oMfAvvq2
BCgIQ4oar/tnOeoeh5n4CWqQ+KYCBmforMFv5roqxr2f3EVckMaXAONtbL1K9WN8
cqCTHBODh3cgDu9t1oCfIb/YyWGWdjEsiBjm7wIDAQABAoIBADon/h4HlO0ZjOw5
l38vI11YaI6DuQn+eWqsk8GPwdAO1U2crWWjtvTmw8SgLbPd53tpsJ4r8kywzB3M
kgxYGwdOm/oeJ/VIoU6+eOLPKTLKUXsWWem3iNsD7a7VYI2B5GpgKyNuZe6FsBlT
3Aq+wBZz9ylvBBspzW/ls+kiJBmADBByQ4LzyHeI8z0wnRsec/3HCBtfvYZgN1Np
IhguvU/qu9E2rZmimoP8SnFKqJBKIKOGvuNOZHJkHPCsf2p4eR/EEsPbsmtoz6IS
ch+/CPB1m9qIMwpCpROTQfXngAZV1l0nEbQAm2aBVHDo1H60K4nz4YhLeImAEhN4
mJFOnwECgYEA1hYT3Ly23o8XPbwYUTq0mVMHgvEVwKMcOxitaDZxe8+Whvv+zxqs
ZFpTUs+EN619na2OOTh8WZVHUlI4jn22HZ82AF8ds1+7hFF3HKJKqv1lztg5iKpt
1uyE+CQWoWSC7aeQsmrMSjxoUTCuKE5q0rbnSHyxSA6AeDLD9IkOlK8CgYEAyA03
yE2KxpsGRUhP69ENdKVloes32oIw7/Mp9uIbyCpK5G1vQXixX86DteCP9B/B9RFi
wGs+SEInu0zSdwfC9ud6Y/H22t2VFB1UtrYEPA0aDu9BzCY0NjdpXA6bsT7ATQtU
XS1ycKmWn1WrZhqhe28SucZrWcU7pdS8Nsij4cECgYEAkxmyZh7JLF13o7SlpNLI
mv2BEMjkoGuzDywiopOeIGt/y5pE+DskrwOdcy4hdDxiLsC9E7YrQ0aeLgNO1yGr
y+jEqzav6rth1kY/qM4eriTVGm5aAfzQ4je8GeB6KEUu7WsQsndNjci6COeBEzLm
lYiVnKoJCjDktzJykIjIGwMCgYEAkFmBS8YwAdjwsGNaT+Vb2TRTXn+0oLXai/mg
6SUEOO3TdnpEkjB9hI0mWsF7/gJAWQ4/fGql2UvrEWqAXyU5mCE1HhMFNa43mPkF
HIXADnjBuc8IYj+a4xgerS9ZRo7qAW3QZR+a+RJVvgj6EUXcCY3/LA+xfGgl/yW+
3aTvI8ECgYEAwruc5SXad5PeJRn6Ciwe6mmsLb/aqupPc1+MByXxz2kDM2vFlDhs
iaWsTXWn2Kt5WDTMMqV6Rq08oiDT//+3MtW+CclVk0eRVHcOaW+xCakEmsExo+mL
e2UIr5j5hmCEznSynewJhCIcL2kK+0H/dRz72evvrfuMc9Vndx/7iM0=
-----END RSA PRIVATE KEY-----
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAp0xNSbYBFh6h1PfiQNgj
MxShQq7jAWkPDEVKtGo0XLCGcZryP+Idzqqsl7K3TCjZpuG6X9XBIA9hlpru6Fv0
wadD21i6TVyZWtCCLuGMSAu8AoPjUvelafkhfqE4gmBBOONRX5lDxOOwnHZz55Rt
SBMH8vHWFpXRUYDLsLC0GghgoeUBiKRBCSO5KIskyynC9IcCl6qLr1sH+5kLu4qy
pRZMMezN2q2eHiyHhpT3oMfAvvq2BCgIQ4oar/tnOeoeh5n4CWqQ+KYCBmforMFv
5roqxr2f3EVckMaXAONtbL1K9WN8cqCTHBODh3cgDu9t1oCfIb/YyWGWdjEsiBjm
7wIDAQAB
-----END PUBLIC KEY-----`)
}

func getDefaultJWTRSA(t *testing.T) *RSA {
	j, err := New(Conf{
		Algorithm: "RSA",
		Method:    "RS256",
		Secret:    generateRSASecret(t),
	})
	require.NoError(t, err)
	return j.(*RSA)
}

func TestNewJWTRSA(t *testing.T) {
	j := getDefaultJWTRSA(t)
	require.NotNil(t, j)
}

func TestJWTHRSA_Unmarshal(t *testing.T) {
	j := getDefaultJWTRSA(t)
	JWTMarshalUnmarshalTest(t, j)
}

func TestJWTHRSA384_Unmarshal(t *testing.T) {
	j, err := New(Conf{
		Algorithm: "RSA",
		Method:    "RS384",
		Secret:    generateRSASecret(t),
	})
	require.NoError(t, err)
	JWTMarshalUnmarshalTest(t, j)
}

func TestJWTHRSA512_Unmarshal(t *testing.T) {
	j, err := New(Conf{
		Algorithm: "RSA",
		Method:    "RS512",
		Secret:    generateRSASecret(t),
	})
	require.NoError(t, err)
	JWTMarshalUnmarshalTest(t, j)
}

func TestJWTRSA_Public_Key_Unmarshal(t *testing.T) {
	c := Conf{
		Algorithm: "RSA",
		Method:    "RS512",
		Secret:    generateRSASecret(t),
	}

	// Unmarshaller
	j, err := NewRSA(c)
	require.NoError(t, err)

	// Produce data
	type Token struct{ Data string }
	tok := Token{Data: "42"}

	// Marshal
	data, err := j.MarshalJWT(&tok)
	require.NoError(t, err)

	// Generate unmarshaller from public key
	require.NotEmpty(t, j.JWKS)
	jwk := j.JWKS[0]
	unmarshaller := NewRSAPublic([]jose.JSONWebKey{jwk})

	// Unmarshal
	var newtok Token
	err = unmarshaller.UnmarshalJWT(data, &newtok)
	require.NoError(t, err)
	require.Equal(t, tok, newtok)
}
