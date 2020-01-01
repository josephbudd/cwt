// +build js, wasm

package keytestpanel

import (
	"errors"
	"fmt"

	"github.com/josephbudd/cwt/rendererprocess/api/markup"
	"github.com/josephbudd/cwt/rendererprocess/widgets"
)

/*

	Panel name: KeyTestPanel

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

	var keyTestH, keyTestKey, keyTestCopy, keyTestStart, keyTestCheck *markup.Element

	// Define the heading.
	if keyTestH = document.ElementByID("keyTestH"); keyTestH == nil {
		err = errors.New("unable to find #keyTestH")
		return
	}

	// Define the key area where the user mouses over with the keyer-mouse when keying.
	if keyTestKey = document.ElementByID("keyTestKey"); keyTestKey == nil {
		err = errors.New("unable to find #keyTestKey")
		return
	}

	// Define the copy area where the user can read the copy to key.
	if keyTestCopy = document.ElementByID("keyTestCopy"); keyTestCopy == nil {
		err = errors.New("unable to find #keyTestCopy")
		return
	}

	// Define the start button for the keyWidget.
	if keyTestStart = document.ElementByID("keyTestStart"); keyTestStart == nil {
		err = errors.New("unable to find #keyTestStart")
		return
	}

	// Define the check button for the keyWidget.
	if keyTestCheck = document.ElementByID("keyTestCheck"); keyTestCheck == nil {
		err = errors.New("unable to find #keyTestCheck")
		return
	}

	// Define the keyWidget.
	controller.keyWidget = widgets.NewKeyWidget(document, keyTestH, keyTestStart, keyTestCheck, nil, keyTestKey, keyTestCopy, controller.messenger, nil)

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
