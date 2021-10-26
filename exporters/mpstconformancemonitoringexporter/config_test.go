package mpstconformancemonitoringexporter

import (
	"go.opentelemetry.io/collector/config"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
)

func TestLoadConfigLts(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factories.Exporters[typeStr] = NewFactory()

	cfg, err := configtest.LoadConfigAndValidate(
		path.Join(".", "testdata", "config.yaml"),
		factories,
	)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	lts := config.NewComponentIDWithName(typeStr, "lts")

	mpstConfig := cfg.Exporters[lts].(*Config)
	assert.Equal(t, "gtype_lts", mpstConfig.SemanticModelType)
	assert.Equal(t, "gtype.sexp", mpstConfig.GlobalTypeSexpFileName)
}

func TestLoadConfigPedro(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factories.Exporters[typeStr] = NewFactory()

	cfg, err := configtest.LoadConfigAndValidate(
		path.Join(".", "testdata", "config.yaml"),
		factories,
	)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	pedro := config.NewComponentIDWithName(typeStr, "pedro")

	mpstConfig := cfg.Exporters[pedro].(*Config)
	assert.Equal(t, "gtype_pedro", mpstConfig.SemanticModelType)
	assert.Equal(t, "pedrolib.so", mpstConfig.PedroSoFileName)
	assert.Equal(t, "MyProto.nuscr", mpstConfig.ProtocolFileName)
	assert.Equal(t, "MyProto", mpstConfig.ProtocolName)
}
