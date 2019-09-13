package specialpanel

import (
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: SpecialPanel

*/

// panelController controls user input.
type panelController struct {
	group     *panelGroup
	presenter *panelPresenter
	caller    *panelCaller
	eventCh   chan viewtools.Event

	/* NOTE TO DEVELOPER. Step 1 of 5.

	// Declare your panelController members.

	*/

	records map[uint64]*record.KeyCode
}

// defineControlsReceiveEvents defines controller members and starts receiving their events.
// Returns the error.
func (controller *panelController) defineControlsReceiveEvents() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(controller *panelController) defineControlsReceiveEvents()")
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 5.

	// Define the controller members by their html elements.
	// Receive their events.

	*/

	return
}

/* NOTE TO DEVELOPER. Step 3 of 5.

// Handlers and other functions.

*/

func (controller *panelController) setup(rr []*record.KeyCode) {
	// fill the table first creating the input elements.
	controller.presenter.fillTable(rr)
	// rebuild controller.records
	controller.records = make(map[uint64]*record.KeyCode)
	cb := tools.RegisterEventCallBack(controller.handleChecked, true, true, true)
	for _, r := range rr {
		controller.records[r.ID] = r
		// set the checkbox on change handler
		checkboxID := controller.presenter.recordIDToCheckBoxID(r.ID)
		checkbox := notJS.GetElementByID(checkboxID)
		notJS.SetOnChange(checkbox, cb)
	}
}

func (controller *panelController) handleChecked(event js.Value) (nilReturn interface{}) {
	target := notJS.GetEventTarget(event)
	notJS.Blur(target)
	checked := notJS.GetChecked(target)
	id := notJS.GetAttributeUint64(target, recordIDAttribute)
	record := controller.records[id]
	record.Selected = checked
	controller.caller.updateKeyCode(record)
	return
}

// dispatchEvents dispatches events from the controls.
// It stops when it receives on the eoj channel.
func (controller *panelController) dispatchEvents() {
	go func() {
		var event viewtools.Event
		for {
			select {
			case <-eojCh:
				return
			case event = <-controller.eventCh:
				// An event that this controller is receiving from one of its members.
				switch event.Target {

				/* NOTE TO DEVELOPER. Step 4 of 5.

				// 4.1.a: Add a case for each controller member
				//          that you are receiving events for.
				// 4.1.b: In that case statement, pass the event to your event handler.

				*/

				}
			}
		}
	}()

	return
}

// initialCalls runs the first code that the controller needs to run.
func (controller *panelController) initialCalls() {

	/* NOTE TO DEVELOPER. Step 5 of 5.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.

	*/
}

// receiveEvent gets this controller listening for element's event.
// Param elements if the controller's element.
// Param event is the event ex: "onclick".
// Param preventDefault indicates if the default behavior of the event must be prevented.
// Param stopPropagation indicates if the event's propogation must be stopped.
// Param stopImmediatePropagation indicates if the event's immediate propogation must be stopped.
func (controller *panelController) receiveEvent(element js.Value, event string, preventDefault, stopPropagation, stopImmediatePropagation bool) {
	tools.SendEvent(controller.eventCh, element, event, preventDefault, stopPropagation, stopImmediatePropagation)
}
