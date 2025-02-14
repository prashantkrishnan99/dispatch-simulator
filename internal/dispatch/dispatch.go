package dispatch

import (
	"math/rand"
	"reflect"
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
	orderqueue    defs.QueueStore
	dispatchqueue defs.QueueStore
	stats         defs.Stats
}

//Config :
type Config struct {
	DispatchArrivalStart int `yaml:"dispatcharrivalstart"`
	DispatchArrivalEnd   int `yaml:"dispatcharrivalend"`
}

//NewDispatch :
func NewDispatch(config Config, log mlog.Joiner, store defs.Store, orderqueue defs.QueueStore, dispatchqueue defs.QueueStore, stats defs.Stats) *dispatch {
	return &dispatch{
		config:        config,
		log:           log.Join("Order Dispatcher"),
		sinkProcessor: make(chan defs.Dispatch),
		stopped:       make(chan interface{}),
		storage:       store,
		orderqueue:    orderqueue,
		dispatchqueue: dispatchqueue,
		stats:         stats,
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
						t := time.Now()
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
								orderreadytime := ready.(time.Time)
								//Order is ready to be picked up in the kitchen
								absoluteTime := t.Sub(orderreadytime)
								dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
									ev.Int("Average wait time for Dispatcher to pick up the order is ", int(absoluteTime.Milliseconds()))
								})

								//Calculate stats
								dispatch.stats.IncrOrdersProcessed()
								dispatch.stats.IncrTotalTime(int(absoluteTime.Milliseconds()))
								dispatch.stats.CalculateAverage()

								dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
									ev.Int("Running Average wait time for Dispatcher to pick up the order is ", dispatch.stats.GetAVerageTime())
								})

								//cleanup dispatch and order queue and map
								dispatch.storage.Delete(defs.ORDERREADY + d.OrderID)
								dispatch.storage.Delete(defs.DISPATCHREADY + d.DispatchID)
								dispatch.storage.Delete(d.OrderID)
							} else {
								dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
									ev.String("Dispatcher ", d.DispatchID)
									ev.String("is ready and waiting for order ", d.OrderID)
									ev.String("to be ready", "")
								})
								//Set dispatch is ready for pickup when order is prepared and ready
								dispatch.storage.Insert(defs.DISPATCHREADY+d.DispatchID, t)
							}
						} else if d.Algo == defs.Fifo {
							//Add dispatcher to dispatch queue
							//check if order queue is not empty
							//if empty wait in dispatch queue
							//if not empty, assign dispatch to any order
							//for our use case we can consider first order to be picked
							//in above scenario, since it may get spoiled if kept for long
							dispatch.dispatchqueue.Enqueue(d.DispatchID)
							if dispatch.orderqueue.Size() != 0 {
								order := dispatch.orderqueue.Dequeue()
								castvar := reflect.ValueOf(*order)
								if ready := dispatch.storage.Get(defs.ORDERREADY + castvar.String()); ready != nil {
									//Order is ready to be picked up in the kitchen
									dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
										ev.String("Dispatcher ", d.DispatchID)
										ev.String("has picked up order ", d.OrderID)
										ev.String("from the kitchen", "")
									})
									orderreadytime := ready.(time.Time)
									//Order is ready to be picked up in the kitchen
									absoluteTime := t.Sub(orderreadytime)
									dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
										ev.Int("Average wait time for Dispatcher to pick up the order is ", int(absoluteTime.Milliseconds()))
									})

									//Calculate stats
									dispatch.stats.IncrOrdersProcessed()
									dispatch.stats.IncrTotalTime(int(absoluteTime.Milliseconds()))
									dispatch.stats.CalculateAverage()

									dispatch.log.Event(mlog.Verbose, func(ev mlog.Event) {
										ev.Int("Running Average wait time for Dispatcher to pick up the order is ", dispatch.stats.GetAVerageTime())
									})

									//cleanup dispatch and order queue and map
									dispatch.storage.Delete(defs.ORDERREADY + d.OrderID)
									dispatch.storage.Delete(defs.DISPATCHREADY + d.DispatchID)

									dispatch.dispatchqueue.Dequeue()
								}
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
