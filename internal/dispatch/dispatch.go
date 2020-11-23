package dispatch

import (
	"math/rand"
	"sync"
	"time"

	"github.com/dispatch-simulator/internal/defs"
	"go.melnyk.org/mlog"
)

type dispatch struct {
	config        Config
	log           mlog.Logger
	sinkProcessor chan defs.Dispatch
	stopped       chan interface{}
	storage       defs.Store
}

//Config :
type Config struct {
	DispatchArrivalStart int `yaml:"dispatcharrivalstart"`
	DispatchArrivalEnd   int `yaml:"dispatcharrivalend"`
}

//NewDispatch :
func NewDispatch(config Config, log mlog.Joiner, store defs.Store) *dispatch {
	return &dispatch{
		config:        config,
		log:           log.Join("Order Dispatcher"),
		sinkProcessor: make(chan defs.Dispatch),
		stopped:       make(chan interface{}),
		storage:       store,
	}
}

func (dispatch *dispatch) Run() error {
	dispatch.log.Verbose("Dispatcher Service Started")
	dispatch.Receive()
	return nil
}

func (dispatch *dispatch) Receive() {
	for {
		select {
		case d := <-dispatch.sinkProcessor:
			dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
				ev.String("Order Received", "")
				ev.String("Order ID", d.OrderID)
				ev.String("Dispacth ID", d.DispatchID)
			})
			rand.Seed(time.Now().UnixNano())
			x := rand.Intn(dispatch.config.DispatchArrivalEnd-dispatch.config.DispatchArrivalStart+1) + dispatch.config.DispatchArrivalStart
			ticker := time.NewTicker(time.Duration(x) * time.Second)

			var wait sync.WaitGroup
			wait.Add(1)

			go func(ticker *time.Ticker, d defs.Dispatch) {
				wait.Done()
				for {
					select {
					case <-ticker.C:
						ticker.Stop()
						dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
							ev.String("Dispatcher ", d.DispatchID)
							ev.String("Arrived at the kitchen at time", time.Now().String())
						})
					}
				}
			}(ticker, d)
		case <-dispatch.stopped:
			return
		}
	}

}

func (dispatch *dispatch) Stop() {

	select {
	case <-dispatch.stopped:
		return
	default:
	}

	dispatch.log.Verbose("Dispatcher Service Stopped")
	dispatch.log.Info("Processor service stopped")
}

func (dispatch *dispatch) Stopped() <-chan interface{} {
	return dispatch.stopped
}
