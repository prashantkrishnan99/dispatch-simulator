package runner

import (
	"fmt"
	"sync"

	"github.com/dispatch-simulator/internal/dispatch"
	"github.com/dispatch-simulator/internal/logic"
	"github.com/dispatch-simulator/internal/process"
)

//Runner :
type Runner interface {
	Run() error
	Stop()
}

//Service :
type Service interface {
	Run() error
	Stop()
	Stopped() <-chan interface{}
}

//Config :
type Config struct {
	Logic    logic.Config    `yaml:"logic"`
	Dispatch dispatch.Config `yaml:"dispatch"`
	Process  process.Config  `yaml:"process"`
}

type runner struct {
	config Config
	wg     sync.WaitGroup

	dispatch Service
	process  Service
}

//NewRunner :
func NewRunner(config Config) Runner {
	runner := &runner{
		config: config,
	}

	dispatch := dispatch.NewDispatch(config.Dispatch)
	runner.dispatch = dispatch

	process := process.NewProcess(config.Process)
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
