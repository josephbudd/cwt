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

// Presenter writes to the panel
type Presenter struct {
	panelGroup *PanelGroup
	controler  *Controler
	caller     *Caller
	tools      *viewtools.Tools // see /renderer/viewtools
	notJS      *notjs.NotJS

	/* NOTE TO DEVELOPER: Step 1 of 3.

	// Declare your Presenter members here.

	*/

	copyPracticeStart js.Value
	copyPracticeCopy  js.Value
	copyPracticeText  js.Value
	copyPracticeCheck js.Value
}

// defineMembers defines the Presenter members by their html elements.
// Returns the error.
func (panelPresenter *Presenter) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelPresenter *Presenter) defineMembers()")
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 3.

	// Define your Presenter members.

	*/

	notJS := panelPresenter.notJS
	null := js.Null()

	// Define the copy input field.
	if panelPresenter.copyPracticeCopy = notJS.GetElementByID("copyPracticeCopy"); panelPresenter.copyPracticeCopy == null {
		err = errors.New("unable to find #copyPracticeCopy")
		return
	}

	// Define the text output field.
	if panelPresenter.copyPracticeText = notJS.GetElementByID("copyPracticeText"); panelPresenter.copyPracticeText == null {
		err = errors.New("unable to find #copyPracticeText")
		return
	}

	// Define the start button.
	if panelPresenter.copyPracticeStart = notJS.GetElementByID("copyPracticeStart"); panelPresenter.copyPracticeStart == null {
		err = errors.New("unable to find #copyPracticeStart")
		return
	}

	// Define the check button.
	if panelPresenter.copyPracticeCheck = notJS.GetElementByID("copyPracticeCheck"); panelPresenter.copyPracticeCheck == null {
		err = errors.New("unable to find #copyPracticeCheck")
		return
	}

	return
}

/* NOTE TO DEVELOPER. Step 3 of 3.

// Define your Presenter functions.

*/

func (panelPresenter *Presenter) started() {
	tools := panelPresenter.tools
	tools.ElementShow(panelPresenter.copyPracticeCheck)
	tools.ElementHide(panelPresenter.copyPracticeStart)
	notJS := panelPresenter.notJS
	notJS.ClearValue(panelPresenter.copyPracticeCopy)
	notJS.RemoveChildNodes(panelPresenter.copyPracticeText)
}

func (panelPresenter *Presenter) checked() {
	tools := panelPresenter.tools
	tools.ElementHide(panelPresenter.copyPracticeCheck)
	tools.ElementShow(panelPresenter.copyPracticeStart)
}

func (panelPresenter *Presenter) ready1() {
	notJS := panelPresenter.notJS
	notJS.ClearValue(panelPresenter.copyPracticeCopy)
}

func (panelPresenter *Presenter) ready2() {
	notJS := panelPresenter.notJS
	notJS.Focus(panelPresenter.copyPracticeCopy)
}

func (panelPresenter *Presenter) keyingFinished() {
	notJS := panelPresenter.notJS
	notJS.Focus(panelPresenter.copyPracticeCopy)
}

func (panelPresenter *Presenter) showResult(correctCount, incorrectCount, keyedCount uint64, testResults [][]types.TestResult) {
	div := panelPresenter.copyPracticeText
	notJS := panelPresenter.notJS
	// heading
	h3 := notJS.CreateElementH3()
	tn := notJS.CreateTextNode(fmt.Sprintf("You got %d correct and %d incorrect out of %d.", correctCount, incorrectCount, keyedCount))
	notJS.AppendChild(h3, tn)
	notJS.AppendChild(div, h3)
	// details
	// put each word in its own table
	for _, mmWord := range testResults {
		// each word in it's own table
		table := notJS.CreateElementTABLE()
		tbody := notJS.CreateElementTBODY()
		// header
		tr := notJS.CreateElementTR()
		th := notJS.CreateElementTH()
		tn = notJS.CreateTextNode("Code")
		notJS.AppendChild(th, tn)
		notJS.AppendChild(tr, th)
		th = notJS.CreateElementTH()
		tn = notJS.CreateTextNode("Keyed")
		notJS.AppendChild(th, tn)
		notJS.AppendChild(tr, th)
		th = notJS.CreateElementTH()
		tn = notJS.CreateTextNode("Copied")
		notJS.AppendChild(th, tn)
		notJS.AppendChild(tr, th)
		notJS.AppendChild(tbody, tr)
		// results
		keyedChars := make([]string, 0, 20)
		copiedChars := make([]string, 0, 20)
		ditdahs := make([]string, 0, 20)
		for _, mmChar := range mmWord {
			ditdahs = append(ditdahs, mmChar.Control.DitDah)
			keyedChars = append(keyedChars, mmChar.Control.Character)
			copiedChars = append(copiedChars, mmChar.Input.Character)
		}
		ddWord := strings.Join(ditdahs, " ")
		keyedWord := strings.Join(keyedChars, ",")
		copiedWord := strings.Join(copiedChars, ",")
		var tdClass string
		if keyedWord == copiedWord {
			tdClass = "match"
		} else {
			tdClass = "mismatch"
		}
		tr = notJS.CreateElementTR()
		td := notJS.CreateElementTD()
		notJS.AppendChild(td, notJS.CreateTextNode(ddWord))
		notJS.ClassListAddClass(td, tdClass)
		notJS.ClassListAddClass(td, "ditdah")
		notJS.AppendChild(tr, td)
		td = notJS.CreateElementTD()
		notJS.AppendChild(td, notJS.CreateTextNode(keyedWord))
		notJS.ClassListAddClass(td, tdClass)
		notJS.ClassListAddClass(td, "keyed")
		notJS.AppendChild(tr, td)
		td = notJS.CreateElementTD()
		notJS.AppendChild(td, notJS.CreateTextNode(copiedWord))
		notJS.ClassListAddClass(td, tdClass)
		notJS.ClassListAddClass(td, "copied")
		notJS.AppendChild(tr, td)
		notJS.AppendChild(tbody, tr)
		notJS.AppendChild(table, tbody)
		notJS.AppendChild(div, table)
	}
}
