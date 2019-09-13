package storing

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/josephbudd/cwt/domain/store/record"
)

const (
	keyID         uint64 = 1
	copyID        uint64 = 2
	wPMBucketName string = "wPM"
)

// WPMLocalBoltStore is the API of the WPM local bolt store.
// It is the implementation of the interface in /domain/store/storerWPM.go.
type WPMLocalBoltStore struct {
	DB    *bolt.DB
	path  string
	perms os.FileMode
}

// NewWPMLocalBoltStore constructs a new WPMLocalBoltStore.
// Param db is an open bolt data-store.
// Returns a pointer to the new WPMLocalBoltStore.
func NewWPMLocalBoltStore(path string, perms os.FileMode) (store *WPMLocalBoltStore) {
	store = &WPMLocalBoltStore{
		path:  path,
		perms: perms,
	}
	return
}

// Open opens the bolt data-store.
// Returns the error.
func (store *WPMLocalBoltStore) Open() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "WPMLocalBoltStore.Open")
		}
	}()

	if store.DB, err = bolt.Open(store.path, store.perms, nil); err != nil {
		err = errors.WithMessage(err, "bolt.Open(path, filepaths.GetFmode(), nil)")
	}
	return
}

// Close closes the bolt data-store.
// Returns the error.
func (store *WPMLocalBoltStore) Close() (err error) {

	defer func() {
		if err != nil {
			err = errors.WithMessage(err, "WPMLocalBoltStore.Close")
		}
	}()

	err = store.DB.Close()
	return
}

// Get retrieves the record.WPM from the bolt data-store.
// Param id is the record id.
// Returns the record and error.
// If no record is found returns a nil record and a nil error.
func (store *WPMLocalBoltStore) Get(id uint64) (r *record.WPM, err error) {
	ids := fmt.Sprintf("%d", id)
	err = store.DB.View(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(wPMBucketName)
		var bucket *bolt.Bucket
		if bucket = tx.Bucket(bucketname); bucket == nil {
			return
		}
		var rbb []byte
		if rbb = bucket.Get([]byte(ids)); rbb == nil {
			// not found
			return
		}
		r = &record.WPM{}
		if er = json.Unmarshal(rbb, r); er != nil {
			r = nil
			return
		}
		return
	})
	return
}

// GetAll retrieves all of the record.WPM from the bolt data-store.
// If no record is found then it calls store.initialize() and tries again. See *WPMLocalBoltStore.initialize().
// Returns the records and error.
// If no record is found returns a zero length records and a nil error.
func (store *WPMLocalBoltStore) GetAll() (rr []*record.WPM, err error) {
	if rr, err = store.getAll(); len(rr) == 0 && err == nil {
		store.initialize()
		rr, err = store.getAll()
	}
	return
}

func (store *WPMLocalBoltStore) getAll() (rr []*record.WPM, err error) {
	err = store.DB.View(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(wPMBucketName)
		var bucket *bolt.Bucket
		if bucket = tx.Bucket(bucketname); bucket == nil {
			return
		}
		c := bucket.Cursor()
		rr = make([]*record.WPM, 0, 1024)
		for k, v := c.First(); k != nil; k, v = c.Next() {
			r := record.NewWPM()
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

// Update adds or updates a record.WPM in the bolt data-store.
// Param r is the record to be updated.
// If r is a new record then r.ID is updated with the new record id.
// Returns the error.
func (store *WPMLocalBoltStore) Update(r *record.WPM) (err error) {
	err = store.update(r)
	return
}

// Remove removes a record.WPM from the bolt data-store.
// Param id the key of the record to be removed.
// If the record is not found returns a nil error.
// Returns the error.
func (store *WPMLocalBoltStore) Remove(id uint64) (err error) {
	err = store.DB.Update(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(wPMBucketName)
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

// GetKeyWPM retrieves the key record.WPM from the db.
// Returns the record and error.
func (store *WPMLocalBoltStore) GetKeyWPM() (r *record.WPM, err error) {
	var rr []*record.WPM
	if rr, err = store.GetAll(); err != nil {
		return
	}
	for _, r = range rr {
		if r.ID == keyID {
			return
		}
	}
	r = nil
	return
}

// GetCopyWPM retrieves the copy record.WPM from the db.
// Returns the record and error.
func (store *WPMLocalBoltStore) GetCopyWPM() (r *record.WPM, err error) {
	var rr []*record.WPM
	if rr, err = store.GetAll(); err != nil {
		return
	}
	for _, r = range rr {
		if r.ID == copyID {
			return
		}
	}
	r = nil
	return
}

// updates the record.WPM in the bolt data-store.
// Param record the record to be updated.
// If the record is new then it's ID is updated.
// Returns the error.
func (store *WPMLocalBoltStore) update(r *record.WPM) (err error) {
	err = store.DB.Update(func(tx *bolt.Tx) (er error) {
		bucketname := []byte(wPMBucketName)
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

// initialize is only useful if you want to add the default records to the bolt data-store.
// otherwise you don't need it to do anything.
func (store *WPMLocalBoltStore) initialize() (err error) {
	// create the key wpm record of 5 wpm.
	r := record.NewWPM()
	r.WPM = 5
	r.ID = keyID
	if err = store.update(r); err != nil {
		return
	}
	// create the copy wpm record of 5 wpm.
	r.ID = copyID
	r.WPM = 15
	err = store.update(r)
	return
}
