package settings

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/data/filepaths"
	"github.com/josephbudd/cwt/domain/types"
	"github.com/josephbudd/cwtsitepack"
)

// NewApplicationSettings makes a new ApplicationSettings.
func NewApplicationSettings() (settings *types.ApplicationSettings, err error) {
	var fpath string
	var contents []byte
	var found bool
	fpath = filepaths.GetShortSettingsPath()
	if contents, found = cwtsitepack.Contents(fpath); !found {
		err = errors.New(fmt.Sprintf("can't find %q", fpath))
		return
	}
	settings = &types.ApplicationSettings{}
	if err = yaml.Unmarshal(contents, settings); err != nil {
		return
	}
	return
}
