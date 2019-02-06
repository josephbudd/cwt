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

// newGetKeyCodesCall is the constructor for the GetKeyCodes Call.
// It should only receive the repos that are needed. In this case the keyCode repo.
// Param keyCodeStorer storer.KeyCodeStorer is the keyCode repo needed to get a keyCode record.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newGetKeyCodesCall(keyCodeStorer storer.KeyCodeStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.GetKeyCodesCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveGetKeyCodes(params, call, keyCodeStorer)
		},
	)
}

// mainProcessReceiveGetKeyCodes is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererGetKeyCodesCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// Param keyCodeStorer is the keyCode repo.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Get the keyCodes from the repo. Call back any errors or not found.
// 3. Call the renderer back with the keyCode records.
func mainProcessReceiveGetKeyCodes(params []byte, callBackToRenderer func(params []byte), keyCodeStorer storer.KeyCodeStorer) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessGetKeyCodesCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		log.Println("mainProcessReceiveGetKeyCodes error is ", err.Error())
		message := fmt.Sprintf("mainProcessGetKeyCodes: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererGetKeyCodesCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2. Get the keyCode from the repo.
	keyCodes, err := keyCodeStorer.GetKeyCodes()
	if err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessGetKeyCodes: keyCodeStorer.GetKeyCodes(rxparams.ID): error is %s\n", err.Error())
		txparams := &types.MainProcessToRendererGetKeyCodesCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 3. Call the renderer back with the keyCode records.
	txparams := &types.MainProcessToRendererGetKeyCodesCallParams{
		Records: keyCodes,
		State:   rxparams.State,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
