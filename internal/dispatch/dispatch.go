package dispatch

type dispatch struct {
	config Config
}

//Config :
type Config struct {
}

//NewDispatch :
func NewDispatch(config Config) *dispatch {
	return &dispatch{
		config: config,
	}
}

func (dispatch *dispatch) Run() error {
	return nil
}

func (dispatch *dispatch) Stop() {
}

func (dispatch *dispatch) Stopped() <-chan interface{} {
	return nil
}
