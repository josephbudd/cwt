package calls

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
)

// newUpdateKeyWPMCall is the constructor for the UpdateKeyWPM Call.
// It should only receive the repos that are needed. In this case the wpm repo.
// Param wpmStorer storer.WPMStorer is the wpm repo needed to get a wpm record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newUpdateKeyWPMCall(wpmStorer storer.WPMStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.UpdateKeyWPMCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveUpdateKeyWPM(params, call, wpmStorer)
		},
	)
}

// mainProcessReceiveUpdateKeyWPM is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererUpdateKeyWPMCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// Param wpmStorer is the wpm repo.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Get the wpm from the repo. Call back any errors or not found.
// 3. Call the renderer back with the wpm record.
func mainProcessReceiveUpdateKeyWPM(params []byte, callBackToRenderer func(params []byte), wpmStorer storer.WPMStorer) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessUpdateKeyWPMCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		log.Println("mainProcessReceiveUpdateKeyWPM error is ", err.Error())
		message := fmt.Sprintf("mainProcessUpdateKeyWPM: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererUpdateKeyWPMCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2. Update the wpm.
	if err := wpmStorer.UpdateWPM(rxparams.Record); err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessUpdateKeyWPM: wpmStorer.UpdateKeyWPM(rxparams.ID): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererUpdateKeyWPMCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 3. Call back with the updated record.
	txparams := &types.MainProcessToRendererUpdateKeyWPMCallParams{
		Record: rxparams.Record,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
