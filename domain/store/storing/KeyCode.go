package storing

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/boltdb/bolt"

	"github.com/josephbudd/cwt/domain/store/record"
)

const keyCodeBucketName string = "keyCode"

// KeyCodeLocalBoltStore is the API of the KeyCode local bolt store.
// It is the implementation of the interface in /domain/store/storerKeyCode.go.
type KeyCodeLocalBoltStore struct {
	DB    *bolt.DB
	path  string
	perms os.FileMode
}

// NewKeyCodeLocalBoltStore constructs a new KeyCodeLocalBoltStore.
// Param db is an open bolt data-store.
// Returns a pointer to the new KeyCodeLocalBoltStore.
func NewKeyCodeLocalBoltStore(path string, perms os.FileMode) (store *KeyCodeLocalBoltStore) {
	store = &KeyCodeLocalBoltStore{
		path:  path,
		perms: perms,
	}
	return
}

// Open opens the bolt data-store.
// Returns the error.
func (store *KeyCodeLocalBoltStore) Open() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeLocalBoltStore.Open: %w", err)
		}
	}()

	if store.DB, err = bolt.Open(store.path, store.perms, nil); err != nil {
		err = fmt.Errorf("bolt.Open(path, filepaths.GetFmode(), nil): %w", err)
	}
	return
}

// Close closes the bolt data-store.
// Returns the error.
func (store *KeyCodeLocalBoltStore) Close() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("KeyCodeLocalBoltStore.Close: %w", err)
		}
	}()

	err = store.DB.Close()
	return
}

// Get retrieves the record.KeyCode from the bolt data-store.
// Param id is the record id.
// Returns the record and error.
// If no record is found returns a nil record and a nil error.
func (store *KeyCodeLocalBoltStore) Get(id uint64) (r *record.KeyCode, err error) {
	ids := fmt.Sprintf("%d", id)
	err = store.DB.View(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(keyCodeBucketName)
		var bucket *bolt.Bucket
		if bucket = tx.Bucket(bucketname); bucket == nil {
			return
		}
		var rbb []byte
		if rbb = bucket.Get([]byte(ids)); rbb == nil {
			// not found
			return
		}
		r = &record.KeyCode{}
		if er = json.Unmarshal(rbb, r); er != nil {
			r = nil
			return
		}
		return
	})
	return
}

// GetAll retrieves all of the record.KeyCode from the bolt data-store.
// If no record is found then it calls store.initialize() and tries again. See *KeyCodeLocalBoltStore.initialize().
// Returns the records and error.
// If no record is found returns a zero length records and a nil error.
func (store *KeyCodeLocalBoltStore) GetAll() (rr []*record.KeyCode, err error) {
	if rr, err = store.getAll(); len(rr) == 0 && err == nil {
		store.initialize()
		rr, err = store.getAll()
	}
	return
}

func (store *KeyCodeLocalBoltStore) getAll() (rr []*record.KeyCode, err error) {
	err = store.DB.View(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(keyCodeBucketName)
		var bucket *bolt.Bucket
		if bucket = tx.Bucket(bucketname); bucket == nil {
			return
		}
		c := bucket.Cursor()
		rr = make([]*record.KeyCode, 0, 1024)
		for k, v := c.First(); k != nil; k, v = c.Next() {
			r := record.NewKeyCode()
			if er = json.Unmarshal(v, r); er != nil {
				rr = nil
				return
			}
			rr = append(rr, r)
		}
		return
	})
	return
}

// Update adds or updates a record.KeyCode in the bolt data-store.
// Param r is the record to be updated.
// If r is a new record then r.ID is updated with the new record id.
// Returns the error.
func (store *KeyCodeLocalBoltStore) Update(r *record.KeyCode) (err error) {
	err = store.update(r)
	return
}

// Remove removes a record.KeyCode from the bolt data-store.
// Param id the key of the record to be removed.
// If the record is not found returns a nil error.
// Returns the error.
func (store *KeyCodeLocalBoltStore) Remove(id uint64) (err error) {
	err = store.DB.Update(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(keyCodeBucketName)
		var bucket *bolt.Bucket
		if bucket = tx.Bucket(bucketname); bucket == nil {
			return
		}
		idbb := []byte(fmt.Sprintf("%d", id))
		er = bucket.Delete(idbb)
		return
	})
	return
}

// updates the record.KeyCode in the bolt data-store.
// Param record the record to be updated.
// If the record is new then it's ID is updated.
// Returns the error.
func (store *KeyCodeLocalBoltStore) update(r *record.KeyCode) (err error) {
	err = store.DB.Update(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(keyCodeBucketName)
		var bucket *bolt.Bucket
		if bucket, er = tx.CreateBucketIfNotExists(bucketname); er != nil {
			return
		}
		if r.ID == 0 {
			if r.ID, er = bucket.NextSequence(); er != nil {
				return
			}
		}
		var rbb []byte
		if rbb, er = json.Marshal(r); er != nil {
			return
		}
		idbb := []byte(fmt.Sprintf("%d", r.ID))
		er = bucket.Put(idbb, rbb)
		return
	})
	return
}
