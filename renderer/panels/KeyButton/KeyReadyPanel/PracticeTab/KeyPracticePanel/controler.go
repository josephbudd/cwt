package keypracticepanel

import (
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
	"github.com/josephbudd/cwt/renderer/widgets"
)

/*

	Panel name: KeyPracticePanel

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

	keyWidget *widgets.KeyWidget
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
	null := js.Null()
	var keyPracticeH, keyPracticeKey, keyPracticeCopy, keyPracticeStart, keyPracticeCheck js.Value

	// Define the heading.
	if keyPracticeH = notJS.GetElementByID("keyPracticeH"); keyPracticeH == null {
		err = errors.New("unable to find #keyPracticeH")
		return
	}

	// Define the key area where the user mouses over with the keyer-mouse when keying.
	if keyPracticeKey = notJS.GetElementByID("keyPracticeKey"); keyPracticeKey == null {
		err = errors.New("unable to find #keyPracticeKey")
		return
	}

	// Define the copy area where the user can read the copy to key.
	if keyPracticeCopy = notJS.GetElementByID("keyPracticeCopy"); keyPracticeCopy == null {
		err = errors.New("unable to find #keyPracticeCopy")
		return
	}

	// Define the start button for the keyWidget.
	if keyPracticeStart = notJS.GetElementByID("keyPracticeStart"); keyPracticeStart == null {
		err = errors.New("unable to find #keyPracticeStart")
		return
	}

	// Define the check button for the keyWidget.
	if keyPracticeCheck = notJS.GetElementByID("keyPracticeCheck"); keyPracticeCheck == null {
		err = errors.New("unable to find #keyPracticeCheck")
		return
	}

	// Define the keyWidget.
	panelControler.keyWidget = widgets.NewKeyWidget(keyPracticeH, keyPracticeStart, keyPracticeCheck, keyPracticeKey, keyPracticeCopy, panelControler.caller, panelControler.caller, panelControler.tools, notJS)
	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// Handlers and other functions.

*/

// initialCalls runs the first code that the controler needs to run.
func (panelControler *Controler) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.

	*/

}
