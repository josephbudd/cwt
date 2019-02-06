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

// newGetCopyWPMCall is the constructor for the GetCopyWPM Call.
// It should only receive the repos that are needed. In this case the keyCode repo.
// Param keyCodeStorer storer.CopyCodeStorer is the keyCode repo needed to get a keyCode record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newGetCopyWPMCall(wPMStorer storer.WPMStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.GetCopyWPMCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveGetCopyWPM(params, call, wPMStorer)
		},
	)
}

// mainProcessReceiveGetCopyWPM is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererGetCopyWPMCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// Param wPMStorer storer.WPMStorer is the wpm repo.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Get the keyCodes from the repo. Call back any errors or not found.
// 3. Call the renderer back with the keyCode records.
func mainProcessReceiveGetCopyWPM(params []byte, callBackToRenderer func(params []byte), wPMStorer storer.WPMStorer) {
	log.Println("mainProcessReceiveGetCopyWPM")
	// 1. ignore params.
	// 2. Get the keyCode from the repo.
	r, err := wPMStorer.GetCopyWPM()
	if err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessGetCopyWPM: keyCodeStorer.GetCopyWPM(): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererGetCopyWPMCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	if r == nil {
		// Calling back the not found error.
		// This will only happen in development. It means that the data store is not getting initialized properly.
		message := "mainProcessGetCopyWPM: keyCodeStorer.GetCopyWPM(): error is Not Found.\n"
		log.Println(message)
		txparams := &types.MainProcessToRendererGetCopyWPMCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 3. Call the renderer back with the keyCode records.
	txparams := &types.MainProcessToRendererGetCopyWPMCallParams{
		Record: r,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
