package stanza_input_operator

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/testutil"
	"testing"
)

func TestNewClient(t *testing.T) {

	testFilters := []string{"all", "policy", "system"}

	for _, testFilter := range testFilters {

		cfg := NewConfig()
		cfg.LogFilter = testFilter
		input_operator, err := cfg.InputConfig.Build(testutil.Logger(t))
		fd, err := NewClient(input_operator, *cfg)
		if err != nil {
			t.Errorf("Failed to create new client: %s", err)
		}

		if fd.server != cfg.Endpoint {
			t.Errorf("Server address does not match configuration")
		}

		if fd.client == nil {
			t.Errorf("Client is not initialized")
		}

		// Test log filter
		if fd.msgStream == nil {
			t.Errorf("Message stream is not initialized")
		}

		if fd.alertStream == nil && (testFilter == "policy" || testFilter == "all") {
			t.Errorf("Alert stream is not initialized")
		}

		if fd.logStream == nil && (testFilter == "system" || testFilter == "all") {
			t.Errorf("Log stream is not initialized")
		}

	}
}

func TestDoHealthCheck(t *testing.T) {

}
