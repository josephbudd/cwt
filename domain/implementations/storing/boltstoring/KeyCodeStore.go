package boltstoring

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/josephbudd/cwt/domain/types"
)

const keyCodeBucketName string = "keyCode"

// KeyCodeBoltDB is the bolt db for key codes.
type KeyCodeBoltDB struct {
	DB    *bolt.DB
	path  string
	perms os.FileMode
}

// NewKeyCodeBoltDB constructs a new KeyCodeBoltDB.
// Param db [in-out] is an open bolt database.
// Returns a pointer to the new KeyCodeBoltDB.
func NewKeyCodeBoltDB(db *bolt.DB, path string, perms os.FileMode) *KeyCodeBoltDB {
	return &KeyCodeBoltDB{
		DB:    db,
		path:  path,
		perms: perms,
	}
}

// KeyCodeBoltDB implements KeyCodeStorer
// which is defined in domain/types/records.go

// Open opens the database.
// Returns the error.
func (keyCodedb *KeyCodeBoltDB) Open() error {
	// the bolt db is already open
	return nil
}

// Close closes the database.
// Returns the error.
func (keyCodedb *KeyCodeBoltDB) Close() error {
	return keyCodedb.DB.Close()
}

// GetKeyCode retrieves the types.KeyCodeRecord from the db.
// Param id [in] is the record id.
// Returns the record and error.
func (keyCodedb *KeyCodeBoltDB) GetKeyCode(id uint64) (*types.KeyCodeRecord, error) {
	var r types.KeyCodeRecord
	ids := fmt.Sprintf("%d", id)
	er := keyCodedb.DB.View(func(tx *bolt.Tx) error {
		bucketname := []byte(keyCodeBucketName)
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

// GetKeyCodes retrieves all of the types.KeyCodeRecord from the db.
// If there are no types.KeyCodeRecords in the db then it calls keyCodedb.initialize().
// See keyCodedb.initialize().
// Returns the records and error.
func (keyCodedb *KeyCodeBoltDB) GetKeyCodes() ([]*types.KeyCodeRecord, error) {
	if rr, err := keyCodedb.getKeyCodes(); len(rr) > 0 && err == nil {
		return rr, err
	}
	keyCodedb.initialize()
	return keyCodedb.getKeyCodes()
}

func (keyCodedb *KeyCodeBoltDB) getKeyCodes() ([]*types.KeyCodeRecord, error) {
	rr := make([]*types.KeyCodeRecord, 0, 5)
	er := keyCodedb.DB.View(func(tx *bolt.Tx) error {
		bucketname := []byte(keyCodeBucketName)
		bucket := tx.Bucket(bucketname)
		if bucket != nil {
			c := bucket.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				r := types.NewKeyCodeRecord()
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

// UpdateKeyCode updates the types.KeyCodeRecord in the database.
// Param record [in-out] the record to be updated.
// if record is new then record.ID is updated as well.
// Returns the error.
func (keyCodedb *KeyCodeBoltDB) UpdateKeyCode(r *types.KeyCodeRecord) error {
	if r.ID > 0 && r.ID < FirstValidID {
		// not a real record so just ignore it.
		return nil
	}
	return keyCodedb.updateKeyCodeBucket(r)
}

// RemoveKeyCode removes the types.KeyCodeRecord from the database.
// Param id [in] the key of the record to be removed.
// Returns the error.
func (keyCodedb *KeyCodeBoltDB) RemoveKeyCode(id uint64) error {
	return keyCodedb.DB.Update(func(tx *bolt.Tx) error {
		bucketname := []byte(keyCodeBucketName)
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

// updates the types.KeyCodeRecord in the database.
// Param record [in-out] the record to be updated
func (keyCodedb *KeyCodeBoltDB) updateKeyCodeBucket(r *types.KeyCodeRecord) error {
	return keyCodedb.DB.Update(func(tx *bolt.Tx) error {
		bucketname := []byte(keyCodeBucketName)
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
