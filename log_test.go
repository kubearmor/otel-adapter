package kubearmor_receiver

import (
	"context"
	input_operator "github.com/Chinwendu20/kubearmor_receiver/stanza_input_operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestLogsReceiverNoErrors(t *testing.T) {
	f := NewFactory()
	lle, err := f.CreateLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), f.CreateDefaultConfig(), nil)
	require.NotNil(t, lle)
	assert.NoError(t, err)

	assert.NoError(t, lle.Start(context.Background(), componenttest.NewNopHost()))

	assert.NoError(t, lle.Shutdown(context.Background()))
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(typeStr, "").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	assert.Equal(t, testdataConfigYaml(), cfg)
}

func testdataConfigYaml() *KubearmorConfig {
	return &KubearmorConfig{
		BaseConfig: adapter.BaseConfig{
			Operators: []operator.Config{},
		},
		InputConfig: func() input_operator.Config {
			c := input_operator.NewConfig()
			c.LogFilter = "policy"
			return *c
		}(),
	}
}
