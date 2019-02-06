package calls

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/mainprocess/services/keyService"
)

// newCheckKeyCall is the constructor for the CheckKey Call.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
// Param keyCodeStorer is the key code storer.
func newCheckKeyCall(wPMStorer storer.WPMStorer, keyCodeStorer storer.KeyCodeStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.CheckKeyCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveCheckKey(params, call, wPMStorer, keyCodeStorer)
		},
	)
}

// mainProcessReceiveCheckKey is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererCheckKeyCallParams.
// Param callBackToRenderer is a func that calls back to the renderer.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Copy: Convert the miliseconds to key code records.
// 2.1 Check the copy against the solution.
func mainProcessReceiveCheckKey(params []byte, callBackToRenderer func(params []byte), wPMStorer storer.WPMStorer, keyCodeStorer storer.KeyCodeStorer) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessCheckKeyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessReceiveCheckKey: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererCheckKeyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2. Copy: Convert the miliseconds to key code records.
	copiedWords, err := keyService.Copy(rxparams.MilliSeconds, rxparams.WPM, wPMStorer, keyCodeStorer)
	if err != nil {
		message := fmt.Sprintf("mainProcessReceiveCheckKey: keyService.Copy(rxparams.MilliSeconds, rxparams.WPM, wPMStorer, keyCodeStorer): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererCheckKeyCallParams{
			Error:        true,
			ErrorMessage: message,
			State:        rxparams.State,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2.1 Check the copy against the solution.
	nCorrect, nIncorrect, nKeyed, misMatches, err := keyService.Check(copiedWords, rxparams.Solution, keyCodeStorer, rxparams.WPM, rxparams.StoreResults)
	if err != nil {
		message := fmt.Sprintf("mainProcessReceiveCheckKey: keyService.Check(copiedWords, rxparams.Solution, keyCodeStorer, rxparams.WPM, rxparams.StoreResults): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererCheckKeyCallParams{
			Error:        true,
			ErrorMessage: message,
			State:        rxparams.State,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// Call back the results.
	txparams := &types.MainProcessToRendererCheckKeyCallParams{
		CorrectCount:   nCorrect,
		IncorrectCount: nIncorrect,
		KeyedCount:     nKeyed,
		MisMatches:     misMatches,
		State:          rxparams.State,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
	return
}
