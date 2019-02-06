package CopyWPMPanel

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

	copyWPMPanel js.Value
}

func (panelGroup *PanelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelGroup *PanelGroup) defineMembers()")
		}
	}()

	notJS := panelGroup.notJS
	null := js.Null()

	if panelGroup.copyWPMPanel = notJS.GetElementByID("tabsMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-WPMTabPanel-inner-CopyWPMPanel"); panelGroup.copyWPMPanel == null {
		err = errors.New("unable to find #tabsMasterView_home_pad_CopyButton_CopyReadyPanel_tab_bar-WPMTabPanel-inner-CopyWPMPanel")
		return
	}


	return
}

/*
	Show panel funcs.

	Call these from the controler, presenter and caller.
*/

// showCopyWPMPanel shows the panel you named CopyWPMPanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
Let the user select the words per miniute for copying.
*/
func (panelGroup *PanelGroup) showCopyWPMPanel() {
	panelGroup.tools.ShowPanelInTabGroup(panelGroup.copyWPMPanel)
}


