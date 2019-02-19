package calls

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/mainprocess/services/copyservice"
)

// newKeyCall is the constructor for the Key Call.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newKeyCall() *calling.MainProcess {
	return calling.NewMainProcess(
		callids.KeyCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveKey(params, call)
		},
	)
}

// mainProcessReceiveKey is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererKeyCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Build the morse code text.
// 3. Key the morse code. Call back any errors.
func mainProcessReceiveKey(params []byte, callBackToRenderer func(params []byte)) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessKeyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessKey: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererKeyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 2. Build the morse code text.
	ditdahs := make([]string, 0, len(rxparams.Solution))
	for _, line := range rxparams.Solution {
		ditdahWord := make([]string, 0, len(line))
		for _, r := range line {
			ditdahWord = append(ditdahWord, r.DitDah)
		}
		ditdahs = append(ditdahs, strings.Join(ditdahWord, " "))
	}
	// 3. Key the morse code.
	if err := copyservice.Key(ditdahs, rxparams.WPM, rxparams.Pause); err != nil {
		message := fmt.Sprintf("mainProcessKey:  ditdah.Key(rxparams.Ditdah, rxparams.WPM, rxparams.Delay): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererKeyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// no error so call back to the renderer.
	log.Println("mainProcessReceiveKey no errors")
	txparams := &types.MainProcessToRendererKeyCallParams{
		State: rxparams.State,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
