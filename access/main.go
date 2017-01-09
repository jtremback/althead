package access

import (
	"bytes"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/jtremback/althea/types"
)

var (
	Accounts  []byte = []byte("Accounts")
	Neighbors []byte = []byte("Neighbors")
)

type NilError struct {
	s string
}

func (e *NilError) Error() string {
	return e.s
}

func MakeBuckets(db *bolt.DB) error {
	var err error
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(Neighbors)
		_, err = tx.CreateBucketIfNotExists(Accounts)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func SetNeighbor(tx *bolt.Tx, neighbor *types.Neighbor) error {
	b, err := json.Marshal(neighbor)
	if err != nil {
		return err
	}

	return tx.Bucket(Neighbors).Put(neighbor.ControlPubkey, b)
}

func GetNeighbor(tx *bolt.Tx, key []byte) (*types.Neighbor, error) {
	b := tx.Bucket(Neighbors).Get(key)

	if bytes.Compare(b, []byte{}) == 0 {
		return nil, &NilError{"neighbor not found"}
	}

	neighbor := &types.Neighbor{}
	err := json.Unmarshal(b, neighbor)
	if err != nil {
		return nil, err
	}

	return neighbor, nil
}

func GetNeighbors(tx *bolt.Tx) ([]*types.Neighbor, error) {
	var err error
	neighbors := []*types.Neighbor{}

	err = tx.Bucket(Neighbors).ForEach(func(k, v []byte) error {
		neighbor := &types.Neighbor{}
		err = json.Unmarshal(v, neighbor)
		if err != nil {
			return err
		}

		neighbors = append(neighbors, neighbor)

		return nil
	})
	if err != nil {
		return nil, err
	}
	return neighbors, nil
}
