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

// newUpdateCopyWPMCall is the constructor for the UpdateCopyWPM Call.
// It should only receive the repos that are needed. In this case the wpm repo.
// Param wpmStorer storer.WPMStorer is the wpm repo needed to get a wpm record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newUpdateCopyWPMCall(wpmStorer storer.WPMStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.UpdateCopyWPMCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveUpdateCopyWPM(params, call, wpmStorer)
		},
	)
}

// mainProcessReceiveUpdateCopyWPM is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererUpdateCopyWPMCallCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// Param wpmStorer is the wpm repo.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Get the wpm from the repo. Call back any errors or not found.
// 3. Call the renderer back with the wpm record.
func mainProcessReceiveUpdateCopyWPM(params []byte, callBackToRenderer func(params []byte), wpmStorer storer.WPMStorer) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessUpdateCopyWPMCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		log.Println("mainProcessReceiveUpdateCopyWPM error is ", err.Error())
		message := fmt.Sprintf("mainProcessUpdateCopyWPM: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererUpdateCopyWPMCallParams{
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
		message := fmt.Sprintf("mainProcessUpdateCopyWPM: wpmStorer.UpdateCopyWPM(rxparams.ID): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererUpdateCopyWPMCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 3. Call back with the updated record.
	txparams := &types.MainProcessToRendererUpdateCopyWPMCallParams{
		Record: rxparams.Record,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
