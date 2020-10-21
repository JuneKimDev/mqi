package mqi

type store struct {
	reqChan    chan chan Channel
	updateChan chan Channel
}

// Store interface
type Store interface {
	ReqChan() chan chan Channel
	UpdateChan() chan Channel
	WithReqChan(ch chan chan Channel) Store
	WithUpdateChan(ch chan Channel) Store
	run()
}

// initStore ...
func initStore() Store {
	mStore := (&store{}).
		WithReqChan(make(chan chan Channel)).
		WithUpdateChan(make(chan Channel))
	go mStore.run()
	return mStore
}

// Getters and withers
func (st *store) ReqChan() chan chan Channel { return st.reqChan }
func (st *store) UpdateChan() chan Channel   { return st.updateChan }
func (st *store) WithReqChan(ch chan chan Channel) Store {
	st.reqChan = ch
	return st
}
func (st *store) WithUpdateChan(ch chan Channel) Store {
	st.updateChan = ch
	return st
}

// run store (runs in goroutine forever)
func (st *store) run() {
	// storage
	state := NewChannel(st)

	for {
		select {
		// Redux selector
		case resChan := <-st.reqChan: // Request to read the state
			resChan <- state // State is returen via resChan
		// Redux reducer
		case newState := <-st.updateChan: // Request to update the state
			state = newState
		}
	}
}
