// +build js, wasm

package copypracticepanel

import (
	"fmt"
	"strings"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
	"github.com/josephbudd/cwt/rendererprocess/api/markup"
)

/*

	Panel name: CopyPracticePanel

*/

// panelPresenter writes to the panel
type panelPresenter struct {
	group          *panelGroup
	controller     *panelController
	messenger      *panelMessenger
	tabPanelHeader *markup.Element
	tabButton      *markup.Element

	/* NOTE TO DEVELOPER: Step 1 of 3.

	// Declare your panelPresenter members here.

	// example:

	editCustomerName *markup.Element

	*/

	copyPracticeStart *markup.Element
	copyPracticeCopy  *markup.Element
	copyPracticeText  *markup.Element
	copyPracticeCheck *markup.Element
}

// defineMembers defines the panelPresenter members by their html elements.
// Returns the error.
func (presenter *panelPresenter) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("(presenter *panelPresenter) defineMembers(): %w", err)
		}
	}()

	/* NOTE TO DEVELOPER. Step 2 of 3.

	// Define your panelPresenter members.

	// example:

	// Define the edit form's customer name input field.
	if presenter.editCustomerName = document.ElementByID("editCustomerName"); presenter.editCustomerName == nil {
		err = fmt.Errorf("unable to find #editCustomerName")
		return
	}

	*/

	// Define the copy input field.
	if presenter.copyPracticeCopy = document.ElementByID("copyPracticeCopy"); presenter.copyPracticeCopy == nil {
		err = fmt.Errorf("unable to find #copyPracticeCopy")
		return
	}

	// Define the text output field.
	if presenter.copyPracticeText = document.ElementByID("copyPracticeText"); presenter.copyPracticeText == nil {
		err = fmt.Errorf("unable to find #copyPracticeText")
		return
	}

	// Define the start button.
	if presenter.copyPracticeStart = document.ElementByID("copyPracticeStart"); presenter.copyPracticeStart == nil {
		err = fmt.Errorf("unable to find #copyPracticeStart")
		return
	}

	// Define the check button.
	if presenter.copyPracticeCheck = document.ElementByID("copyPracticeCheck"); presenter.copyPracticeCheck == nil {
		err = fmt.Errorf("unable to find #copyPracticeCheck")
		return
	}

	return
}

// Tab button label.

func (presenter *panelPresenter) getTabLabel() (label string) {
	label = presenter.tabButton.InnerText()
	return
}

func (presenter *panelPresenter) setTabLabel(label string) {
	presenter.tabButton.SetInnerText(label)
}

// Tab panel heading.

func (presenter *panelPresenter) getTabPanelHeading() (heading string) {
	heading = presenter.tabPanelHeader.InnerText()
	return
}

func (presenter *panelPresenter) setTabPanelHeading(heading string) {
	heading = strings.TrimSpace(heading)
	if len(heading) == 0 {
		presenter.tabPanelHeader.Hide()
	} else {
		presenter.tabPanelHeader.Show()
	}
	presenter.tabPanelHeader.SetInnerText(heading)
}

/* NOTE TO DEVELOPER. Step 3 of 3.

// Define your panelPresenter functions.

// example:

// displayCustomer displays the customer in the edit customer form panel.
func (presenter *panelPresenter) displayCustomer(record *types.CustomerRecord) {
	presenter.editCustomerName.SetValue(record.Name)
}

*/

func (presenter *panelPresenter) started() {
	presenter.copyPracticeCheck.Show()
	presenter.copyPracticeStart.Hide()
	presenter.copyPracticeCopy.ClearValue()
	presenter.copyPracticeText.RemoveChildren()
}

func (presenter *panelPresenter) checked() {
	presenter.copyPracticeCheck.Hide()
	presenter.copyPracticeStart.Show()
}

func (presenter *panelPresenter) ready1() {
	presenter.copyPracticeCopy.ClearValue()
}

func (presenter *panelPresenter) ready2() {
	presenter.copyPracticeCopy.Focus()
}

func (presenter *panelPresenter) keyingFinished() {
	presenter.copyPracticeCopy.Focus()
}

func (presenter *panelPresenter) keyingStopped() {
	presenter.copyPracticeCopy.Focus()
}

func (presenter *panelPresenter) showResult(correctCount, incorrectCount, keyedCount uint64, testResults [][]data.TestResult) {
	// heading
	h3 := document.NewH3()
	h3.SetInnerText(fmt.Sprintf("You got %d correct and %d incorrect out of %d.", correctCount, incorrectCount, keyedCount))
	presenter.copyPracticeText.AppendChild(h3)
	// calculate max columns
	max := 0
	for _, mmWord := range testResults {
		l := len(mmWord)
		if max < l {
			max = l
		}
	}
	tbody := presenter.startTable()
	// details
	// put each word in its own table
	for _, mmWord := range testResults {
		// results
		l := len(mmWord)
		controlChars := make([]string, l)
		copiedChars := make([]string, l)
		ditdahs := make([]string, l)
		matches := make([]bool, l)
		for i, mmChar := range mmWord {
			matches[i] = mmChar.Control.ID == mmChar.Input.ID
			controlChars[i] = mmChar.Control.Character
			copiedChars[i] = mmChar.Input.Character
			ditdahs[i] = mmChar.Control.DitDah
		}
		presenter.addTableRow(tbody, max, matches, controlChars, ditdahs, copiedChars)
	}
	presenter.endTable(tbody)
	display.Resize()
}

func (presenter *panelPresenter) startTable() (tbody *markup.Element) {
	tbody = document.NewTBODY()
	return
}

func (presenter *panelPresenter) endTable(tbody *markup.Element) {
	table := document.NewTABLE()
	table.AppendChild(tbody)
	presenter.copyPracticeText.AppendChild(table)
}

func (presenter *panelPresenter) addTableRow(tbody *markup.Element, nColumns int, matches []bool, controls, ditDahs, copieds []string) {
	var tr, th, td *markup.Element
	// control characters
	tr = document.NewTR()
	th = document.NewTH()
	th.AddClass("keyed")
	th.SetInnerText("Text")
	tr.AppendChild(th)
	var i int
	var s string
	for i, s = range controls {
		td = document.NewTD()
		td.SetInnerText(s)
		td.AddClass("keyed")
		if matches[i] {
			td.AddClass("match")
		} else {
			td.AddClass("mismatch")
		}
		tr.AppendChild(td)
	}
	for i++; i < nColumns; i++ {
		td = document.NewTD()
		td.SetInnerText(" ")
		td.AddClasses("keyed", "mismatch")
		tr.AppendChild(td)
	}
	tbody.AppendChild(tr)
	// dit-dahs
	tr = document.NewTR()
	th = document.NewTH()
	th.AddClass("ditdah")
	th.SetInnerText("Morse Code")
	tr.AppendChild(th)
	for i, s = range ditDahs {
		td = document.NewTD()
		td.SetInnerText(s)
		td.AddClass("ditdah")
		if matches[i] {
			td.AddClass("match")
		} else {
			td.AddClass("mismatch")
		}
		tr.AppendChild(td)
	}
	for i++; i < nColumns; i++ {
		td = document.NewTD()
		td.SetInnerText(" ")
		td.AddClasses("ditdah", "mismatch")
		tr.AppendChild(td)
	}
	tbody.AppendChild(tr)
	// copy
	tr = document.NewTR()
	th = document.NewTH()
	th.SetInnerText("You Copyed")
	th.AddClass("copied")
	tr.AppendChild(th)
	for i, s = range copieds {
		td = document.NewTD()
		td.SetInnerText(s)
		td.AddClass("copied")
		if matches[i] {
			td.AddClass("match")
		} else {
			td.AddClass("mismatch")
		}
		tr.AppendChild(td)
	}
	for i++; i < nColumns; i++ {
		td = document.NewTD()
		td.SetInnerText(" ")
		td.AddClasses("copied", "mismatch")
		tr.AppendChild(td)
	}
	tbody.AppendChild(tr)
}
