package dispatch

import (
	"github.com/dispatch-simulator/internal/defs"
)

func (dispatch *dispatch) DispatchSink() chan<- defs.Dispatch {
	return dispatch.sinkProcessor
}
