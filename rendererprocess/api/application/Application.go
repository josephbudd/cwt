// +build js, wasm

package application

import (
	"github.com/josephbudd/cwt/rendererprocess/api/display"
	"github.com/josephbudd/cwt/rendererprocess/api/event"
	"github.com/josephbudd/cwt/rendererprocess/framework/callback"
)

// NewGracefullyCloseHandler makes an event handler which gracefully closes the application for you.
// Use it in your panel controllers to handle your own application closing buttons.
// Param quitCh is the panel controller's quit channel.
func NewGracefullyCloseHandler(quitCh chan struct{}) (handler func(e event.Event) (nilReturn interface{})) {
	handler = func(e event.Event) (nilReturn interface{}) {
		callback.RemoveApplicationOnCloseHandler()
		title := "Closing"
		msg := "Closing <q>Linux CW Trainer</q>."
		cb := func() { quitCh <- struct{}{} }
		display.Inform(msg, title, cb)
		return
	}
	return
}
