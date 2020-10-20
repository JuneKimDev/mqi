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
	Run()
	GetChannel() Channel
	Update(ch Channel)
}

// InitStore ...
func InitStore() Store {
	mStore := (&store{}).
		WithReqChan(make(chan chan Channel)).
		WithUpdateChan(make(chan Channel))
	go mStore.Run()
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

// Run starts store (runs in goroutine forever)
func (st *store) Run() {
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

// GetChannel returns currently Channel
func (st *store) GetChannel() Channel {
	resChan := make(chan Channel)
	st.reqChan <- resChan
	return <-resChan
}

// Update updates State
func (st *store) Update(ch Channel) {
	st.updateChan <- ch
}
