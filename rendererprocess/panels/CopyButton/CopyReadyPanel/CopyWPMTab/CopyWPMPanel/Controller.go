// +build js, wasm

package copywpmpanel

import (
	"errors"
	"fmt"

	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/event"
	"github.com/josephbudd/cwt/rendererprocess/api/markup"
)

/*

	Panel name: CopyWPMPanel

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

	recordWPM *record.WPM
	copyWPM   *markup.Element
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

	// Define the wpm input field.
	if controller.copyWPM = document.ElementByID("copyWPM"); controller.copyWPM == nil {
		err = errors.New("unable to find #copyWPM")
		return
	}
	// Handle the wpm's onchange event.
	controller.copyWPM.SetEventHandler(controller.handleOnChange, "click", false)

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

// Event handlers.

func (controller *panelController) handleOnChange(e event.Event) (nilReturn interface{}) {
	controller.recordWPM.WPM, _ = controller.copyWPM.ValueUint64()
	controller.messenger.updateCopyWPM(controller.recordWPM)
	return
}

// Process information from the main process.

func (controller *panelController) processWPM(r *record.WPM) {
	controller.recordWPM = r
	controller.presenter.displayWPM(r.WPM)
}

// initialCalls runs the first code that the controller needs to run.
func (controller *panelController) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.

	// example:

	controller.customerSelectWidget.start()

	*/
}
