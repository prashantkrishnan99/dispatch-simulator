package runner

import (
	"fmt"
	"sync"

	"github.com/dispatch-simulator/internal/defs"
	"github.com/dispatch-simulator/internal/dispatch"
	"github.com/dispatch-simulator/internal/process"
	"github.com/dispatch-simulator/internal/stats"
	"go.melnyk.org/mlog"
)

//Runner :
type Runner interface {
	Run() error
	Stop()
	Stopped() <-chan interface{}
}

//Service :
type Service interface {
	Run() error
	Stop()
	Stopped() <-chan interface{}
}

//Config :
type Config struct {
	Dispatch dispatch.Config `yaml:"dispatch"`
	Process  process.Config  `yaml:"process"`
}

type runner struct {
	config    Config
	wg        sync.WaitGroup
	log       mlog.Logger
	logjoiner mlog.Joiner
	store     defs.Store
	dispatch  Service
	process   Service
}

//NewRunner :
func NewRunner(config Config, log mlog.Joiner) Runner {
	runner := &runner{
		config:    config,
		log:       log.Join("runner"),
		logjoiner: log,
	}

	store := NewStorage()
	orderqueue := NewQueue()
	dispatchqueue := NewQueue()

	stats := stats.NewStats()

	dispatch := dispatch.NewDispatch(config.Dispatch, log, store, orderqueue, dispatchqueue, stats)
	runner.dispatch = dispatch

	process := process.NewProcess(config.Process, log, dispatch, store, orderqueue, dispatchqueue, stats)
	runner.process = process

	return runner
}

func (runner *runner) Run() error {
	if err := runner.process.Run(); err != nil {
		return err
	}
	if err := runner.dispatch.Run(); err != nil {
		return err
	}

	runner.wg.Add(1)
	processorders := func() {
		defer runner.wg.Done()
	forever:
		for {
			select {
			case <-runner.dispatch.Stopped():
				fmt.Println("Dispatch service stopped")
				break forever
			case <-runner.process.Stopped():
				fmt.Println("Process service stopped")
				break forever
			}
		}
		runner.dispatch.Stop()
		runner.process.Stop()
	}
	go processorders()

	return nil
}

func (runner *runner) Stop() {
	runner.dispatch.Stop()
}

func (runner *runner) Stopped() <-chan interface{} {
	return nil
}
