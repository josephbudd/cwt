package storer

import (
	"github.com/josephbudd/cwt/domain/types"
)

// WPMStorer defines the behavior of a WPM database.
type WPMStorer interface {

	// Open opens the database.
	// Returns the error.
	Open() error

	// Close closes the database.
	// Returns the error.
	Close() error

	// GetWPM retrieves one *types.WPMRecord from the db.
	// Param id [in] is the record id.
	// Returns a record pointer and error.
	// Returns (nil, nil) when the record is not found.
	GetWPM(id uint64) (*types.WPMRecord, error)

	// GetCopyWPM retrieves the copy types.WPMRecord from the db.
	// Returns a record pointer and error.
	GetCopyWPM() (*types.WPMRecord, error)

	// GetKeyWPM retrieves the key types.WPMRecord from the db.
	// Returns a record pointer and error.
	GetKeyWPM() (*types.WPMRecord, error)

	// GetWPMs retrieves all of the *types.WPMRecords from the db.
	// Returns a slice of record pointers and error.
	// When no records found, the slice length is 0 and error is nil.
	GetWPMs() ([]*types.WPMRecord, error)

	// UpdateWPM updates the *types.WPMRecord in the database.
	// Param r [in-out] the pointer to the record to be updated.
	// If r is a new record then r.ID is updated as well.
	// Returns the error.
	UpdateWPM(r *types.WPMRecord) error

	// RemoveWPM removes the types.WPMRecord from the database.
	// Param id [in] the key of the record to be removed.
	// Returns the error.
	RemoveWPM(id uint64) error
}
