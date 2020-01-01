// +build js, wasm

package keypracticepanel

import (
	"errors"
	"fmt"

	"github.com/josephbudd/cwt/rendererprocess/api/markup"
	"github.com/josephbudd/cwt/rendererprocess/widgets"
)

/*

	Panel name: KeyPracticePanel

*/

// panelController controls user input.
type panelController struct {
	group     *panelGroup
	presenter *panelPresenter
	messenger *panelMessenger

	/* NOTE TO DEVELOPER. Step 1 of 4.

	// Declare your panelController fields.

	// example:

	import "github.com/josephbudd/cwt/rendererprocess/api/markup"

	addCustomerName   *markup.Element
	addCustomerSubmit *markup.Element

	*/

	keyWidget *widgets.KeyWidget
}

// defineControlsHandlers defines the GUI's controllers and their event handlers.
// Returns the error.
func (controller *panelController) defineControlsHandlers() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("(controller *panelController) defineControlsHandlers(): %w", err)
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 4.

	// Define each controller in the GUI by it's html element.
	// Handle each controller's events.

	// example:

	// Define the customer name text input GUI controller.
	if controller.addCustomerName = document.ElementByID("addCustomerName"); controller.addCustomerName == nil {
		err = fmt.Errorf("unable to find #addCustomerName")
		return
	}

	// Define the submit button GUI controller.
	if controller.addCustomerSubmit = document.ElementByID("addCustomerSubmit"); controller.addCustomerSubmit == nil {
		err = fmt.Errorf("unable to find #addCustomerSubmit")
		return
	}
	// Handle the submit button's onclick event.
	controller.addCustomerSubmit.SetEventHandler(controller.handleSubmit, "click", false)

	*/

	var keyPracticeH, keyPracticeKey, keyPracticeCopy, keyPracticeStart, keyPracticeCheck, keyPracticeMetronomeOn *markup.Element

	// Define the heading.
	if keyPracticeH = document.ElementByID("keyPracticeH"); keyPracticeH == nil {
		err = errors.New("unable to find #keyPracticeH")
		return
	}

	// Define the key area where the user mouses over with the keyer-mouse when keying.
	if keyPracticeKey = document.ElementByID("keyPracticeKey"); keyPracticeKey == nil {
		err = errors.New("unable to find #keyPracticeKey")
		return
	}

	// Define the copy area where the user can read the copy to key.
	if keyPracticeCopy = document.ElementByID("keyPracticeCopy"); keyPracticeCopy == nil {
		err = errors.New("unable to find #keyPracticeCopy")
		return
	}

	// Define the start button for the keyWidget.
	if keyPracticeStart = document.ElementByID("keyPracticeStart"); keyPracticeStart == nil {
		err = errors.New("unable to find #keyPracticeStart")
		return
	}

	// Define the check button for the keyWidget.
	if keyPracticeCheck = document.ElementByID("keyPracticeCheck"); keyPracticeCheck == nil {
		err = errors.New("unable to find #keyPracticeCheck")
		return
	}

	// Define the metronome on check box for the keyWidget.
	if keyPracticeMetronomeOn = document.ElementByID("keyPracticeMetronomeOn"); keyPracticeMetronomeOn == nil {
		err = errors.New("unable to find #keyPracticeMetronomeOn")
		return
	}

	// Define the keyWidget.
	controller.keyWidget = widgets.NewKeyWidget(document, keyPracticeH, keyPracticeStart, keyPracticeCheck, keyPracticeMetronomeOn, keyPracticeKey, keyPracticeCopy, controller.messenger, controller.messenger)

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// Handlers and other functions.

// example:

import "github.com/josephbudd/cwt/domain/store/record"
import "github.com/josephbudd/cwt/rendererprocess/api/event"
import "github.com/josephbudd/cwt/rendererprocess/api/display"

func (controller *panelController) handleSubmit(e event.Event) (nilReturn interface{}) {
	// See rendererprocess/api/event/event.go.
	// The event.Event funcs.
	//   e.PreventDefaultBehavior()
	//   e.StopCurrentPhasePropagation()
	//   e.StopAllPhasePropagation()
	//   target := e.JSTarget
	//   event := e.JSEvent
	// You must use the javascript event e.JSEvent, as a js.Value.
	// However, you can use the target as a *markup.Element
	//   target := document.NewElementFromJSValue(e.JSTarget)

	name := strings.TrimSpace(controller.addCustomerName.Value())
	if len(name) == 0 {
		display.Error("Customer Name is required.")
		return
	}
	r := &record.Customer{
		Name: name,
	}
	controller.messenger.AddCustomer(r)
	return
}

*/

// initialCalls runs the first code that the controller needs to run.
func (controller *panelController) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.

	// example:

	controller.customerSelectWidget.start()

	*/
}
