package CopyTestPanel

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

	copyTestPanel js.Value
}

func (panelGroup *PanelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelGroup *PanelGroup) defineMembers()")
		}
	}()

	notJS := panelGroup.notJS
	null := js.Null()

	if panelGroup.copyTestPanel = notJS.GetElementByID("tabsMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-TestTabPanel-inner-CopyTestPanel"); panelGroup.copyTestPanel == null {
		err = errors.New("unable to find #tabsMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-TestTabPanel-inner-CopyTestPanel")
		return
	}


	return
}

/*
	Show panel funcs.

	Call these from the controler, presenter and caller.
*/

// showCopyTestPanel shows the panel you named CopyTestPanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
A page to let the user copy and record the results to the repo.
*/
func (panelGroup *PanelGroup) showCopyTestPanel() {
	panelGroup.tools.ShowPanelInTabGroup(panelGroup.copyTestPanel)
}


