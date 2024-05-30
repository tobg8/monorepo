package awsx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awstrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/aws/aws-sdk-go-v2/aws"

	"github.com/monorepo/common/logging"
)

// Doer is a simple standard interface to manage HTTP requests
// For usage with httputils, you can declare your preferred httpClient compatible with httputils.Doer
// and use the httputils.Client.Client exposed by the implementation.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// AWSRole is a type to allow easy configuration
type AWSRole string

// AWSRegionRegexp the regexp to validate aws regions
var AWSRegionRegexp = regexp.MustCompile(`(us(-gov)?|ap|ca|cn|eu|sa)-(central|(north|south)?(east|west)?)-\d`)

// Config represents a configuration for sqsx services
type Config struct {
	Region        string `mapstructure:"region"`
	TokenFilePath string `mapstructure:"token_file_path"`
	Profile       string `mapstructure:"profile"`
}

// CredentialsCacheConfig holds the credentials cache options
// see: https://github.com/aws/aws-sdk-go-v2/blob/main/aws/credential_cache.go#L14
type CredentialsCacheConfig struct {
	ExpiryWindow           time.Duration `mapstructure:"expiry-window"`
	ExpiryWindowJitterFrac float64       `mapstructure:"expiry-window-jitter-frac"`
}

// AWSConfigBuilder is a Builder for aws.Config.
type AWSConfigBuilder struct {
	credsProviderBuilder func(ctx context.Context, opts ...func(*awsconfig.LoadOptions) error) (aws.CredentialsProvider, error)
	config               *Config
	logger               logging.Logger
	httpClient           Doer
	endpoint             *string
	withTracing          bool
	serviceName          string
	err                  error
}

// NewAWSConfigBuilder will return a AWSConfigBuilder carrying the aws.Config
// this struct allows to customise some options (assumerole, logging, ...) before building the aws.Config
func NewAWSConfigBuilder(config *Config) *AWSConfigBuilder {
	return &AWSConfigBuilder{
		config: config,
		err:    CheckMandatoryGlobalConfig(config),
	}
}

// Build builds the aws.Config to be used in filestore or sqsx
func (b *AWSConfigBuilder) Build(ctx context.Context) (aws.Config, error) {
	if b.err != nil {
		return aws.Config{}, b.err
	}

	opts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(b.config.Region),
		awsconfig.WithSharedConfigProfile(b.config.Profile),
	}

	if b.logger != nil {
		opts = append(opts, awsconfig.WithLogger(AwsLoggerfromLogger(b.logger)))
	}

	if b.httpClient != nil {
		opts = append(opts, awsconfig.WithHTTPClient(b.httpClient))
	}

	if b.endpoint != nil {
		opts = append(opts, awsconfig.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service string, region string, opts ...any) (aws.Endpoint, error) {
					return aws.Endpoint{
						PartitionID:       "aws",
						URL:               *b.endpoint,
						SigningRegion:     b.config.Region,
						HostnameImmutable: true,
					}, nil
				},
			),
		))
	}

	if b.credsProviderBuilder != nil {
		credsProvider, err := b.credsProviderBuilder(ctx, opts...)
		if err != nil {
			return aws.Config{}, err
		}

		opts = append(opts, awsconfig.WithCredentialsProvider(credsProvider))
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("can't create config: %v", err)
	}

	if b.withTracing {
		awstrace.AppendMiddleware(&cfg,
			awstrace.WithServiceName(b.serviceName),
		)
	}

	return cfg, nil
}

// WithAssumeRole will override the config credential provider to switch to a stscreds.AssumeRoleProvider
// will check if config has been Built before setting credentials.
// This option is mutually exclusive with `WithStaticCredentials()`.
func (b *AWSConfigBuilder) WithAssumeRole(role AWSRole) *AWSConfigBuilder {
	if role == "" {
		return b
	}

	return b.addCredsProviderBuilder(func(ctx context.Context, opts ...func(*awsconfig.LoadOptions) error) (aws.CredentialsProvider, error) {
		assumeCfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("can't create assume config: %v", err)
		}

		if b.withTracing {
			awstrace.AppendMiddleware(&assumeCfg,
				awstrace.WithServiceName(b.serviceName),
			)
		}

		stsClient := sts.NewFromConfig(assumeCfg)

		return stscreds.NewAssumeRoleProvider(stsClient, string(role)), nil
	})
}

// WithStaticCredentials configures the static credentials to use.
// This option is mutually exclusive with `WithAssumeRole()`.
func (b *AWSConfigBuilder) WithStaticCredentials(accessKeyID string, secretAccessKey string, sessionToken string) *AWSConfigBuilder {
	if accessKeyID == "" && secretAccessKey == "" && sessionToken == "" {
		return b
	}

	return b.addCredsProviderBuilder(func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.CredentialsProvider, error) {
		return credentials.NewStaticCredentialsProvider(
			accessKeyID, secretAccessKey, sessionToken,
		), nil
	})
}

func (b *AWSConfigBuilder) addCredsProviderBuilder(
	credsProviderBuilder func(_ context.Context, _ ...func(*awsconfig.LoadOptions) error) (aws.CredentialsProvider, error),
) *AWSConfigBuilder {
	if b.credsProviderBuilder != nil {
		return b.withErrorf("credentials provider is already registered")
	}

	b.credsProviderBuilder = credsProviderBuilder

	return b
}

// WithCredentialsCache allows to specify a cache strategy for the credentials.
// A credentials provider, defined by `WithAssumeRole()` or `WithStaticCredentials()`,
// must be defined before to call this option.
func (b *AWSConfigBuilder) WithCredentialsCache(conf CredentialsCacheConfig) *AWSConfigBuilder {
	if b.credsProviderBuilder == nil {
		return b.withErrorf("can't define a credentials cache without a credentials provider")
	}

	oldCredsProviderBuilder := b.credsProviderBuilder

	b.credsProviderBuilder = func(ctx context.Context, opts ...func(*awsconfig.LoadOptions) error) (aws.CredentialsProvider, error) {
		credsProvider, err := oldCredsProviderBuilder(ctx, opts...)
		if err != nil {
			return credsProvider, err
		}

		return aws.NewCredentialsCache(credsProvider, func(options *aws.CredentialsCacheOptions) {
			options.ExpiryWindow = conf.ExpiryWindow
			options.ExpiryWindowJitterFrac = conf.ExpiryWindowJitterFrac
		}), nil
	}

	return b
}

// WithLogger allows to specify a specific logger. Compatible with logging.logger/ glog / ...
func (b *AWSConfigBuilder) WithLogger(logger logging.Logger) *AWSConfigBuilder {
	b.logger = logger
	return b
}

// WithHTTPClient allows to specify a httpclient. Compatible with httputils.Client.Client.
func (b *AWSConfigBuilder) WithHTTPClient(httpClient Doer) *AWSConfigBuilder {
	b.httpClient = httpClient
	return b
}

// WithEndpoint allows to specify an endpoint, useful for testing purpose.
func (b *AWSConfigBuilder) WithEndpoint(endpoint *string) *AWSConfigBuilder {
	b.endpoint = endpoint
	return b
}

// WithTracing activates the tracing and sets the given service name for the dialled connection.
func (b *AWSConfigBuilder) WithTracing(serviceName string) *AWSConfigBuilder {
	b.withTracing = true
	b.serviceName = serviceName
	return b
}

func (b *AWSConfigBuilder) withErrorf(format string, a ...any) *AWSConfigBuilder {
	b.err = errors.Join(b.err, fmt.Errorf(format, a...))
	return b
}

// DefaultRegion represents the default aws region
const DefaultRegion string = "eu-west-3"

// CheckMandatoryGlobalConfig checks the basis of the aws configuration
func CheckMandatoryGlobalConfig(c *Config) error {
	switch {
	case len(c.Region) == 0 && len(c.Profile) == 0:
		c.Region = DefaultRegion
		break
	case len(c.Region) > 0:
		if !AWSRegionRegexp.MatchString(c.Region) {
			return ErrRegionIsInvalid
		}
		break
	}

	return nil
}

// Error is an error type for aws errors
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrQueueMandatory error to pop when queue is not filled
	ErrQueueMandatory = Error("queue_name is mandatory")

	// ErrRegionIsInvalid error to pop when region is filled but incorrect
	ErrRegionIsInvalid = Error("region is not a valid region")

	// ErrTopicIsMandatory error to pop when topic is not filled
	ErrTopicIsMandatory = Error("topic-arn is mandatory")
)
