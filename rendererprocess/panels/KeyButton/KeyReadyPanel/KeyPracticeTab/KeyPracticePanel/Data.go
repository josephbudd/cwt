// +build js, wasm

package keypracticepanel

import (
	"github.com/josephbudd/cwt/rendererprocess/api/dom"
	"github.com/josephbudd/cwt/rendererprocess/framework/lpc"
)

/*

	Panel name: KeyPracticePanel

*/

var (
	// quitCh will close the application
	quitCh chan struct{}

	// eojCh will signal go routines to stop and return because the application is ending.
	eojCh chan struct{}

	// receiveCh receives messages from the main process.
	receiveCh lpc.Receiving

	// sendCh sends messages to the main process.
	sendCh lpc.Sending

	// The document object module.
	document *dom.DOM

	// Practice or test.
	state uint64
)
