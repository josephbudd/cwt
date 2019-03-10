package numberspanel

import (
	"syscall/js"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

// PanelGroup is a group of 1 panel.
// It also has a show panel func for each panel in this panel group.
type PanelGroup struct {
	tools *viewtools.Tools
	notJS *notjs.NotJS

	numbersPanel js.Value
}

func (panelGroup *PanelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelGroup *PanelGroup) defineMembers()")
		}
	}()

	notJS := panelGroup.notJS
	null := js.Null()

	if panelGroup.numbersPanel = notJS.GetElementByID("tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-NumbersTabPanel-inner-NumbersPanel"); panelGroup.numbersPanel == null {
		err = errors.New("unable to find #tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-NumbersTabPanel-inner-NumbersPanel")
		return
	}


	return
}

/*
	Show panel funcs.

	Call these from the controler, presenter and caller.
*/

// showNumbersPanel shows the panel you named NumbersPanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
This panel is displayed when the "Numbers" tab button is clicked.
This panel is the only panel in it's panel group.

*/
func (panelGroup *PanelGroup) showNumbersPanel() {
	panelGroup.tools.ShowPanelInTabGroup(panelGroup.numbersPanel)
}


