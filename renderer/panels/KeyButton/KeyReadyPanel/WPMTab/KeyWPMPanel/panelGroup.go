package keywpmpanel

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

	keyWPMPanel js.Value
}

func (panelGroup *PanelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelGroup *PanelGroup) defineMembers()")
		}
	}()

	notJS := panelGroup.notJS
	null := js.Null()

	if panelGroup.keyWPMPanel = notJS.GetElementByID("tabsMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-WPMTabPanel-inner-KeyWPMPanel"); panelGroup.keyWPMPanel == null {
		err = errors.New("unable to find #tabsMasterView_home_pad_KeyButton_KeyReadyPanel_tab_bar-WPMTabPanel-inner-KeyWPMPanel")
		return
	}


	return
}

/*
	Show panel funcs.

	Call these from the controler, presenter and caller.
*/

// showKeyWPMPanel shows the panel you named KeyWPMPanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
Let the user select the words per miniute for Keying.
*/
func (panelGroup *PanelGroup) showKeyWPMPanel() {
	panelGroup.tools.ShowPanelInTabGroup(panelGroup.keyWPMPanel)
}


