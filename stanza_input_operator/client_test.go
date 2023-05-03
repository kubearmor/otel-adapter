package kubearmor_receiver

//
//import (
//	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/testutil"
//	"testing"
//)
//
//func TestNewClient(t *testing.T) {
//	mockOperator := testutil.
//	cfg := Config{
//		Endpoint:  "localhost:50051",
//		LogFilter: "all",
//	}
//
//	fd, err := NewClient(operator, cfg)
//	if err != nil {
//		t.Errorf("Failed to create new client: %s", err)
//	}
//
//	if !fd.Running {
//		t.Errorf("Client should be running, but Running is false")
//	}
//
//	if fd.server != cfg.Endpoint {
//		t.Errorf("Server address does not match configuration")
//	}
//
//	if fd.client == nil {
//		t.Errorf("Client is not initialized")
//	}
//
//	// Test log filter
//	if fd.msgStream == nil {
//		t.Errorf("Message stream is not initialized")
//	}
//
//	if fd.alertStream == nil {
//		t.Errorf("Alert stream is not initialized")
//	}
//
//	if fd.logStream == nil {
//		t.Errorf("Log stream is not initialized")
//	}
//}
