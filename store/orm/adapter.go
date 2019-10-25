package orm

import "github.com/cosmos/cosmos-sdk/store/types"

var _ KVStore = adapter{}

// methods used from the sdk KVStore and implementd in the memory store
type sdkKVStore interface {
	Get(key []byte) []byte
	Has(key []byte) bool
	Set(key, value []byte)
	Delete(key []byte)
	Iterator(start, end []byte) types.Iterator
	ReverseIterator(start, end []byte) types.Iterator
}

// adapter bridges the old panic affine store to a more Go style api with some new functionality.
type adapter struct {
	db sdkKVStore
}

func Wrap(db sdkKVStore) adapter {
	return adapter{db: db}
}
func (a adapter) Get(key []byte) ([]byte, error) {
	// todo: prevent panics
	return a.db.Get(key), nil
}

func (a adapter) Has(key []byte) (bool, error) {
	// todo: prevent panics
	return a.db.Has(key), nil
}

func (a adapter) Set(key, value []byte) error {
	// todo: prevent panics
	a.db.Set(key, value)
	return nil
}

func (a adapter) Delete(key []byte) error {
	// todo: prevent panics
	a.db.Delete(key)
	return nil
}

func (a adapter) Iterator(start, end []byte) (Iterator, error) {
	panic("implement me")
}
func (a adapter) ReverseIterator(start, end []byte) (Iterator, error) {
	panic("implement me")
}

func (a adapter) NewBatch() Batch {
	panic("implement me")
}
