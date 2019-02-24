package calls

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/mainprocess/services/keyservice"
)

// newMetronomeCall is the constructor for the Metronome Call.
// It should only receive the repos that are needed. In this case the wpm repo.
func newMetronomeCall() *calling.MainProcess {
	return calling.NewMainProcess(
		callids.MetronomeCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveMetronome(params, call)
		},
	)
}

// mainProcessReceiveMetronome is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererMetronomeCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Turn the metronome on or off.
// 3. Call the renderer back.
func mainProcessReceiveMetronome(params []byte, callBackToRenderer func(params []byte)) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessMetronomeCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		log.Println("mainProcessReceiveMetronome error is ", err.Error())
		message := fmt.Sprintf("mainProcessMetronome: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererMetronomeCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2. Turn the metronome on or off.
	if rxparams.Run {
		errCh := make(chan error)
		keyservice.StartMetronome(rxparams.WPM, errCh)
		go handleMetronomeError(rxparams, callBackToRenderer, errCh)
	} else {
		keyservice.StopMetronome()
	}
	// 3. Call back.
	txparams := &types.MainProcessToRendererMetronomeCallParams{
		State: rxparams.State,
		Run:   rxparams.Run,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}

func handleMetronomeError(rxparams *types.RendererToMainProcessMetronomeCallParams, callBackToRenderer func(params []byte), errCh chan error) {
	select {
	case err := <-errCh:
		if err != nil {
			txparams := &types.MainProcessToRendererMetronomeCallParams{
				State:        rxparams.State,
				Run:          rxparams.Run,
				Error:        true,
				ErrorMessage: err.Error(),
			}
			txparamsbb, _ := json.Marshal(txparams)
			callBackToRenderer(txparamsbb)
		}
	}
}
