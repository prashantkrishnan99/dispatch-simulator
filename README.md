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
