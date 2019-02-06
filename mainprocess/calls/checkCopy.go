package calls

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/josephbudd/cwt/domain/data/callids"
	"github.com/josephbudd/cwt/domain/implementations/calling"
	"github.com/josephbudd/cwt/domain/interfaces/storer"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/mainprocess/services/copyService"
)

// newCheckCopyCall is the constructor for the CheckCopy Call.
// Param rendererSendPayload: is a kickasm generated renderer func that sends data to the main process.
func newCheckCopyCall(keyCodeStore storer.KeyCodeStorer) *calling.MainProcess {
	return calling.NewMainProcess(
		callids.CheckCopyCallID,
		func(params []byte, call func([]byte)) {
			mainProcessReceiveCheckCopy(params, call, keyCodeStore)
		},
	)
}

// mainProcessReceiveCheckCopy is a main process func.
// This is how the main process receives a call from the renderer.
// Param params is a []byte of a MainProcessToRendererCheckCopyCallParams
// Param callBackToRenderer is a func that calls back to the renderer.
// The func is simple:
// 1. Unmarshall the params. Call back any errors.
// 2. Convert the copy string to key code records.
// 3. Check the copy.
// 4. Call the renderer back with the results.
func mainProcessReceiveCheckCopy(params []byte, callBackToRenderer func(params []byte), keyCodeStore storer.KeyCodeStorer) {
	// 1. Unmarshall the params.
	rxparams := &types.RendererToMainProcessCheckCopyCallParams{}
	if err := json.Unmarshal(params, rxparams); err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessCheckCopy: json.Unmarshal(params, rxparams): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererCheckCopyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	//log.Printf("len(rxparams.Solution) is %d", len(rxparams.Solution))
	// 2. Convert the copy string to key code records.
	// convert copy [][]string to [][]*types.KeyCodeRecord
	rr, err := keyCodeStore.GetKeyCodes()
	if err != nil {
		// Calling back the error.
		message := fmt.Sprintf("mainProcessCheckCopy: keyCodeStore.GetKeyCodes(): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererCheckCopyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	copy := make([][]*types.KeyCodeRecord, 0, len(rxparams.Solution))
	for _, rxCopyLine := range rxparams.Copy {
		copyLine := make([]*types.KeyCodeRecord, 0, len(rxCopyLine))
		for _, rxCopyChar := range rxCopyLine {
			uc := strings.ToUpper(string(rxCopyChar))
			found := false
			for _, r := range rr {
				if r.Character == uc {
					copyLine = append(copyLine, r)
					found = true
					break
				}
			}
			if !found {
				copyLine = append(copyLine, nil)
			}
		}
		copy = append(copy, copyLine)
	}
	// 3. Check the copy.
	nCorrect, nIncorrect, nKeyed, misMatches, err := copyService.Check(copy, rxparams.Solution, keyCodeStore, rxparams.WPM, rxparams.StoreResults)
	if err != nil {
		message := fmt.Sprintf("mainProcessCheckCopy: copyService.Check(copy, rxparams.Solution, keyCodeStore, rxparams.WPM, rxparams.StoreResults): error is %s\n", err.Error())
		log.Println(message)
		txparams := &types.MainProcessToRendererCheckCopyCallParams{
			Error:        true,
			ErrorMessage: message,
		}
		txparamsbb, _ := json.Marshal(txparams)
		callBackToRenderer(txparamsbb)
		return
	}
	// 4. Call the renderer back with the results.
	txparams := &types.MainProcessToRendererCheckCopyCallParams{
		IncorrectCount: nIncorrect,
		CorrectCount:   nCorrect,
		KeyedCount:     nKeyed,
		MisMatches:     misMatches,
		State:          rxparams.State,
	}
	txparamsbb, _ := json.Marshal(txparams)
	callBackToRenderer(txparamsbb)
}
