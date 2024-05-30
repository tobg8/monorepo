package awsx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckMandatoryGlobalConfig(t *testing.T) {
	t.Run("should fill region with default region if not setted", func(t *testing.T) {
		c := Config{}
		err := CheckMandatoryGlobalConfig(&c)
		assert.NoError(t, err)
		assert.Equal(t, DefaultRegion, c.Region)
	})

	t.Run("should fill profile", func(t *testing.T) {
		c := Config{
			Profile: "test",
		}
		err := CheckMandatoryGlobalConfig(&c)
		assert.NoError(t, err)
		assert.Equal(t, "test", c.Profile)
	})

	t.Run("should pop a ErrRegionIsInvalid if region is incorrectly filled", func(t *testing.T) {
		c := Config{
			Region: "not a real region",
		}
		err := CheckMandatoryGlobalConfig(&c)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrRegionIsInvalid)
	})
}

func TestBuilderOptions(t *testing.T) {
	t.Run("credentials cache without credentials provider", func(t *testing.T) {
		cfg, err := NewAWSConfigBuilder(&Config{}).
			WithCredentialsCache(CredentialsCacheConfig{}).
			Build(context.Background())
		require.Error(t, err)
		assert.Zero(t, cfg)
	})

	t.Run("assume role with static credentials", func(t *testing.T) {
		cfg, err := NewAWSConfigBuilder(&Config{}).
			WithStaticCredentials("test", "test", "").
			WithAssumeRole("test").
			Build(context.Background())
		require.Error(t, err)
		assert.Zero(t, cfg)
	})

	t.Run("static credentials with assume role", func(t *testing.T) {
		cfg, err := NewAWSConfigBuilder(&Config{}).
			WithAssumeRole("test").
			WithStaticCredentials("test", "test", "").
			Build(context.Background())
		require.Error(t, err)
		assert.Zero(t, cfg)
	})
}
