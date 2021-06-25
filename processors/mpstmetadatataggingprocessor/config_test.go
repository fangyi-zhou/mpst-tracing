package mpstmetadatataggingprocessor

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configtest"
	"path"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factories.Processors[typeStr] = NewFactory()

	cfg, err := configtest.LoadConfigAndValidate(
		path.Join(".", "testdata", "config.yaml"),
		factories,
	)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	mpstConfig := cfg.Processors[config.NewID(typeStr)].(*Config)
	require.NotNil(t, mpstConfig)

	roles := mpstConfig.Roles
	client, exists := roles["client"]
	assert.True(t, exists)
	assert.Equal(t, "frontend", client.Name)

}

func TestConfigValidateEmpty(t *testing.T) {
	cfg := Config{
		ProcessorSettings: config.ProcessorSettings{},
		Roles:             make(map[string]metadataTag),
	}

	err := cfg.Validate()
	assert.Error(t, err)
}

func TestConfigValidateDupe(t *testing.T) {
	cfg := Config{
		ProcessorSettings: config.ProcessorSettings{},
		Roles:             make(map[string]metadataTag),
	}

	cfg.Roles["client"] = metadataTag{Name: "client"}
	cfg.Roles["server"] = metadataTag{Name: "client"}

	err := cfg.Validate()
	assert.Error(t, err)
}

func TestConfigValidateCorrect(t *testing.T) {
	cfg := Config{
		ProcessorSettings: config.ProcessorSettings{},
		Roles:             make(map[string]metadataTag),
	}

	cfg.Roles["client"] = metadataTag{Name: "client"}

	err := cfg.Validate()
	assert.NoError(t, err)
}
