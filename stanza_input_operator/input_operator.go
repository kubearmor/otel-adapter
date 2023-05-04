package stanza_input_operator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"
)

const operatorType = "kubearmor_input"

func init() {
	operator.Register(operatorType, func() operator.Builder { return NewConfig() })
}

// NewConfig creates a new input config with default values
func NewConfig() *Config {
	return NewConfigWithID(operatorType)
}

// NewConfigWithID creates a new input config with default values
func NewConfigWithID(operatorID string) *Config {
	var gRPC string
	if val, ok := os.LookupEnv("KUBEARMOR_SERVICE"); ok {
		gRPC = val
	} else {
		gRPC = "localhost:32767"
	}
	return &Config{
		InputConfig: helper.NewInputConfig(operatorID, operatorType),
		Endpoint:    gRPC,
		LogFilter:   "all",
	}
}

// Config is the configuration of a kubearmor input operator
type Config struct {
	helper.InputConfig `mapstructure:",squash"`

	Endpoint  string `mapstructure:"endpoint,omitempty"`
	LogFilter string `mapstructure:"logfilter,omitempty"`
}

// Build will build a kubearmor input operator from the supplied configuration
func (cfg Config) Build(logger *zap.SugaredLogger) (operator.Operator, error) {
	inputOperator, err := cfg.InputConfig.Build(logger)
	if err != nil {
		return nil, err
	}

	logClient, err := NewClient(inputOperator, cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to create kubearmor relay client: %s", err.Error())
	}
	inputOperator.Logger().Infof("Created kubearmor relay client (%s)\n", cfg.Endpoint)
	return &Input{
		InputOperator:   inputOperator,
		json:            jsoniter.ConfigFastest,
		Config:          cfg,
		kubearmorClient: logClient,
	}, nil
}

// Input is an operator that processes kubearmor logs
type Input struct {
	helper.InputOperator
	Config
	kubearmorClient *Feeder
	json            jsoniter.API
	cancel          context.CancelFunc
	wg              sync.WaitGroup
}

// Start will start generating log entries.
func (operator *Input) Start(_ operator.Persister) error {
	ctx, cancel := context.WithCancel(context.Background())
	operator.cancel = cancel

	logfilter := operator.LogFilter
	logClient := operator.kubearmorClient

	// do healthcheck
	if ok := logClient.doHealthCheck(operator); !ok {
		return fmt.Errorf("Failed to check the liveness of the gRPC server")
	}
	operator.Logger().Info("Checked the liveness of the gRPC server")

	if logfilter == "all" || logfilter == "kubearmorLogs" {
		// watch messages
		operator.wg.Add(1)
		go func() {
			defer operator.wg.Done()
			for logClient.Running {
				err, log := logClient.recvMsg(operator)
				if err != nil {
					operator.Logger().Errorf("%s", err.Error())
					return
				}
				arr, _ := json.Marshal(log)
				str := string(arr)
				entry, err := operator.parseLogEntry([]byte(str))
				if err != nil {
					operator.Logger().Errorf("%s", err.Error())
					return
				}
				operator.Write(ctx, entry)
			}

		}()

	}
	if logfilter == "all" || logfilter == "policy" {
		// watch alerts
		operator.wg.Add(1)
		go func() {
			defer operator.wg.Done()
			for logClient.Running {
				err, log := logClient.recvAlerts(operator)
				if err != nil {
					operator.Logger().Errorf("%s", err.Error())
					return
				}
				arr, _ := json.Marshal(log)
				str := string(arr)
				entry, err := operator.parseLogEntry([]byte(str))
				if err != nil {
					operator.Logger().Errorf("%s", err.Error())
					return
				}
				operator.Write(ctx, entry)
			}

		}()
	}

	if logfilter == "all" || logfilter == "system" {
		// watch logs
		operator.wg.Add(1)
		go func() {
			defer operator.wg.Done()
			for logClient.Running {
				err, log := logClient.recvLogs(operator)
				if err != nil {
					operator.Logger().Errorf("%s", err.Error())
					return
				}
				arr, _ := json.Marshal(log)
				str := string(arr)
				entry, err := operator.parseLogEntry([]byte(str))
				if err != nil {
					operator.Logger().Errorf("%s", err.Error())
					return
				}
				operator.Write(ctx, entry)
			}
		}()
	}

	return nil
}

func (operator *Input) parseLogEntry(line []byte) (*entry.Entry, error) {
	var body map[string]interface{}
	err := operator.json.Unmarshal(line, &body)

	if err != nil {
		return nil, err
	}

	timestamp, ok := body["Timestamp"]

	if !ok {
		return nil, errors.New("log body missing timestamp field")
	}

	timestampFloat, ok := timestamp.(float64)
	if !ok {
		return nil, errors.New("log body field for timestamp is not of string type")
	}

	delete(body, "Timestamp")

	entry, err := operator.NewEntry(body)
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}

	entry.Timestamp = time.Unix(int64(timestampFloat), 0)

	return entry, nil
}

// Stop will stop generating logs.
func (operator *Input) Stop() error {
	logClient := operator.kubearmorClient
	logClient.Running = false
	time.Sleep(2 * time.Second)
	if err := logClient.DestroyClient(); err != nil {

		return fmt.Errorf("Failed to destroy the kubearmor relay gRPC client (%s)\n", err.Error())

	}

	operator.Logger().Info("Destroyed kubearmor relay gRPC client")
	operator.wg.Wait()
	operator.cancel()
	operator.Logger().Info("Stopped kubearmor receiver")
	return nil
}
