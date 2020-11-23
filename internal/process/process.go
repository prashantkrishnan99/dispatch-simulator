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
	queue        defs.QueueStore
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
func NewProcess(config Config, log mlog.Joiner, dispatchsink DispatchSink, store defs.Store, queue defs.QueueStore) *process {
	return &process{
		config:       config,
		processQueue: make(chan defs.Order),
		log:          log.Join("Order Processor"),
		dispatchsink: dispatchsink.DispatchSink(),
		storage:      store,
		queue:        queue,
	}
}

func (process *process) CreateDispatchID() string {
	dispID := make([]byte, 16)
	rand.Read(dispID)
	return "dispatch_" + fmt.Sprintf("%x-%x-%x-%x-%x", dispID[0:4], dispID[4:6], dispID[6:8], dispID[8:10], dispID[10:])
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

			did := ""
			if process.config.Mode == defs.Matched {
				//Get the dispatcher from the map for the orderid
				//Check if the dispatcher is ready
				//Is maintained as "ready_<dispatch> : true"
				if didforoid := process.storage.Get(o.ID); didforoid != nil {
					//dispatchID
					did = didforoid.(string)
				} else {
					process.log.Event(mlog.Error, func(ev mlog.Event) {
						ev.Error("msg", fmt.Errorf("DispatchID for an order should be present"))
					})
					continue
				}
				//Check if dispatcher is ready to pick up
				//if not set order ready for dispatcher to collect
				//later
				if ready := process.storage.Get(defs.DISPATCHREADY + did); ready != nil {
					//Dispatcher is ready and waiting in the kitchen
					process.log.Event(mlog.Verbose, func(ev mlog.Event) {
						ev.String("Order ", o.ID)
						ev.String("has been picked up by ", did)
						ev.String("from the kitchen", "")
					})
					//cleanup dispatch and order queue and map
					process.storage.Delete(defs.ORDERREADY + o.ID)
					process.storage.Delete(defs.DISPATCHREADY + did)
					process.storage.Delete(o.ID)
					//print statistics
				} else {
					process.log.Event(mlog.Verbose, func(ev mlog.Event) {
						ev.String("Order ", o.ID)
						ev.String("is ready and waiting for dispatcher ", did)
						ev.String("to arrive", "")
					})
					//Set order is ready for pickup when dispatcher arrives later
					process.storage.Insert(defs.ORDERREADY+o.ID, true)
				}
			}
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
	//Queue the order for processing
	go process.Queue(order)
	return nil
}

//Prepare: API to set preptime (start cooking)
func (process *process) Prepare(order defs.Order) {
	ticker := time.NewTicker(order.PrepTime * time.Second)

	var wait sync.WaitGroup
	wait.Add(1)

	go func(ticker *time.Ticker, o defs.Order) {
		wait.Done()
		for {
			select {
			case <-ticker.C:
				ticker.Stop()
				process.processQueue <- o
			}
		}
	}(ticker, order)

	wait.Wait()
	return
}

//Dispatch: API to invoke a dispatcher to collect the order on complete
func (process *process) Dispatch(order defs.Order) {
	//create DispatchID
	dID := process.CreateDispatchID()
	//before invoking dispatcher; maintain orderid to dispatch id map
	process.storage.Insert(order.ID, dID)
	//Send dispatch details to dispatch service
	process.dispatchsink <- defs.Dispatch{
		DispatchID: dID,
		OrderID:    order.ID,
		Algo:       process.config.Mode,
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
		//On receiving an order, Invoke a dispatcher to pick up the order when ready
		process.Dispatch(o)
		//On receiving an order, Prepare the order
		process.Prepare(o)
	}
}

func (process *process) Stop() {
	process.log.Verbose("Stopping Order Process Service")
}

func (process *process) Stopped() <-chan interface{} {
	return nil
}
