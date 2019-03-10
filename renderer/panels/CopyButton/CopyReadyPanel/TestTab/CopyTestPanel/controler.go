package copytestpanel

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: CopyTestPanel

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

	copyTestStart js.Value
	copyTestCopy  js.Value
	copyTestCheck js.Value

	solution      [][]*types.KeyCodeRecord
	userIsCopying bool
	codeIsKeying  bool
	delaySeconds  uint64
	wpm           uint64
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

	// Define the copy input field.
	if panelControler.copyTestCopy = notJS.GetElementByID("copyTestCopy"); panelControler.copyTestCopy == null {
		err = errors.New("unable to find #copyTestCopy")
		return
	}

	// Define the start button and set it's handler.
	if panelControler.copyTestStart = notJS.GetElementByID("copyTestStart"); panelControler.copyTestStart == null {
		err = errors.New("unable to find #copyTestStart")
		return
	}
	cb := tools.RegisterEventCallBack(panelControler.handleStart, true, true, true)
	notJS.SetOnClick(panelControler.copyTestStart, cb)

	// Define the check button and set it's handler.
	if panelControler.copyTestCheck = notJS.GetElementByID("copyTestCheck"); panelControler.copyTestCheck == null {
		err = errors.New("unable to find #copyTestCheck")
		return
	}
	cb = tools.RegisterEventCallBack(panelControler.handleCheck, true, true, true)
	notJS.SetOnClick(panelControler.copyTestCheck, cb)

	panelControler.delaySeconds = 5

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// Handlers and other functions.

*/

func (panelControler *Controler) processWPM(wpm uint64) {
	panelControler.wpm = wpm
}

func (panelControler *Controler) handleStart(event js.Value) interface{} {
	panelControler.userIsCopying = true
	panelControler.presenter.started()
	panelControler.caller.getTextToCopy()
	return nil
}

func (panelControler *Controler) handleCheck(event js.Value) interface{} {
	if panelControler.codeIsKeying {
		panelControler.tools.Error("Can't stop yet. Still keying.")
		return nil
	}
	panelControler.userIsCopying = false
	panelControler.presenter.checked()
	copy := strings.TrimSpace(panelControler.notJS.GetValue(panelControler.copyTestCopy))
	if len(copy) == 0 {
		panelControler.tools.Error("You didn't enter any copy.")
		return nil
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
	panelControler.caller.checkCopy(practiceWords, panelControler.solution, panelControler.wpm)
	return nil
}

func (panelControler *Controler) processTextToCopy(solution [][]*types.KeyCodeRecord) {
	panelControler.solution = solution
	panelControler.presenter.ready1()
	panelControler.tools.GoModal(
		fmt.Sprintf("The CW will begin %d seconds after you click close. Enter your copy into the red square.", panelControler.delaySeconds),
		"Copy Test",
		func() {
			panelControler.presenter.ready2()
			panelControler.codeIsKeying = true
			panelControler.caller.key(panelControler.solution, panelControler.wpm, panelControler.delaySeconds)
		},
	)
}

func (panelControler *Controler) processKeyFinished() {
	panelControler.codeIsKeying = false
	panelControler.presenter.keyingFinished()
}

// initialCalls runs the first code that the controler needs to run.
func (panelControler *Controler) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.

	*/

}
