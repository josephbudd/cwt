package copypracticepanel

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/store/record"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: CopyPracticePanel

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

	copyPracticeStart js.Value
	copyPracticeCopy  js.Value
	copyPracticeCheck js.Value

	solution      [][]*record.KeyCode
	userIsCopying bool
	codeIsKeying  bool
	delaySeconds  uint64
	wpm           uint64

	lockMessageTitle string
	lockMessage      string
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

	// Define the copy input field.
	if controller.copyPracticeCopy = notJS.GetElementByID("copyPracticeCopy"); controller.copyPracticeCopy == null {
		err = errors.New("unable to find #copyPracticeCopy")
		return
	}

	// Define the start button and set it's handler.
	if controller.copyPracticeStart = notJS.GetElementByID("copyPracticeStart"); controller.copyPracticeStart == null {
		err = errors.New("unable to find #copyPracticeStart")
		return
	}
	// Receive the submit button's onclick event.
	controller.receiveEvent(controller.copyPracticeStart, "onclick", true, true, true)

	// Define the check button and set it's handler.
	if controller.copyPracticeCheck = notJS.GetElementByID("copyPracticeCheck"); controller.copyPracticeCheck == null {
		err = errors.New("unable to find #copyPracticeCheck")
		return
	}
	// Receive the submit button's onclick event.
	controller.receiveEvent(controller.copyPracticeCheck, "onclick", true, true, true)

	controller.delaySeconds = 5
	controller.lockMessageTitle = "Oops!"
	controller.lockMessage = "You are still copying."

	return
}

/* NOTE TO DEVELOPER. Step 3 of 5.

// Handlers and other functions.

*/

func (controller *panelController) processWPM(wpm uint64) {
	controller.wpm = wpm
}

func (controller *panelController) handleStart(event js.Value) {
	controller.userIsCopying = true
	controller.presenter.started()
	controller.caller.getTextToCopy()
	tools.LockButtonsWithMessage(controller.lockMessage, controller.lockMessageTitle)
}

func (controller *panelController) handleCheck(event js.Value) {
	if controller.codeIsKeying {
		tools.Error("Can't stop yet. Still keying.")
		return
	}
	tools.UnLockButtons()
	controller.userIsCopying = false
	controller.presenter.checked()
	userCopy := strings.TrimSpace(notJS.GetValue(controller.copyPracticeCopy))
	if len(userCopy) == 0 {
		tools.Error("You didn't enter any copy.")
		return
	}
	practiceWords := make([]string, 0, len(userCopy))
	lines := strings.Split(userCopy, "\n")
	for _, line := range lines {
		words := strings.Split(line, " ")
		for _, w := range words {
			if len(w) > 0 {
				practiceWords = append(practiceWords, w)
			}
		}
	}
	controller.caller.checkCopy(practiceWords, controller.solution, controller.wpm)
}

func (controller *panelController) processTextToCopy(solution [][]*record.KeyCode) {
	controller.solution = solution
	controller.presenter.ready1()
	tools.GoModal(
		fmt.Sprintf("The CW will begin %d seconds after you click close. Enter your copy into the red square.", controller.delaySeconds),
		"Copy Practice",
		func() {
			controller.presenter.ready2()
			controller.codeIsKeying = true
			controller.caller.key(controller.solution, controller.wpm, controller.delaySeconds)
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

				case controller.copyPracticeStart:
					controller.handleStart(event.Event)
				case controller.copyPracticeCheck:
					controller.handleCheck(event.Event)
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
	// example:

	controller.customerSelectWidget.start()

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
