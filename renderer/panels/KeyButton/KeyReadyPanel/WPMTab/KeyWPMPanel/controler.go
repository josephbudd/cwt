package keywpmpanel

import (
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: KeyWPMPanel

*/

// Controler is a HelloPanel Controler.
type Controler struct {
	panelGroup *PanelGroup
	presenter  *Presenter
	caller     *Caller
	quitCh     chan struct{}    // send an empty struct to start the quit process.
	tools      *viewtools.Tools // see /renderer/viewtools
	notJS      *notjs.NotJS

	/* NOTE TO DEVELOPER. Step 1 of 4.

	// Declare your Controler members.

	*/

	record *types.WPMRecord
	keyWPM js.Value
}

// defineControlsSetHandlers defines controler members and sets their handlers.
// Returns the error.
func (panelControler *Controler) defineControlsSetHandlers() (err error) {

	defer func() {
		if err != nil {
			errors.WithMessage(err, "(panelControler *Controler) defineControlsSetHandlers()")
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 4.

	// Define the Controler members by their html elements.
	// Set their handlers.

	*/

	notJS := panelControler.notJS
	tools := panelControler.tools
	null := js.Null()

	// Define the wpm input field.
	if panelControler.keyWPM = notJS.GetElementByID("keyWPM"); panelControler.keyWPM == null {
		err = errors.New("unable to find #keyWPM")
		return
	}
	cb := tools.RegisterEventCallBack(panelControler.handleOnChange, true, true, true)
	notJS.SetOnChange(panelControler.keyWPM, cb)

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// Handlers and other functions.

*/

func (panelControler *Controler) handleOnChange(event js.Value) interface{} {
	panelControler.record.WPM = panelControler.notJS.GetValueUint64(panelControler.keyWPM)
	panelControler.caller.updateKeyWPM(panelControler.record)
	return nil
}

func (panelControler *Controler) processWPM(record *types.WPMRecord) {
	panelControler.record = record
	panelControler.presenter.displayWPM(record.WPM)
}

// initialCalls runs the first code that the controler needs to run.
func (panelControler *Controler) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.

	*/

}
