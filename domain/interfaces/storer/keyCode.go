package storer

import (
	"github.com/josephbudd/cwt/domain/types"
)

// KeyCodeStorer defines the behavior of a KeyCode database.
type KeyCodeStorer interface {

	// Open opens the database.
	// Returns the error.
	Open() error

	// Close closes the database.
	// Returns the error.
	Close() error

	// GetKeyCode retrieves one *types.KeyCodeRecord from the db.
	// Param id [in] is the record id.
	// Returns a record pointer and error.
	// Returns (nil, nil) when the record is not found.
	GetKeyCode(id uint64) (*types.KeyCodeRecord, error)

	// GetKeyCodes retrieves all of the *types.KeyCodeRecords from the db.
	// Returns a slice of record pointers and error.
	// When no records found, the slice length is 0 and error is nil.
	GetKeyCodes() ([]*types.KeyCodeRecord, error)

	// UpdateKeyCode updates the *types.KeyCodeRecord in the database.
	// Param r [in-out] the pointer to the record to be updated.
	// If r is a new record then r.ID is updated as well.
	// Returns the error.
	UpdateKeyCode(r *types.KeyCodeRecord) error

	// RemoveKeyCode removes the types.KeyCodeRecord from the database.
	// Param id [in] the key of the record to be removed.
	// Returns the error.
	RemoveKeyCode(id uint64) error
}

