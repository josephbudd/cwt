package boltstoring

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/josephbudd/cwt/domain/types"
)

const (
	keyID  uint64 = 1
	copyID uint64 = 2
)

const wPMBucketName string = "wPM"

// WPMBoltDB is the bolt db for key codes.
type WPMBoltDB struct {
	DB    *bolt.DB
	path  string
	perms os.FileMode
}

// NewWPMBoltDB constructs a new WPMBoltDB.
// Param db [in-out] is an open bolt database.
// Returns a pointer to the new WPMBoltDB.
func NewWPMBoltDB(db *bolt.DB, path string, perms os.FileMode) *WPMBoltDB {
	return &WPMBoltDB{
		DB:    db,
		path:  path,
		perms: perms,
	}
}

// WPMBoltDB implements WPMStorer
// which is defined in domain/types/records.go

// Open opens the database.
// Returns the error.
func (wPMdb *WPMBoltDB) Open() error {
	// the bolt db is already open
	return nil
}

// Close closes the database.
// Returns the error.
func (wPMdb *WPMBoltDB) Close() error {
	return wPMdb.DB.Close()
}

// GetKeyWPM retrieves the key types.WPMRecord from the db.
// Returns the record and error.
func (wPMdb *WPMBoltDB) GetKeyWPM() (record *types.WPMRecord, err error) {
	var rr []*types.WPMRecord
	if rr, err = wPMdb.GetWPMs(); err != nil {
		return
	}
	for _, record = range rr {
		if record.ID == keyID {
			return
		}
	}
	record = nil
	return
}

// GetCopyWPM retrieves the copy types.WPMRecord from the db.
// Returns the record and error.
func (wPMdb *WPMBoltDB) GetCopyWPM() (record *types.WPMRecord, err error) {
	var rr []*types.WPMRecord
	if rr, err = wPMdb.GetWPMs(); err != nil {
		return
	}
	for _, record = range rr {
		if record.ID == copyID {
			return
		}
	}
	record = nil
	return
}

// GetWPM retrieves the types.WPMRecord from the db.
// Param id [in] is the record id.
// Returns the record and error.
func (wPMdb *WPMBoltDB) GetWPM(id uint64) (*types.WPMRecord, error) {
	var r types.WPMRecord
	ids := fmt.Sprintf("%d", id)
	er := wPMdb.DB.View(func(tx *bolt.Tx) error {
		bucketname := []byte(wPMBucketName)
		bucket := tx.Bucket(bucketname)
		if bucket != nil {
			bb := bucket.Get([]byte(ids))
			if bb != nil {
				// found
				err := json.Unmarshal(bb, &r)
				if err == nil {
					r.ID = id
				}
				return err
			}
		}
		// no bucket or not found
		return errNotFound
	})
	if er == nil {
		// found
		return &r, nil
	} else if er == errNotFound {
		// not found
		return nil, nil
	}
	return nil, er
}

// GetWPMs retrieves all of the types.WPMRecord from the db.
// If there are no types.WPMRecords in the db then it calls wPMdb.initialize().
// See wPMdb.initialize().
// Returns the records and error.
func (wPMdb *WPMBoltDB) GetWPMs() ([]*types.WPMRecord, error) {
	if rr, err := wPMdb.getWPMs(); len(rr) > 0 && err == nil {
		return rr, err
	}
	wPMdb.initialize()
	return wPMdb.getWPMs()
}

func (wPMdb *WPMBoltDB) getWPMs() ([]*types.WPMRecord, error) {
	rr := make([]*types.WPMRecord, 0, 5)
	er := wPMdb.DB.View(func(tx *bolt.Tx) error {
		bucketname := []byte(wPMBucketName)
		bucket := tx.Bucket(bucketname)
		if bucket != nil {
			c := bucket.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				r := types.NewWPMRecord()
				err := json.Unmarshal(v, r)
				if err != nil {
					return err
				}
				r.ID, err = strconv.ParseUint(string(k), 10, 64)
				if err != nil {
					return err
				}
				rr = append(rr, r)
			}
		}
		return nil
	})
	return rr, er
}

// UpdateWPM updates the types.WPMRecord in the database.
// Param record [in-out] the record to be updated.
// if record is new then record.ID is updated as well.
// Returns the error.
func (wPMdb *WPMBoltDB) UpdateWPM(r *types.WPMRecord) error {
	return wPMdb.updateWPMBucket(r)
}

// RemoveWPM removes the types.WPMRecord from the database.
// Param id [in] the key of the record to be removed.
// Returns the error.
func (wPMdb *WPMBoltDB) RemoveWPM(id uint64) error {
	return wPMdb.DB.Update(func(tx *bolt.Tx) error {
		bucketname := []byte(wPMBucketName)
		bucket := tx.Bucket(bucketname)
		if bucket != nil {
			idbb := []byte(fmt.Sprintf("%d", id))
			col := bucket.Get(idbb)
			if col != nil {
				err := bucket.Delete(idbb)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// updates the types.WPMRecord in the database.
// Param record [in-out] the record to be updated
func (wPMdb *WPMBoltDB) updateWPMBucket(r *types.WPMRecord) error {
	return wPMdb.DB.Update(func(tx *bolt.Tx) error {
		bucketname := []byte(wPMBucketName)
		bucket, err := tx.CreateBucketIfNotExists(bucketname)
		if err == nil {
			if r.ID == 0 {
				id, err := bucket.NextSequence()
				if err == nil {
					r.ID = id
				}
			}
			if err == nil {
				bb, err := json.Marshal(r)
				if err == nil {
					idbb := []byte(fmt.Sprintf("%d", r.ID))
					err = bucket.Put(idbb, bb)
				}
			}
		}
		return err
	})
}

// initialize is only useful if you want to add the default records to the db.
// otherwise you don't need it to do anything.
func (wPMdb *WPMBoltDB) initialize() (err error) {
	// create the key wpm record of 5 wpm.
	record := types.NewWPMRecord()
	record.WPM = 5
	record.ID = keyID
	if err = wPMdb.updateWPMBucket(record); err != nil {
		return
	}
	// create the copy wpm record of 5 wpm.
	record.ID = copyID
	record.WPM = 15
	err = wPMdb.updateWPMBucket(record)
	return
}
