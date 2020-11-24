# Dispatch-Simulator


## Features

- Dispatch-Simulator is a Golang library for dispatch management of food orders

## Usage
``` 
 $ dispatch --help
 
      Usage:
		  dispatch [flags]
		  dispatch [command]

	  Available Commands:
		  config      Config information
		  help        Help about any command
		  run         Run service

	  Flags:
		  -h, --help   help for dispatch
```

## Build Binary 
go build -a -o bin/dispatch ./cmd/dispatch


### Basic usage

```bash
$ bin/dispatch run
```

To attach configuration to run use `dispatch.yaml` file at any time, simply run:

```bash
$ bin/dispatch run dispatch.yaml
```

### System Design

Use https://sequencediagram.org/ and copy below paragraphs to view the design

1. Matched Algorithm

title Matched Dispatch Algorithm

Process->Process:Receive order
Process->Process:Generate dispatch id for an order
Process->Process:Start cooking
Process->Storage:Store order id to dispatch id map
Process->Dispatch:Pass Dispatch object (order info + courier info)
Dispatch->Dispatch:Start Dispatch Order invocation
Process->Process:Async receive cooking completed callback
Process->Process:On receiving the callback; check if corresponding dispatch id has entered kitchen; If present then order can be picked up; if not wait in the queue
Process->Storage:Store order id received entry
Process->Storage:On above 2 step success, Cleanup order and dispatch details from storage
Dispatch->Dispatch:Async receive Dispatcher in the kitchen callback
Dispatch->Dispatch:On receiving the callback; check if corresponding order id has been completed in kitchen; If present then order can be picked up; if not wait in the queue
Dispatch->Storage:Store dispatch id received entry
Dispatch->Storage:On above 2 step success, Cleanup order and dispatch details from storage

2. FIFO Algorithm

title FIFO Dispatch Algorithm

Process->Process:Receive order
Process->Process:Generate dispatch id for an order
Process->Process:Start cooking
Process->Storage:Enqueue Order id in the order queue
Process->Dispatch:Pass Dispatch object (order info + courier info)
Dispatch->Dispatch:Start Dispatch Order invocation
Process->Process:Async receive cooking completed callback
Process->Process:On receiving the callback; check if the dispatch queue length is greater than 0 (i.e if any dispatchers are available); if present dequeue from the dispatch queue and assign to the order; if not then wait in the queue
Process->Storage:Store order id received entry
Process->Storage:On above 2 step success, Cleanup order and dispatch details from storage
Dispatch->Dispatch:Async receive Dispatcher in the kitchen callback
Dispatch->Dispatch:On receiving the callback; check if the order queue length is greater than 0 (i.e if any dispatchers are available); if present dequeue from the order queue and assign to the dispatcher; if not then wait in the queue
Dispatch->Storage:Store dispatch id received entry
Dispatch->Storage:On above 2 step success, Cleanup order and dispatch details from storage