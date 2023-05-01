package kubearmor_receiver

import (
	"context"
	"fmt"
	pb "github.com/kubearmor/KubeArmor/protobuf"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"google.golang.org/grpc"
	"math/rand"
)

// Feeder Structure
type Feeder struct {
	// flag
	Running bool

	// server
	server string

	// connection
	conn *grpc.ClientConn

	// client
	client pb.LogServiceClient

	// messages
	msgStream pb.LogService_WatchMessagesClient

	// alerts
	alertStream pb.LogService_WatchAlertsClient

	// logs
	logStream pb.LogService_WatchLogsClient
}

// NewClient Function
func NewClient(operator helper.InputOperator, cfg Config) (*Feeder, error) {
	fd := &Feeder{}

	fd.Running = true

	fd.server = cfg.Endpoint

	logFilter := cfg.LogFilter

	conn, err := grpc.Dial(fd.server, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to kubearmor relay server (%s)\n", err.Error())
	}
	fd.conn = conn

	fd.client = pb.NewLogServiceClient(fd.conn)

	operator.Logger().Info("Connected to kubearmor relay server")

	msgIn := pb.RequestMessage{}
	msgIn.Filter = ""

	if logFilter == "all" || logFilter == "kubearmorLogs" {
		msgStream, err := fd.client.WatchMessages(context.Background(), &msgIn)
		if err != nil {
			return nil, fmt.Errorf("Failed to call WatchMessages() (%s)\n", err.Error())
		}
		fd.msgStream = msgStream
		operator.Logger().Info("Watching kubearmor logs......")
	}

	alertIn := pb.RequestMessage{}
	alertIn.Filter = logFilter

	if alertIn.Filter == "all" || alertIn.Filter == "policy" {
		alertStream, err := fd.client.WatchAlerts(context.Background(), &alertIn)
		if err != nil {
			return nil, fmt.Errorf("Failed to call WatchAlerts() (%s)\n", err.Error())
		}
		fd.alertStream = alertStream
		operator.Logger().Info("Watching kubearmor alerts.....")
	}

	logIn := pb.RequestMessage{}
	logIn.Filter = logFilter

	if logIn.Filter == "all" || logIn.Filter == "system" {
		logStream, err := fd.client.WatchLogs(context.Background(), &logIn)
		if err != nil {
			return nil, fmt.Errorf("Failed to call WatchLogs() (%s)\n", err.Error())
		}
		fd.logStream = logStream
		operator.Logger().Info("Watching system insights received from kubearmor ......")
	}

	return fd, nil
}

func (fd *Feeder) recvMsg(operator *Input) (error, *pb.Message) {
	res, err := fd.msgStream.Recv()
	if err != nil {
		return fmt.Errorf("Failed to receive a message (%s)\n", err.Error()), nil
	}
	return nil, res
}

func (fd *Feeder) recvAlerts(operator *Input) (error, *pb.Alert) {
	res, err := fd.alertStream.Recv()
	if err != nil {
		return fmt.Errorf("Failed to receive an alert (%s)\n", err.Error()), res
	}
	return nil, res
}

func (fd *Feeder) recvLogs(operator *Input) (error, *pb.Log) {
	res, err := fd.logStream.Recv()
	if err != nil {
		return fmt.Errorf("Failed to receive an log (%s)\n", err.Error()), nil
	}
	return nil, res
}

func (fd *Feeder) doHealthCheck(operator *Input) bool {
	// #nosec
	randNum := rand.Int31()

	// send a nonce
	nonce := pb.NonceMessage{Nonce: randNum}
	res, err := fd.client.HealthCheck(context.Background(), &nonce)
	if err != nil {
		operator.Logger().Info("Failed to call HealthCheck() (%s)\n", err.Error())
		return false
	}

	// check nonce
	if randNum != res.Retval {
		return false
	}

	return true
}

func (fd *Feeder) DestroyClient() error {
	if err := fd.conn.Close(); err != nil {
		return err
	}
	return nil
}
