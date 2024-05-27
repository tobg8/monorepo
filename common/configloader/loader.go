package configloader

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/monorepo/common/secret"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
)

var (
	// ErrFileRead is the error returned when a file related error happens. aezaezae
	ErrFileRead = errors.New("file error")

	// ErrRequiredSecretGetter is the error returned when a secret getter is required to fetch config value.
	ErrRequiredSecretGetter = errors.New("a secret getter is required")
)

// DefaultSetter define an interface which will be used by Loader object to set its default values.
// This interface will be silently called on configuration loading if the configuration object respect the interface.
type DefaultSetter interface {
	Defaults(loader *Loader)
}

// EnvBinder define an interface which will be used by Loader object to call the Envs functions.
// This function aims to define the exported configuration variable to the environment.
type EnvBinder interface {
	Envs(loader *Loader)
}

// SecretBinder define an interface which will be used by Loader object to set secrets values.
// This interface will be silently called on configuration loading if the configuration object respect the interface.
type SecretBinder interface {
	Secrets(loader *Loader)
}

// PostSetter define an interface which will be used by Loader object to set its values after add other configurations (secrets, defaults and file properties).
// This interface will be silently called on configuration loading if the configuration object respect the interface.
type PostSetter interface {
	Post(loader *Loader)
}

// SecretGetter is an interface to retrieve a secret from a key
type secretGetter interface {
	GetSecret(key string) (secret.String, error)
}

// configFileReader groups a reader to a config data
// with its associated format
type configFileReader struct {
	// viper config type such as "yaml", "json", "toml", etc
	format string
	file   io.Reader
}

var (
	defaultSetterType = reflect.TypeOf((*DefaultSetter)(nil)).Elem()
	postSetterType    = reflect.TypeOf((*PostSetter)(nil)).Elem()
	envBinderType     = reflect.TypeOf((*EnvBinder)(nil)).Elem()
	secretBinderType  = reflect.TypeOf((*SecretBinder)(nil)).Elem()
	secretStringKind  = reflect.TypeOf(secret.String("")).Kind()
)

// Loader configures loads a configuration
type Loader struct {
	configFilePaths   []string
	configFileReaders map[string]configFileReader
	v                 *viper.Viper
	prefix            string
	secretClient      secretGetter
	secretErr         error
}

// New creates a new loader
func New(envPrefix string, paths ...string) *Loader {
	v := viper.New()
	replacer := strings.NewReplacer(" ", "_", ".", "_", "-", "_")
	v.SetEnvPrefix(strings.ToUpper(replacer.Replace(envPrefix)))
	v.SetEnvKeyReplacer(replacer)
	l := &Loader{
		configFileReaders: make(map[string]configFileReader),
		v:                 v,
	}
	for _, p := range paths {
		l.AddConfigFile(p)
	}
	return l
}

// WithSecretGetter adds a secretGetter to bind secrets
func (l *Loader) WithSecretGetter(secretClient secretGetter) *Loader {
	l.secretClient = secretClient
	return l
}

// SetDefault sets default value for the given key
func (l *Loader) SetDefault(key string, value interface{}) *Loader {
	l.v.SetDefault(addPrefix(l.prefix, key), value)
	return l
}

// GetWithPrefix value from prefixed key
func (l *Loader) GetWithPrefix(key string) interface{} {
	return l.v.Get(addPrefix(l.prefix, key))
}

// BindEnv binds an env variable
func (l *Loader) BindEnv(input string) *Loader {
	_ = l.v.BindEnv(addPrefix(l.prefix, input))
	return l
}

// BindSecret binds a secret to the given key.
func (l *Loader) BindSecret(key string, secretPath string) *Loader {
	if l.secretClient == nil {
		return l
	}
	if l.secretErr != nil {
		return l
	}

	res, err := l.secretClient.GetSecret(secretPath)
	if err != nil {
		l.secretErr = fmt.Errorf("GetSecret %s: %w", secretPath, err)
		return l
	}

	l.v.Set(key, res)

	return l
}

func addPrefix(prefix, name string) string {
	if len(prefix) == 0 {
		return name
	}
	if len(name) == 0 {
		return prefix
	}
	return prefix + "." + name
}

// Set sets the value for a key
func (l *Loader) Set(key string, value interface{}) *Loader {
	l.v.Set(key, value)
	return l
}

// SetWithPrefix sets the value for a key
func (l *Loader) SetWithPrefix(key string, value interface{}) *Loader {
	l.v.Set(addPrefix(l.prefix, key), value)
	return l
}

// AddConfigFile adds a configuration file path to be read with loading
// the configuration. `path` must be a full path with extension.
func (l *Loader) AddConfigFile(path string) *Loader {
	l.configFilePaths = append(l.configFilePaths, path)
	return l
}

// AddConfigFileReader adds configuration data using io.Reader interface
// This config data will be read when loading the configuration
func (l *Loader) AddConfigFileReader(name, format string, in io.Reader) *Loader {
	l.configFileReaders[name] = configFileReader{
		format: format,
		file:   in,
	}
	return l
}

// Load loads a configuration.
func (l *Loader) Load(configuration interface{}) error {
	v := reflect.ValueOf(configuration)
	l.prepareConfiguration(v, v.Type())
	fileError := l.mergeConfigFiles()
	l.postConfiguration(v, v.Type())

	if err := l.v.Unmarshal(configuration); err != nil {
		return multierr.Combine(fileError, err)
	}

	if fileError != nil {
		return fmt.Errorf("%s: %w", fileError.Error(), ErrFileRead)
	}

	l.fetchVaultPathsFromConf(v, v.Type())

	return l.secretErr
}

// LoadExact loads a configuration, erroring if the target configuration
// param does not contain some field.
func (l *Loader) LoadExact(configuration interface{}) error {
	v := reflect.ValueOf(configuration)
	l.prepareConfiguration(v, v.Type())
	fileError := l.mergeConfigFiles()
	l.postConfiguration(v, v.Type())

	if err := l.v.UnmarshalExact(configuration); err != nil {
		return multierr.Combine(fileError, err)
	}

	if fileError != nil {
		return fmt.Errorf("%s: %w", fileError.Error(), ErrFileRead)
	}
	l.fetchVaultPathsFromConf(v, v.Type())
	return l.secretErr
}
func (l *Loader) mergeConfigFiles() error {
	var fileError error
	for _, p := range l.configFilePaths {
		l.v.SetConfigFile(p)
		if err := l.v.MergeInConfig(); err != nil {
			multierr.AppendInto(&fileError, fmt.Errorf("in file %s: %w", p, err))
		}
	}
	for name, r := range l.configFileReaders {
		l.v.SetConfigType(r.format)
		if err := l.v.MergeConfig(r.file); err != nil {
			multierr.AppendInto(&fileError, fmt.Errorf("in reader %s: %w", name, err))
		}
	}

	return fileError
}

// ExtractConfigFilesRawSettings get configuration files values ignoring defaults/secrets/post values
func (l *Loader) ExtractConfigFilesRawSettings() (map[string]interface{}, error) {
	v := viper.New()
	var fileError error
	for _, p := range l.configFilePaths {
		v.SetConfigFile(p)
		if err := v.MergeInConfig(); err != nil {
			multierr.AppendInto(&fileError, fmt.Errorf("in file %s: %w", p, err))
		}
	}
	for name, r := range l.configFileReaders {
		v.SetConfigType(r.format)
		if err := v.MergeConfig(r.file); err != nil {
			multierr.AppendInto(&fileError, fmt.Errorf("in reader %s: %w", name, err))
		}
	}

	if fileError != nil {
		return nil, fileError
	}

	return v.AllSettings(), nil
}

func (l *Loader) prepareConfiguration(v reflect.Value, t reflect.Type) {
	switch t.Kind() {
	case reflect.Ptr:
		var sv reflect.Value
		if v.IsValid() {
			sv = v.Elem()
		}

		l.prepareConfiguration(sv, t.Elem())

	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			var sv reflect.Value
			if v.IsValid() {
				sv = v.Field(i)
			}

			st := t.Field(i)

			tag, isSquash := cleanTag(st.Tag.Get("mapstructure"))
			if len(tag) == 0 && !isSquash {
				tag = strings.ToLower(st.Name)
			}

			l.withSubPrefix(tag, func() {
				l.prepareConfiguration(sv, st.Type)
			})
		}

		if o, ok := implements(t, defaultSetterType); ok {
			o.(DefaultSetter).Defaults(l)
		}

		if o, ok := implements(t, envBinderType); ok {
			o.(EnvBinder).Envs(l)
		}

		if o, ok := implements(t, secretBinderType); ok {
			o.(SecretBinder).Secrets(l)
		}

	case reflect.Map:
		if v.IsValid() {
			iter := v.MapRange()
			for iter.Next() {
				k := iter.Key()
				if k.Kind() != reflect.String {
					continue
				}

				ks := k.Interface().(string)
				l.withSubPrefix(ks, func() {
					sv := iter.Value()
					l.prepareConfiguration(sv, sv.Type())
				})
			}
		}
	}
}

func (l *Loader) postConfiguration(v reflect.Value, t reflect.Type) {
	switch t.Kind() {
	case reflect.Ptr:
		var sv reflect.Value
		if v.IsValid() {
			sv = v.Elem()
		}

		l.postConfiguration(sv, t.Elem())

	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			var sv reflect.Value
			if v.IsValid() {
				sv = v.Field(i)
			}

			st := t.Field(i)

			tag, isSquash := cleanTag(st.Tag.Get("mapstructure"))
			if len(tag) == 0 && !isSquash {
				tag = strings.ToLower(st.Name)
			}

			l.withSubPrefix(tag, func() {
				l.postConfiguration(sv, st.Type)
			})
		}

		if o, ok := implements(t, postSetterType); ok {
			o.(PostSetter).Post(l)
		}

	case reflect.Map:
		if v.IsValid() {
			iter := v.MapRange()
			for iter.Next() {
				k := iter.Key()
				if k.Kind() != reflect.String {
					continue
				}

				ks := k.Interface().(string)
				l.withSubPrefix(ks, func() {
					sv := iter.Value()
					l.postConfiguration(sv, sv.Type())
				})
			}
		}
	}
}

func (l *Loader) fetchVaultPathsFromConf(v reflect.Value, t reflect.Type) {
	switch t.Kind() {
	case reflect.Ptr:
		var sv reflect.Value
		if v.IsValid() {
			sv = v.Elem()
		}

		l.fetchVaultPathsFromConf(sv, t.Elem())

	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			var sv reflect.Value
			if v.IsValid() {
				sv = v.Field(i)
			}

			st := t.Field(i)
			l.fetchVaultPathsFromConf(sv, st.Type)
		}

	case reflect.Map:
		if v.IsValid() {
			iter := v.MapRange()
			for iter.Next() {
				k := iter.Key()
				if k.Kind() != reflect.String {
					continue
				}

				sv := iter.Value()
				l.fetchVaultPathsFromConf(sv, sv.Type())
			}
		}
	case secretStringKind:
		secretPath := v.String()
		if strings.HasPrefix(secretPath, "VAULT:") {
			if l.secretClient == nil {
				multierr.AppendInto(&l.secretErr, fmt.Errorf("GetSecret %q: %w", secretPath, ErrRequiredSecretGetter))
				return
			}

			secretPath = strings.TrimPrefix(secretPath, "VAULT:")
			res, err := l.secretClient.GetSecret(secretPath)
			if err != nil {
				multierr.AppendInto(&l.secretErr, fmt.Errorf("GetSecret %q: %v", secretPath, err))
			} else {
				v.SetString(string(res))
			}
		}
	}
}

func implements(t reflect.Type, u reflect.Type) (interface{}, bool) {
	switch {
	case t.Implements(u):
		return reflect.Zero(t).Interface(), true
	case reflect.PtrTo(t).Implements(u):
		return reflect.New(t).Interface(), true
	default:
		return nil, false
	}
}

func (l *Loader) withSubPrefix(subprefix string, cb func()) {
	if subprefix != "" {
		prev := l.prefix

		l.prefix = addPrefix(l.prefix, subprefix)

		defer func() {
			l.prefix = prev
		}()
	}

	cb()
}

func cleanTag(tag string) (name string, isSquash bool) {
	t := strings.Split(tag, ",")
	if len(t) > 0 {
		name = t[0]
	}
	for _, s := range t {
		if s == "squash" {
			isSquash = true
			break
		}
	}
	return
}
