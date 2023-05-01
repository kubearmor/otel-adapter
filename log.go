package kubearmor_receiver

import (
	input_operator "github.com/Chinwendu20/kubearmor_receiver/stanza_input_operator"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
)

const (
	typeStr   = "kubearmor_receiver"
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a factory for kubearmor receiver
func NewFactory() receiver.Factory {
	return adapter.NewFactory(ReceiverType{}, stability)
}

// ReceiverType implements adapter.LogReceiverType
// to create a kubearmor receiver
type ReceiverType struct{}

// Type is the receiver type
func (f ReceiverType) Type() component.Type {
	return typeStr
}

// CreateDefaultConfig creates a config with type and version
func (f ReceiverType) CreateDefaultConfig() component.Config {
	return &KubearmorConfig{
		BaseConfig: adapter.BaseConfig{
			Operators: []operator.Config{},
		},
		InputConfig: *input_operator.NewConfig(),
	}
}

// BaseConfig gets the base config from config, for now
func (f ReceiverType) BaseConfig(cfg component.Config) adapter.BaseConfig {
	return cfg.(*KubearmorConfig).BaseConfig
}

// InputConfig unmarshals the input operator
func (f ReceiverType) InputConfig(cfg component.Config) operator.Config {
	return operator.NewConfig(&cfg.(*KubearmorConfig).InputConfig)
}

// KubearmorConfig defines configuration for the journald receiver
type KubearmorConfig struct {
	adapter.BaseConfig `mapstructure:",squash"`
	InputConfig        input_operator.Config `mapstructure:",squash"`
}
