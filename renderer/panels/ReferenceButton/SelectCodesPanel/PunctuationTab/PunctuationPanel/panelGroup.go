package punctuationpanel

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

	punctuationPanel js.Value
}

func (panelGroup *PanelGroup) defineMembers() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "(panelGroup *PanelGroup) defineMembers()")
		}
	}()

	notJS := panelGroup.notJS
	null := js.Null()

	if panelGroup.punctuationPanel = notJS.GetElementByID("tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-PunctuationTabPanel-inner-PunctuationPanel"); panelGroup.punctuationPanel == null {
		err = errors.New("unable to find #tabsMasterView_home_pad_ReferenceButton_SelectCodesPanel_tab_bar-PunctuationTabPanel-inner-PunctuationPanel")
		return
	}


	return
}

/*
	Show panel funcs.

	Call these from the controler, presenter and caller.
*/

// showPunctuationPanel shows the panel you named PunctuationPanel while hiding any other panels in this panel group.
// This panel will become visible only when this group of panels becomes visible.
/* Your note for this panel is:
This panel is displayed when the "Punctuation" tab button is clicked.
This panel is the only panel in it's panel group.

*/
func (panelGroup *PanelGroup) showPunctuationPanel() {
	panelGroup.tools.ShowPanelInTabGroup(panelGroup.punctuationPanel)
}


