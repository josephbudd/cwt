// +build js, wasm

package copytestpanel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
	"github.com/josephbudd/cwt/rendererprocess/api/event"
	"github.com/josephbudd/cwt/rendererprocess/api/markup"
)

/*

	Panel name: CopyTestPanel

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

	copyTestStart *markup.Element
	copyTestCopy  *markup.Element
	copyTestCheck *markup.Element

	solution      [][]*record.KeyCode
	userIsCopying bool
	codeIsKeying  bool
	wpm           uint64
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

	// Define the copy input field.
	if controller.copyTestCopy = document.ElementByID("copyTestCopy"); controller.copyTestCopy == nil {
		err = errors.New("unable to find #copyTestCopy")
		return
	}

	// Define the start button.
	if controller.copyTestStart = document.ElementByID("copyTestStart"); controller.copyTestStart == nil {
		err = errors.New("unable to find #copyTestStart")
		return
	}
	// Handle the submit button's onclick event.
	controller.copyTestStart.SetEventHandler(controller.handleStart, "click", false)

	// Define the check button.
	if controller.copyTestCheck = document.ElementByID("copyTestCheck"); controller.copyTestCheck == nil {
		err = errors.New("unable to find #copyTestCheck")
		return
	}
	// Handle the submit button's onclick event.
	controller.copyTestCheck.SetEventHandler(controller.handleCheck, "click", false)

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

// Onclick handlers.

func (controller *panelController) handleStart(e event.Event) (nilReturn interface{}) {
	controller.userIsCopying = true
	controller.presenter.started()
	controller.messenger.getTextToCopy()
	display.BlockButtonsWithMessage("You are still copying.", "Oops!")
	return
}

func (controller *panelController) handleCheck(e event.Event) (nilReturn interface{}) {
	if controller.codeIsKeying {
		display.Error("Can't stop yet. Still keying.")
		return
	}
	display.UnBlockButtons()
	controller.userIsCopying = false
	controller.presenter.checked()
	copy := strings.TrimSpace(controller.copyTestCopy.Value())
	if len(copy) == 0 {
		display.Error("You didn't enter any copy.")
		return
	}
	practiceWords := make([]string, 0, len(copy))
	lines := strings.Split(copy, "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		for _, w := range words {
			if len(w) > 0 {
				practiceWords = append(practiceWords, w)
			}
		}
	}
	controller.messenger.checkCopy(practiceWords, controller.solution, controller.wpm)
	return
}

// Process data received from the main process.

func (controller *panelController) processWPM(wpm uint64) {
	controller.wpm = wpm
}

func (controller *panelController) processTextToCopy(solution [][]*record.KeyCode) {
	controller.solution = solution
	controller.presenter.ready1()
	display.Inform(
		fmt.Sprintf("The CW will begin %d seconds after you click close. Enter your copy into the red square.", delaySeconds),
		"Copy Practice",
		func() {
			controller.presenter.ready2()
			controller.codeIsKeying = true
			controller.messenger.key(controller.solution, controller.wpm, delaySeconds)
		},
	)
}

func (controller *panelController) processKeyFinished() {
	controller.codeIsKeying = false
	controller.presenter.keyingFinished()
}

func (controller *panelController) processKeyStopped() {
	controller.codeIsKeying = false
	controller.presenter.keyingStopped()
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
