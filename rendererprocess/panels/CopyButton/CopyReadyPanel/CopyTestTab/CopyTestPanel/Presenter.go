// +build js, wasm

package copytestpanel

import (
	"errors"
	"fmt"
	"strings"

	"github.com/josephbudd/cwt/domain/data"
	"github.com/josephbudd/cwt/rendererprocess/api/display"
	"github.com/josephbudd/cwt/rendererprocess/api/markup"
)

/*

	Panel name: CopyTestPanel

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

	copyTestStart *markup.Element
	copyTestCopy  *markup.Element
	copyTestText  *markup.Element
	copyTestCheck *markup.Element
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
	if presenter.copyTestCopy = document.ElementByID("copyTestCopy"); presenter.copyTestCopy == nil {
		err = errors.New("unable to find #copyTestCopy")
		return
	}

	// Define the text output field.
	if presenter.copyTestText = document.ElementByID("copyTestText"); presenter.copyTestText == nil {
		err = errors.New("unable to find #copyTestText")
		return
	}

	// Define the start button.
	if presenter.copyTestStart = document.ElementByID("copyTestStart"); presenter.copyTestStart == nil {
		err = errors.New("unable to find #copyTestStart")
		return
	}

	// Define the check button.
	if presenter.copyTestCheck = document.ElementByID("copyTestCheck"); presenter.copyTestCheck == nil {
		err = errors.New("unable to find #copyTestCheck")
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
	presenter.copyTestCheck.Show()
	presenter.copyTestStart.Hide()
	presenter.copyTestCopy.ClearValue()
	presenter.copyTestText.RemoveChildren()
}

func (presenter *panelPresenter) checked() {
	presenter.copyTestCheck.Hide()
	presenter.copyTestStart.Show()
}

func (presenter *panelPresenter) ready1() {
	presenter.copyTestCopy.ClearValue()
}

func (presenter *panelPresenter) ready2() {
	presenter.copyTestCopy.Focus()
}

func (presenter *panelPresenter) keyingFinished() {
	presenter.copyTestCopy.Focus()
}

func (presenter *panelPresenter) keyingStopped() {
	presenter.copyTestCopy.Focus()
}

func (presenter *panelPresenter) showResult(correctCount, incorrectCount, keyedCount uint64, testResults [][]data.TestResult) {
	// heading
	h3 := document.NewH3()
	h3.SetInnerText(fmt.Sprintf("You got %d correct and %d incorrect out of %d.", correctCount, incorrectCount, keyedCount))
	presenter.copyTestText.AppendChild(h3)
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
	presenter.copyTestText.AppendChild(table)
}

func (presenter *panelPresenter) addTableRow(tbody *markup.Element, nColumns int, matches []bool, controls, ditDahs, copieds []string) {
	var tr, th, td *markup.Element
	// control characters
	tr = document.NewTR()
	th = document.NewTH()
	th.SetInnerText("keyed")
	tr.AppendChild(th)
	var i int
	var s string
	for i, s = range controls {
		td = document.NewTD()
		td.AppendText(s)
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
		td.AppendText(" ")
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
	th.AddClass("copied")
	th.SetInnerText("You Copyed")
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
