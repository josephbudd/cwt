package lpc

import (
	"encoding/json"
	"fmt"

	"github.com/josephbudd/cwt/domain/lpc"
	"github.com/josephbudd/cwt/domain/lpc/message"
)

// Sending is a chanel that sends to the renderer.
type Sending chan interface{}

// Receiving is a channel that receives from the renderer.
type Receiving chan interface{}

var (
	send    Sending
	receive Receiving
)

func init() {
	send = make(chan interface{}, 1024)
	receive = make(chan interface{})
}

// Channels returns the renderer connection channels.
func Channels() (sendChan Sending, receiveChan Receiving) {
	sendChan = send
	receiveChan = receive
	return
}

// Payload converts unmarshalled msg to the correct marshalled payload.
func (sending Sending) Payload(msg interface{}) (payload []byte, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("sending.Payload: %w", err)
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
	case *message.InitMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 1
	case *message.CheckCopyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 2
	case *message.CheckKeyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 3
	case *message.GetCopyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 4
	case *message.GetKeyCodesMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 5
	case *message.GetKeyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 6
	case *message.GetTextToCopyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 7
	case *message.GetTextToKeyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 8
	case *message.KeyMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 9
	case *message.MetronomeMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 10
	case *message.UpdateCopyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 11
	case *message.UpdateKeyCodeMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 12
	case *message.UpdateKeyWPMMainProcessToRenderer:
		if bb, err = json.Marshal(msg); err != nil {
			return
		}
		id = 13
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
			err = fmt.Errorf("receiving.Cargo: %w", err)
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
		msg := &message.InitRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 2:
		msg := &message.CheckCopyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 3:
		msg := &message.CheckKeyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 4:
		msg := &message.GetCopyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 5:
		msg := &message.GetKeyCodesRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 6:
		msg := &message.GetKeyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 7:
		msg := &message.GetTextToCopyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 8:
		msg := &message.GetTextToKeyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 9:
		msg := &message.KeyRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 10:
		msg := &message.MetronomeRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 11:
		msg := &message.UpdateCopyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 12:
		msg := &message.UpdateKeyCodeRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	case 13:
		msg := &message.UpdateKeyWPMRendererToMainProcess{}
		if err = json.Unmarshal(payload.Cargo, msg); err != nil {
			return
		}
		cargo = msg
	default:
		err = fmt.Errorf("no case found for payload id %d", payload.ID)
	}
	return
}
