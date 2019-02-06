package NumbersPanel

import (
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: NumbersPanel

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

	records map[uint64]*types.KeyCodeRecord
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
	// example:

	// import "syscall/js"

	notJS := panelControler.notJS
	tools := panelPresenter.tools
	null := js.Null()

	// Define the customer name input field.
	if panelControler.customerName = notJS.GetElementByID("customerName"); panelControler.customerName == null {
		err = errors.New("unable to find #customerName")
		return
	}

	// Define the submit button and set it's handler.
	if panelControler.addCustomerSubmit = notJS.GetElementByID("addCustomerSubmit"); panelControler.addCustomerSubmit == null {
		err = errors.New("unable to find #addCustomerSubmit")
		return
	}
	cb := notJS.RegisterCallBack(panelControler.handleSubmit)
	notJS.SetOnClick(panelControler.addCustomerSubmit, cb)

	*/

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// Handlers and other functions.

*/

func (panelControler *Controler) setup(records []*types.KeyCodeRecord) {
	notJS := panelControler.notJS
	// fill the table first creating the input elements.
	panelControler.presenter.fillTable(records)
	// rebuild panelControler.records
	panelControler.records = make(map[uint64]*types.KeyCodeRecord)
	cb := notJS.RegisterEventCallBack(false, false, false, panelControler.handleChecked)
	for _, r := range records {
		panelControler.records[r.ID] = r
		// set the checkbox on change handler
		checkboxID := panelControler.presenter.recordIDToCheckBoxID(r.ID)
		checkbox := notJS.GetElementByID(checkboxID)
		notJS.SetOnChange(checkbox, cb)
	}
}

func (panelControler *Controler) handleChecked(event js.Value) {
	notJS := panelControler.notJS
	target := notJS.GetEventTarget(event)
	notJS.Blur(target)
	checked := notJS.GetChecked(target)
	id := notJS.GetAttributeUint64(target, recordIDAttribute)
	record := panelControler.records[id]
	record.Selected = checked
	panelControler.caller.updateKeyCode(record)
}

// initialCalls runs the first code that the controler needs to run.
func (panelControler *Controler) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.
	// example:

	panelControler.customerSelectWidget.start()

	*/

}
