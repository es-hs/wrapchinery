package wrapchinery

import (
	"context"
	"reflect"
	"time"

	"github.com/RichardKnop/machinery/v2"
	backendsiface "github.com/RichardKnop/machinery/v2/backends/iface"
	"github.com/RichardKnop/machinery/v2/backends/result"
	brokersiface "github.com/RichardKnop/machinery/v2/brokers/iface"
	"github.com/RichardKnop/machinery/v2/config"
	lockiface "github.com/RichardKnop/machinery/v2/locks/iface"
	"github.com/RichardKnop/machinery/v2/tasks"
)

type Server struct {
	*machinery.Server
}

// NewServer creates Server instance
func NewServer(
	cnf *config.Config, brokerServer brokersiface.Broker,
	backendServer backendsiface.Backend, lock lockiface.Lock,
) *Server {
	server := &Server{
		machinery.NewServer(cnf, brokerServer, backendServer, lock),
	}

	return server
}

func (m *Server) WrapSendTask(taskName string, delay time.Duration, retry int, args ...interface{}) (*result.AsyncResult, error) {
	task := GetTaskSignature(taskName, delay, retry, args)
	return m.SendTask(task)
}

func (m *Server) WrapSendTaskWithContext(
	taskName string, ctx context.Context, delay time.Duration, retry int, args ...interface{},
) (*result.AsyncResult, error) {
	task := GetTaskSignature(taskName, delay, retry, args)
	return m.SendTaskWithContext(ctx, task)
}

func GetTaskSignature(taskName string, delay time.Duration, retry int, args ...interface{}) *tasks.Signature {
	task := tasks.Signature{
		Name: taskName,
		Args: []tasks.Arg{},
	}
	if delay > 0 {
		timeETA := time.Now().UTC().Add(delay)
		task.ETA = &timeETA
	}
	task.RetryCount = retry
	if len(args) > 0 {
		task.Args = parseArgs(args)
	}
	return &task
}

func parseArgs(args ...interface{}) []tasks.Arg {
	taskArgs := []tasks.Arg{}
	for k := range args {
		taskArgs = append(taskArgs, tasks.Arg{
			Type:  reflect.TypeOf(args[k]).String(),
			Value: args[k],
		})
	}
	return taskArgs
}