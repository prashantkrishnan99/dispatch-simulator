package process

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/dispatch-simulator/internal/defs"
	"go.melnyk.org/mlog"
)

type process struct {
	config       Config
	processQueue chan defs.Order
	log          mlog.Logger
	dispatchsink chan<- defs.Dispatch
	storage      defs.Store
}

//DispatchSink :
type DispatchSink interface {
	DispatchSink() chan<- defs.Dispatch
}

//Config :
type Config struct {
	CF   string        `yaml:"config"`
	Time time.Duration `yaml:"time"`
	Mode int           `yaml:"mode"`
}

//NewProcess :
func NewProcess(config Config, log mlog.Joiner, dispatchsink DispatchSink, store defs.Store) *process {
	return &process{
		config:       config,
		processQueue: make(chan defs.Order),
		log:          log.Join("Order Processor"),
		dispatchsink: dispatchsink.DispatchSink(),
		storage:      store,
	}
}

func (process *process) CreateDispatchID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (process *process) Listen() {
	for {
		select {
		case o := <-process.processQueue:
			process.log.Event(mlog.Verbose, func(ev mlog.Event) {
				ev.String("msg", "Order Ready!!!")
				ev.String("Order ID", o.ID)
				ev.String("Order Name", o.Name)
				ev.String("Order ready at", time.Now().String())
			})
		}
	}
}

func (process *process) Run() error {
	process.log.Verbose("Starting Order Processor Service")
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
	//listen in order queue
	go process.Listen()
	//Based on Mode, select the algorithm you want to run
	//0: Stands for Matched ALgorithm
	//1: Stands for FiFO ALgorithm
	process.SelectAlgo()
	//Queue the order for processing
	go process.Queue(order)
	return nil
}

//Prepare: API to set preptime (start cooking)
func (process *process) Prepare(order defs.Order) {
	ticker := time.NewTicker(order.PrepTime * time.Second)

	var wait sync.WaitGroup
	wait.Add(1)

	go func(ticker *time.Ticker) {
		wait.Done()
		for {
			select {
			case <-ticker.C:
				ticker.Stop()
				process.processQueue <- order
			}
		}
	}(ticker)

	wait.Wait()
	return
}

//Dispatch: API to invoke a dispatcher to collect the order on complete
func (process *process) Dispatch(order defs.Order) {
	process.dispatchsink <- defs.Dispatch{
		DispatchID: process.CreateDispatchID(),
		OrderID:    order.ID,
	}
}

func (process *process) Queue(order []defs.Order) {
	process.log.Event(mlog.Verbose, func(ev mlog.Event) {
		ev.String("msg", "Start receiving orders")
		ev.String("time", time.Now().String())
	})
	n := len(order)
	if n == 0 {
		return
	}
	limiter := time.Tick(process.config.Time * time.Millisecond)
	for _, o := range order {
		<-limiter
		process.log.Event(mlog.Verbose, func(ev mlog.Event) {
			ev.String("Order Received", "")
			ev.String("ID", o.ID)
			ev.String("Name", o.Name)
			ev.String("ready at", time.Now().String())
		})
		//On receiving an order, we have to do 2 things
		//dispatch: Invoke a dispatcher to pick up the order when ready
		//prepare: Prepare the order
		process.Dispatch(o)
		process.Prepare(o)
	}
}

func (process *process) Stop() {
	process.log.Verbose("Stopping Order Process Service")
}

func (process *process) Stopped() <-chan interface{} {
	return nil
}
