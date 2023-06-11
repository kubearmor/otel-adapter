package stanza_input_operator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/testutil"
)

func TestInputKubearmor(t *testing.T) {
	cfg := NewConfig()
	cfg.LogFilter = "system"
	cfg.OutputIDs = []string{"output"}

	op, err := cfg.Build(testutil.Logger(t))
	require.NoError(t, err)
	mockOutput := testutil.NewMockOperator("output")
	received := make(chan *entry.Entry)
	mockOutput.On("Process", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		received <- args.Get(1).(*entry.Entry)
	}).Return(nil)

	err = op.SetOutputs([]operator.Operator{mockOutput})
	require.NoError(t, err)
	err = op.Start(testutil.NewMockPersister("test"))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, op.Stop())
	}()

	select {
	case e := <-received:
		require.NotNil(t, e.Body)

	case <-time.After(3 * time.Second):
		require.FailNow(t, "Timed out waiting for entry to be read")
	}
}
