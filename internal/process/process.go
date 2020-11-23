package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/dispatch-simulator/internal/defs"
)

type process struct {
	config       Config
	processQueue chan defs.Order
}

//Config :
type Config struct {
	CF   string        `yaml:"config"`
	Time time.Duration `yaml:"time"`
}

//NewProcess :
func NewProcess(config Config) *process {
	return &process{
		config:       config,
		processQueue: make(chan defs.Order),
	}
}

func (process *process) Run() error {
	//TODO: Currently statically reading from json file
	//Can be mounted to API in future
	jsonFile, err := os.Open(process.config.CF)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	v, _ := ioutil.ReadFile(jsonFile.Name())
	var order []defs.Order
	json.Unmarshal(v, &order)
	//Queue the order for processing
	process.Queue(order)
	return nil
}

//Prepare: API to set preptime (start cooking)
func (process *process) Prepare(order defs.Order) {
	ticker := time.NewTicker(order.PrepTime)

	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		wait.Done()
		for {
			select {
			case <-ticker.C:
				process.processQueue <- order
			default:
			}
		}
	}()

	wait.Wait()
	return
}

func (process *process) Queue(order []defs.Order) {
	n := len(order)
	if n == 0 {
		return
	}
	limiter := time.Tick(process.config.Time * time.Millisecond)
	for _, o := range order {
		<-limiter
		process.Prepare(o)
	}
}

func (process *process) Stop() {
}

func (process *process) Stopped() <-chan interface{} {
	return nil
}
