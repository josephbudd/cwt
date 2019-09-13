package copypracticepanel

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/pkg/errors"
)

/*

	Panel name: CopyPracticePanel

*/

// panelPresenter writes to the panel
type panelPresenter struct {
	group          *panelGroup
	controller     *panelController
	caller         *panelCaller
	tabPanelHeader js.Value

	/* NOTE TO DEVELOPER: Step 1 of 3.

	// Declare your panelPresenter members here.

	*/

	copyPracticeStart js.Value
	copyPracticeCopy  js.Value
	copyPracticeText  js.Value
	copyPracticeCheck js.Value
}

// defineMembers defines the panelPresenter members by their html elements.
// Returns the error.
func (presenter *panelPresenter) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(presenter *panelPresenter) defineMembers()")
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 3.

	// Define your panelPresenter members.

	*/

	// Define the copy input field.
	if presenter.copyPracticeCopy = notJS.GetElementByID("copyPracticeCopy"); presenter.copyPracticeCopy == null {
		err = errors.New("unable to find #copyPracticeCopy")
		return
	}

	// Define the text output field.
	if presenter.copyPracticeText = notJS.GetElementByID("copyPracticeText"); presenter.copyPracticeText == null {
		err = errors.New("unable to find #copyPracticeText")
		return
	}

	// Define the start button.
	if presenter.copyPracticeStart = notJS.GetElementByID("copyPracticeStart"); presenter.copyPracticeStart == null {
		err = errors.New("unable to find #copyPracticeStart")
		return
	}

	// Define the check button.
	if presenter.copyPracticeCheck = notJS.GetElementByID("copyPracticeCheck"); presenter.copyPracticeCheck == null {
		err = errors.New("unable to find #copyPracticeCheck")
		return
	}

	return
}

// Tab panel heading.

func (presenter *panelPresenter) getTabPanelHeading() (heading string) {
	heading = notJS.GetInnerText(presenter.tabPanelHeader)
	return
}

func (presenter *panelPresenter) setTabPanelHeading(heading string) {
	heading = strings.TrimSpace(heading)
	if len(heading) == 0 {
		tools.ElementHide(presenter.tabPanelHeader)
	} else {
		tools.ElementShow(presenter.tabPanelHeader)
	}
	notJS.SetInnerText(presenter.tabPanelHeader, heading)
}

/* NOTE TO DEVELOPER. Step 3 of 3.

// Define your panelPresenter functions.

*/

func (presenter *panelPresenter) started() {
	tools.ElementShow(presenter.copyPracticeCheck)
	tools.ElementHide(presenter.copyPracticeStart)
	notJS.ClearValue(presenter.copyPracticeCopy)
	notJS.RemoveChildNodes(presenter.copyPracticeText)
}

func (presenter *panelPresenter) checked() {
	tools.ElementHide(presenter.copyPracticeCheck)
	tools.ElementShow(presenter.copyPracticeStart)
}

func (presenter *panelPresenter) ready1() {
	notJS.ClearValue(presenter.copyPracticeCopy)
}

func (presenter *panelPresenter) ready2() {
	notJS.Focus(presenter.copyPracticeCopy)
}

func (presenter *panelPresenter) keyingFinished() {
	notJS.Focus(presenter.copyPracticeCopy)
}

func (presenter *panelPresenter) keyingStopped() {
	notJS.Focus(presenter.copyPracticeCopy)
}

func (presenter *panelPresenter) showResult(correctCount, incorrectCount, keyedCount uint64, testResults [][]data.TestResult) {
	div := presenter.copyPracticeText
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
