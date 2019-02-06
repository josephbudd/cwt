package calls

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/mainprocess/services/copyService"
)

// newGetTextToCopyCall is the constructor for the GetTextToCopy Call.
// It should only receive the repos that are needed. In this case the customer repo.
// Param keyCodeStorer storer.KeyCodeStorer is the keycode store.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newGetTextToCopyCall(keyCodeStorer storer.KeyCodeStorer, wPMStorer storer.WPMStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.GetTextToCopyCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveGetTextToCopy(params, call, keyCodeStorer, wPMStorer)
		},
	)
}

// mainProcessReceiveGetTextToCopy is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererGetTextToCopyCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// Param customerStorer is the customer repo.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Get the text and ditdah to copy.
// 3. Get the copy WPM for rendering the ditdahs.
// 4. Call the renderer back with the text.
func mainProcessReceiveGetTextToCopy(params []byte, callBackToRenderer func(params []byte), keyCodeStorer storer.KeyCodeStorer, wPMStorer storer.WPMStorer) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessGetTextToCopyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		log.Println("mainProcessReceiveGetTextToCopy error is ", err.Error())
		message := fmt.Sprintf("mainProcessGetTextToCopy: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererGetTextToCopyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2. Get the text and ditdah to copy.
	solution, err := copyService.GetKeyCodes(keyCodeStorer)
	if err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessGetTextToCopy: keyService.GetTextToCopy(keyCodeStorer): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererGetTextToCopyCallParams{
			Error:        true,
			ErrorMessage: message,
			State:        rxparams.State,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	log.Printf("mainProcessReceiveGetTextToCopy solution is %#v", solution)
	// 3. Get the copy WPM for rendering the ditdahs.
	r, err := wPMStorer.GetCopyWPM()
	if err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessGetTextToCopy: wPMStorer.GetCopyWPM(): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererGetTextToCopyCallParams{
			Error:        true,
			ErrorMessage: message,
			State:        rxparams.State,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 4. Call the renderer back with the text.
	txparams := &types.MainProcessToRendererGetTextToCopyCallParams{
		Solution: solution,
		WPM:      r.WPM,
		State:    rxparams.State,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
