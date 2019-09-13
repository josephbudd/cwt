package keytestpanel

import (
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/renderer/viewtools"
	"github.com/josephbudd/cwt/renderer/widgets"
)

/*

	Panel name: KeyTestPanel

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

	keyWidget *widgets.KeyWidget
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

	var keyTestH, keyTestKey, keyTestCopy, keyTestStart, keyTestCheck js.Value

	// Define the heading.
	if keyTestH = notJS.GetElementByID("keyTestH"); keyTestH == null {
		err = errors.New("unable to find #keyTestH")
		return
	}

	// Define the key area where the user mouses over with the keyer-mouse when keying.
	if keyTestKey = notJS.GetElementByID("keyTestKey"); keyTestKey == null {
		err = errors.New("unable to find #keyTestKey")
		return
	}

	// Define the copy area where the user can read the copy to key.
	if keyTestCopy = notJS.GetElementByID("keyTestCopy"); keyTestCopy == null {
		err = errors.New("unable to find #keyTestCopy")
		return
	}

	// Define the start button for the keyWidget.
	if keyTestStart = notJS.GetElementByID("keyTestStart"); keyTestStart == null {
		err = errors.New("unable to find #keyTestStart")
		return
	}

	// Define the check button for the keyWidget.
	if keyTestCheck = notJS.GetElementByID("keyTestCheck"); keyTestCheck == null {
		err = errors.New("unable to find #keyTestCheck")
		return
	}

	// Define the keyWidget.
	controller.keyWidget = widgets.NewKeyWidget(keyTestH, keyTestStart, keyTestCheck, null, keyTestKey, keyTestCopy, controller.caller, nil, tools, notJS)

	return
}

/* NOTE TO DEVELOPER. Step 3 of 5.

// Handlers and other functions.

*/

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
