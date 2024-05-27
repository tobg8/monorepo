package configloader

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/monorepo/common/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type Base struct {
	Debug bool `mapstructure:"debug"`
	Value int  `mapstructure:"value"`
}

func (*Base) Envs(l *Loader) {
	l.BindEnv("value")
}

func (*Base) Defaults(l *Loader) {
	l.SetDefault("debug", true)
	l.SetDefault("value", 1)
}

type Conf struct {
	Base `mapstructure:",squash"`

	Key        int                   `mapstructure:"key"`
	Subconf    SubConf               `mapstructure:"subconf"`
	PtrSubconf *SubConf              `mapstructure:"ptrsubconf"`
	Params     map[string]ConfParams `mapstructure:"params"`
}

func (c *Conf) Defaults(l *Loader) {
	l.SetDefault("key", 10)
}

type SubConf struct {
	Key        int         `mapstructure:"key"`
	Subconf    SubSubConf  `mapstructure:"subconf"`
	PtrSubconf *SubSubConf `mapstructure:"ptrsubconf"`
}

func (*SubConf) Defaults(l *Loader) {
	l.SetDefault("key", 200)
}

type SubSubConf struct {
	Key                   int            `mapstructure:"key"`
	SubConfWithoutDefault *SubSubSubConf `mapstructure:"subconf-without-default"`
}

func (SubSubConf) Defaults(l *Loader) {
	l.SetDefault("key", 3000)
}

type SubSubSubConf struct {
	Value int `mapstructure:"value"`
}

func (SubSubSubConf) Envs(l *Loader) {
	l.BindEnv("value")
}

type ConfParams struct {
	A int    `mapstructure:"a"`
	B string `mapstructure:"b"`
}

func (ConfParams) Defaults(l *Loader) {
	l.SetDefault("a", 100)
}

func TestLoader_AddFileReader(t *testing.T) {
	cfgData := strings.NewReader(`---
key: 1
subconf:
  key: 2`)

	loader := New("").
		AddConfigFileReader("test_reader", "yaml", cfgData)

	var conf Conf
	err := loader.Load(&conf)
	require.NoError(t, err)

	assert.Equal(t, 1, conf.Key)
	assert.Equal(t, 2, conf.Subconf.Key)
}

func TestLoader_AddFileReadersInDifferentFormats(t *testing.T) {
	cfgDataYAML := strings.NewReader(`---
key: 1`)
	cfgDataJSON := strings.NewReader(`{
		"subconf": {"key": 2}
	}`)

	loader := New("").
		AddConfigFileReader("test_yaml_reader", "yaml", cfgDataYAML).
		AddConfigFileReader("test_json_reader", "json", cfgDataJSON)

	var conf Conf
	err := loader.Load(&conf)
	require.NoError(t, err)

	assert.Equal(t, 1, conf.Key)
	assert.Equal(t, 2, conf.Subconf.Key)
}

func TestLoader_Load(t *testing.T) {
	loader := New("")
	expectedValue, subconfExpectedValue := 42, 43

	for varname, value := range map[string]int{
		"VALUE": expectedValue,
		"SUBCONF_SUBCONF_SUBCONF_WITHOUT_DEFAULT_VALUE": subconfExpectedValue,
	} {
		_ = os.Setenv(varname, strconv.Itoa(value))
	}

	conf := Conf{
		Params: map[string]ConfParams{
			"foo": {},
			"bar": {},
		},
	}

	require.NoError(t, loader.Load(&conf))
	assert.Equal(t, Conf{
		Base: Base{
			Debug: true,
			Value: expectedValue,
		},
		Key: 10,
		Subconf: SubConf{
			Key: 200,
			Subconf: SubSubConf{
				Key: 3000,
				SubConfWithoutDefault: &SubSubSubConf{
					Value: subconfExpectedValue,
				},
			},
			PtrSubconf: &SubSubConf{
				Key: 3000,
			},
		},
		PtrSubconf: &SubConf{
			Key: 200,
			Subconf: SubSubConf{
				Key: 3000,
			},
			PtrSubconf: &SubSubConf{
				Key: 3000,
			},
		},
		Params: map[string]ConfParams{
			"foo": {
				A: 100,
			},
			"bar": {
				A: 100,
			},
		},
	}, conf)
}

type ConfWithSecrets struct {
	Token secret.String `mapstructure:"token"`
}

func (*ConfWithSecrets) Secrets(l *Loader) {
	l.BindSecret("token", "/path/to/my/secret")
}

func TestLoader_GetSecrets(t *testing.T) {
	secretGetter := vaultClientMock{}
	loader := New("").WithSecretGetter(&secretGetter)

	secretGetter.On("GetSecret", "/path/to/my/secret").Return(
		secret.String("mytoken"),
		nil,
	)

	var conf ConfWithSecrets
	require.NoError(t, loader.Load(&conf))
	assert.Equal(t,
		ConfWithSecrets{
			Token: "mytoken",
		},
		conf,
	)

	secretGetter.AssertExpectations(t)
}

type vaultClientMock struct {
	mock.Mock
}

func (vc *vaultClientMock) GetSecret(path string) (secret.String, error) {
	args := vc.Called(path)
	return args.Get(0).(secret.String), args.Error(1)
}

type ConfWithPost struct {
	DefaultOverride     string `mapstructure:"default-override"`
	SecretOverride      string `mapstructure:"secret-override"`
	FromConfOverride    string `mapstructure:"from-conf-override"`
	DefaultNotOverride  string `mapstructure:"default-not-override"`
	SecretNotOverride   string `mapstructure:"secret-not-override"`
	FromConfNotOverride string `mapstructure:"from-conf-not-override"`
	FromPost            string `mapstructure:"from-post"`
}

func (*ConfWithPost) Defaults(l *Loader) {
	l.SetDefault("default-override", "not-override")
	l.SetDefault("default-not-override", "not-override")
}

func (*ConfWithPost) Secrets(l *Loader) {
	l.BindSecret("secret-override", "/path/to/my/secret")
	l.BindSecret("secret-not-override", "/path/to/my/secret")
}

func (*ConfWithPost) Post(l *Loader) {
	l.SetWithPrefix("default-override", "override")
	l.SetWithPrefix("secret-override", "override")
	l.SetWithPrefix("from-conf-override", "override")
	defaultNotOverride := l.GetWithPrefix("default-not-override")
	if defaultNotOverride == nil {
		l.SetWithPrefix("default-not-override", "override")
	}
	secretNotOverride := l.GetWithPrefix("secret-not-override")
	if secretNotOverride == nil {
		l.SetWithPrefix("secret-not-override", "override")
	}
	fromConfNotOverride := l.GetWithPrefix("from-conf-not-override")
	if fromConfNotOverride == nil {
		l.SetWithPrefix("from-conf-not-override", "override")
	}
	fromPost := l.GetWithPrefix("from-post")
	if fromPost == nil {
		l.SetWithPrefix("from-post", "override")
	}
}

func TestLoader_PostMethod(t *testing.T) {
	cfgData := strings.NewReader(`---
from-conf-override: not-override
from-conf-not-override: not-override`)

	secretGetter := vaultClientMock{}
	loader := New("").
		AddConfigFileReader("test_reader", "yaml", cfgData).
		WithSecretGetter(&secretGetter)

	secretGetter.On("GetSecret", "/path/to/my/secret").Return(
		secret.String("not-override"),
		nil,
	)

	var conf ConfWithPost
	err := loader.Load(&conf)
	require.NoError(t, err)

	assert.Equal(t, "override", conf.DefaultOverride)
	assert.Equal(t, "override", conf.FromConfOverride)
	assert.Equal(t, "override", conf.SecretOverride)
	assert.Equal(t, "not-override", conf.DefaultNotOverride)
	assert.Equal(t, "not-override", conf.FromConfNotOverride)
	assert.Equal(t, "not-override", conf.SecretNotOverride)
	assert.Equal(t, "override", conf.FromPost)
}

func TestLoader_PostMethodExact(t *testing.T) {
	cfgData := strings.NewReader(`---
from-conf-override: not-override
from-conf-not-override: not-override`)

	secretGetter := vaultClientMock{}
	loader := New("").
		AddConfigFileReader("test_reader", "yaml", cfgData).
		WithSecretGetter(&secretGetter)

	secretGetter.On("GetSecret", "/path/to/my/secret").Return(
		secret.String("not-override"),
		nil,
	)

	var conf ConfWithPost
	err := loader.LoadExact(&conf)
	require.NoError(t, err)

	assert.Equal(t, "override", conf.DefaultOverride)
	assert.Equal(t, "override", conf.FromConfOverride)
	assert.Equal(t, "override", conf.SecretOverride)
	assert.Equal(t, "not-override", conf.DefaultNotOverride)
	assert.Equal(t, "not-override", conf.FromConfNotOverride)
	assert.Equal(t, "not-override", conf.SecretNotOverride)
	assert.Equal(t, "override", conf.FromPost)
}
