package lpc

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/lpc"
	"github.com/josephbudd/cwt/domain/lpc/message"
)

// Sending is a chanel that sends to the renderer.
type Sending chan interface{}

// Receiving is a channel that receives from the renderer.
type Receiving chan interface{}

// EOJer signals lpc go process to quit.
type EOJer interface {
	Signal()
	NewEOJ() (ch chan struct{})
	Release()
}

// EOJing is has a channel with a dynami size.
// The channel signals go routines to stop.
// EOJing implements EOJer.
type EOJing struct {
	ch          chan struct{}
	count       int
	signaled    bool
	countMutex  *sync.Mutex
	signalMutex *sync.Mutex
}

var (
	send    Sending
	receive Receiving
	eoj     EOJer
)

func init() {
	send = make(chan interface{}, 1024)
	receive = make(chan interface{})
	eoj = &EOJing{
		ch:          make(chan struct{}, 1024),
		count:       0,
		countMutex:  &sync.Mutex{},
		signalMutex: &sync.Mutex{},
	}
}

// Channels returns the renderer connection channels.
func Channels() (sendChan Sending, receiveChan Receiving, eojChan EOJer) {
	sendChan = send
	receiveChan = receive
	eojChan = eoj
	return
}

// Payload converts unmarshalled msg to the correct marshalled payload.
func (sending Sending) Payload(msg interface{}) (payload []byte, err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "sending.Payload")
		}
	}()

	var bb []byte
	var id uint64
	switch msg := msg.(type) {
	case *message.LogMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 0
	case *message.CheckCopyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 1
	case *message.CheckKeyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 2
	case *message.GetCopyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 3
	case *message.GetKeyCodesMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 4
	case *message.GetKeyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 5
	case *message.GetTextToCopyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 6
	case *message.GetTextToKeyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 7
	case *message.KeyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 8
	case *message.MetronomeMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 9
	case *message.UpdateCopyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 10
	case *message.UpdateKeyCodeMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 11
	case *message.UpdateKeyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 12
	default:
		bb = []byte("Unknown!")
		id = 999
	}
	pl := &lpc.Payload{
		ID:    id,
		Cargo: bb,
	}
	payload, err = json.Marshal(pl)
	return
}

// Cargo returns a marshalled payload's unmarshalled cargo.
func (receiving Receiving) Cargo(payloadbb []byte) (cargo interface{}, err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "receiving.Cargo")
		}
	}()

	payload := lpc.Payload{}
	if err = json.Unmarshal(payloadbb, &payload); err != nil {
		return
	}
	switch payload.ID {
	case 0:
		msg := &message.LogRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 1:
		msg := &message.CheckCopyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 2:
		msg := &message.CheckKeyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 3:
		msg := &message.GetCopyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 4:
		msg := &message.GetKeyCodesRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 5:
		msg := &message.GetKeyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 6:
		msg := &message.GetTextToCopyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 7:
		msg := &message.GetTextToKeyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 8:
		msg := &message.KeyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 9:
		msg := &message.MetronomeRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 10:
		msg := &message.UpdateCopyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 11:
		msg := &message.UpdateKeyCodeRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 12:
		msg := &message.UpdateKeyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	default:
		errMsg := fmt.Sprintf("no case found for payload id %d", payload.ID)
		err = errors.New(errMsg)
	}
	return
}

// Signal sends on the eoj channel signaling lpc go funcs to quit.
func (eoj EOJing) Signal() {
	eoj.signalMutex.Lock()
	if !eoj.signaled {
		eoj.signaled = true
		end := struct{}{}
		for i := 0; i < eoj.count; i++ {
			eoj.ch <- end
		}
	}
	eoj.signalMutex.Unlock()
}

// NewEOJ returns a new eoj channel and increments the usage count.
func (eoj EOJing) NewEOJ() (ch chan struct{}) {
	eoj.countMutex.Lock()
	eoj.count++
	ch = eoj.ch
	eoj.countMutex.Unlock()
	return
}

// Release decrements the usage count.
// Call this at the end of your lpc handler func.
func (eoj EOJing) Release() {
	eoj.countMutex.Lock()
	if eoj.count > 0 {
		eoj.count--
	}
	eoj.countMutex.Unlock()
}
