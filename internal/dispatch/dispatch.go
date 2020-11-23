package dispatch

import (
	"math/rand"
	"sync"
	"time"

	"github.com/dispatch-simulator/internal/defs"
	"github.com/dispatch-simulator/internal/helper"
	"go.melnyk.org/mlog"
)

type dispatch struct {
	config        Config
	log           mlog.Logger
	sinkProcessor chan defs.Dispatch
	stopped       chan interface{}
	storage       defs.Store
	queue         defs.QueueStore
}

//Config :
type Config struct {
	DispatchArrivalStart int `yaml:"dispatcharrivalstart"`
	DispatchArrivalEnd   int `yaml:"dispatcharrivalend"`
}

//NewDispatch :
func NewDispatch(config Config, log mlog.Joiner, store defs.Store, queue defs.QueueStore) *dispatch {
	return &dispatch{
		config:        config,
		log:           log.Join("Order Dispatcher"),
		sinkProcessor: make(chan defs.Dispatch),
		stopped:       make(chan interface{}),
		storage:       store,
		queue:         queue,
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
			//Call dispatcher to pickup order after "n" seconds
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
							ev.String("Has arrived at the kitchen at time", time.Now().String())
						})
						//Time when the dispatcher has arrived
						t := time.Now().UnixNano() / 100000
						//On arrival of dispatcher check if order is present to be dispatched
						//For matched algorithm, match order to dispatch
						if d.Algo == defs.Matched {
							//Check if the order is ready
							//Is maintained as "ready_<orderid> : true"
							if ready := dispatch.storage.Get(defs.ORDERREADY + d.OrderID); ready != nil {
								//Order is ready to be picked up in the kitchen
								dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
									ev.String("Dispatcher ", d.DispatchID)
									ev.String("has picked up order ", d.OrderID)
									ev.String("from the kitchen", "")
								})
								orderreadytime := ready.(int64)
								//Order is ready to be picked up in the kitchen
								dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
									ev.Int("Average wait time for Dispatcher to pick up the order is ", int(helper.Abs(time.Now().UnixNano()/1000000-orderreadytime)))
								})

								//cleanup dispatch and order queue and map
								dispatch.storage.Delete(defs.ORDERREADY + d.OrderID)
								dispatch.storage.Delete(defs.DISPATCHREADY + d.DispatchID)
								dispatch.storage.Delete(d.OrderID)
								//print statistics
							} else {
								dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
									ev.String("Dispatcher ", d.DispatchID)
									ev.String("is ready and waiting for order ", d.OrderID)
									ev.String("to be ready", "")
								})
								//Set dispatch is ready for pickup when order is prepared and ready
								dispatch.storage.Insert(defs.DISPATCHREADY+d.DispatchID, t)
							}
						}
						dispatch.storage.Insert(d.DispatchID, d.OrderID)
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
}

func (dispatch *dispatch) Stopped() <-chan interface{} {
	return dispatch.stopped
}
