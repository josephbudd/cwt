package KeyPracticePanel

import (
	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/renderer/notjs"
	"github.com/josephbudd/cwt/renderer/viewtools"
)

/*

	Panel name: KeyPracticePanel

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

	return
}

/* NOTE TO DEVELOPER. Step 3 of 3.

// Define your Presenter functions.

*/
