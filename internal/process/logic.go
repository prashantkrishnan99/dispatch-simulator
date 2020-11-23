package process

func (process *process) SelectAlgo() {
	if process.config.Mode == 0 {
		go process.Matched()
	} else if process.config.Mode == 1 {
		go process.Fifo()
	}
}
