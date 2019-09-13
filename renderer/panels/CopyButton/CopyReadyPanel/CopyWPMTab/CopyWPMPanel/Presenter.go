package copywpmpanel

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/pkg/errors"
)

/*

	Panel name: CopyWPMPanel

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

	copyWPM js.Value
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

	// Define the wpm input field.
	if presenter.copyWPM = notJS.GetElementByID("copyWPM"); presenter.copyWPM == null {
		err = errors.New("unable to find #copyWPM")
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

// displayWPM displays the wpm in the select.
func (presenter *panelPresenter) displayWPM(wpm uint64) {
	notJS.SetValue(presenter.copyWPM, fmt.Sprintf("%d", wpm))
}
