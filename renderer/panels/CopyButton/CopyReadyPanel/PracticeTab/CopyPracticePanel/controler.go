package CopyPracticePanel

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

	Panel name: CopyPracticePanel

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

	copyPracticeStart js.Value
	copyPracticeCopy  js.Value
	copyPracticeCheck js.Value

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
	null := js.Null()

	// Define the copy input field.
	if panelControler.copyPracticeCopy = notJS.GetElementByID("copyPracticeCopy"); panelControler.copyPracticeCopy == null {
		err = errors.New("unable to find #copyPracticeCopy")
		return
	}

	// Define the start button and set it's handler.
	if panelControler.copyPracticeStart = notJS.GetElementByID("copyPracticeStart"); panelControler.copyPracticeStart == null {
		err = errors.New("unable to find #copyPracticeStart")
		return
	}
	cb := notJS.RegisterCallBack(panelControler.handleStart)
	notJS.SetOnClick(panelControler.copyPracticeStart, cb)

	// Define the check button and set it's handler.
	if panelControler.copyPracticeCheck = notJS.GetElementByID("copyPracticeCheck"); panelControler.copyPracticeCheck == null {
		err = errors.New("unable to find #copyPracticeCheck")
		return
	}
	cb = notJS.RegisterCallBack(panelControler.handleCheck)
	notJS.SetOnClick(panelControler.copyPracticeCheck, cb)

	panelControler.delaySeconds = 5

	return
}

/* NOTE TO DEVELOPER. Step 3 of 4.

// Handlers and other functions.

*/

func (panelControler *Controler) processWPM(wpm uint64) {
	panelControler.wpm = wpm
}

func (panelControler *Controler) handleStart([]js.Value) {
	panelControler.userIsCopying = true
	panelControler.presenter.started()
	panelControler.caller.getTextToCopy()
}

func (panelControler *Controler) handleCheck([]js.Value) {
	if panelControler.codeIsKeying {
		panelControler.tools.Error("Can't stop yet. Still keying.")
		return
	}
	panelControler.userIsCopying = false
	panelControler.presenter.checked()
	copy := strings.TrimSpace(panelControler.notJS.GetValue(panelControler.copyPracticeCopy))
	if len(copy) == 0 {
		panelControler.tools.Error("You didn't enter any copy.")
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
	panelControler.caller.checkCopy(practiceWords, panelControler.solution, panelControler.wpm)
}

func (panelControler *Controler) processTextToCopy(solution [][]*types.KeyCodeRecord) {
	panelControler.solution = solution
	panelControler.presenter.ready1()
	panelControler.tools.GoModal(
		fmt.Sprintf("The CW will begin %d seconds after you click close. Enter your copy into the red square.", panelControler.delaySeconds),
		"Copy Practice",
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

func (panelControler *Controler) processKeyStopped() {
	panelControler.codeIsKeying = false
	panelControler.presenter.keyingStopped()
}

// initialCalls runs the first code that the controler needs to run.
func (panelControler *Controler) initialCalls() {

	/* NOTE TO DEVELOPER. Step 4 of 4.

	// Make the initial calls.
	// I use this to start up widgets. For example a virtual list widget.

	*/

}
