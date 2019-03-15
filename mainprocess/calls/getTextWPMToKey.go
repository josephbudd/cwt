package calls

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/mainprocess/services/keyservice"
)

// newGetTextWPMToKeyCall is the constructor for the GetTextWPMToKey Call.
// It should only receive the repos that are needed. In this case the customer repo.
// Param keyCodeStorer storer.KeyCodeStorer is the keycode store.
// Param wPMStorer storer.WPMStorer is the wpm storer.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newGetTextWPMToKeyCall(keyCodeStorer storer.KeyCodeStorer, wPMStorer storer.WPMStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.GetTextWPMToKeyCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveGetTextWPMToKey(params, call, keyCodeStorer, wPMStorer)
		},
	)
}

// mainProcessReceiveGetTextWPMToKey is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererGetTextWPMToKeyCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// Param keyCodeStorer is the key code storer.
// Param wPMStorer is the wpm storer.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Get the wpm for the user to key. Call back any errors or not found.
// 3. Get the text for the user to key from the repo and the help. Call back any errors or not found.
// 4. Get the help for each word in the text.
// 5. Call the renderer back with the text, wpm.
func mainProcessReceiveGetTextWPMToKey(params []byte, callBackToRenderer func(params []byte), keyCodeStorer storer.KeyCodeStorer, wPMStorer storer.WPMStorer) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessGetTextWPMToKeyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		log.Println("mainProcessReceiveGetTextWPMToKey error is ", err.Error())
		message := fmt.Sprintf("mainProcessGetTextWPMToKey: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererGetTextWPMToKeyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2 Get the wpm to key.
	r, err := wPMStorer.GetKeyWPM()
	if err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessGetTextWPMToKey: keyCodeStorer.GetKeyWPM(): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererGetTextWPMToKeyCallParams{
			State:        rxparams.State,
			Practice:     rxparams.Practice,
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 3. Get the text for the user to key from the repo and the help.
	var text [][]*types.KeyCodeRecord
	var help [][]types.HowTo
	if rxparams.Practice {
		// practicing
		text, help, err = keyservice.GetPracticeKeyCodes(keyCodeStorer, r.WPM)
		if err != nil {
			// Calling back the error.
			message := fmt.Sprintf("mainProcessGetTextWPMToKey: keyservice.GetPracticeKeyCodes(keyCodeStorer, r.WPM): error is %s\n", err.Error())
			txparams := &types.MainProcessToRendererGetTextWPMToKeyCallParams{
				State:        rxparams.State,
				Practice:     rxparams.Practice,
				Error:        true,
				ErrorMessage: message,
			}
			txparamsbb, _ := json.Marshal(txparams)
			callBackToRenderer(txparamsbb)
			return
		}
	} else {
		// testing not practicing
		text, err = keyservice.GetTestKeyCodes(keyCodeStorer)
		if err != nil {
			// Calling back the error.
			message := fmt.Sprintf("mainProcessGetTextWPMToKey: keyservice.GetTestKeyCodes(keyCodeStorer): error is %s\n", err.Error())
			txparams := &types.MainProcessToRendererGetTextWPMToKeyCallParams{
				State:        rxparams.State,
				Practice:     rxparams.Practice,
				Error:        true,
				ErrorMessage: message,
			}
			txparamsbb, _ := json.Marshal(txparams)
			callBackToRenderer(txparamsbb)
			return
		}
	}
	// . Call the renderer back with the text, wpm.
	txparams := &types.MainProcessToRendererGetTextWPMToKeyCallParams{
		Solution: text,
		Help:     help,
		WPM:      r.WPM,
		Practice: rxparams.Practice,
		State:    rxparams.State,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
