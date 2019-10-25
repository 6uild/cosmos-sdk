package orm

// ReadOnlyKVStore is a simple interface to query data.
type ReadOnlyKVStore interface {
	// Get returns nil iff key doesn't exist. Panics on nil key.
	Get(key []byte) ([]byte, error)

	// Has checks if a key exists. Panics on nil key.
	Has(key []byte) (bool, error)

	// Iterator over a domain of keys in ascending order. End is exclusive.
	// Start must be less than end, or the Iterator is invalid.
	// CONTRACT: No writes may happen within a domain while an iterator exists over it.
	Iterator(start, end []byte) (Iterator, error)

	// ReverseIterator over a domain of keys in descending order. End is exclusive.
	// Start must be less than end, or the Iterator is invalid.
	// CONTRACT: No writes may happen within a domain while an iterator exists over it.
	ReverseIterator(start, end []byte) (Iterator, error)
}

// SetDeleter is a minimal interface for writing,
// Unifying KVStore and Batch
type SetDeleter interface {
	Set(key, value []byte) error // CONTRACT: key, value readonly []byte
	Delete(key []byte) error     // CONTRACT: key readonly []byte
}

// KVStore is a simple interface to get/set data
//
// For simplicity, we require all backing stores to implement this
// interface. They *may* implement other methods as well, but
// at least these are required.
type KVStore interface {
	ReadOnlyKVStore
	SetDeleter
	// NewBatch returns a batch that can write multiple ops atomically
	NewBatch() Batch
}

// Batch can write multiple ops atomically to an underlying KVStore
type Batch interface {
	SetDeleter
	Write() error
}

/*
Iterator allows us to access a set of items within a range of
keys. These may all be preloaded, or loaded on demand.

  Usage:

  var itr Iterator = ...
  defer itr.Release()

  k, v, err := itr.Next()
  for err == nil {
	// ... do stuff with k, v
	k, v, err = itr.Next()
  }
  // ErrIteratorDone means we hit the end, otherwise this is a real error
  if !errors.ErrIteratorDone.Is(err) {
	  return err
  }
*/
type Iterator interface {
	// Next moves the iterator to the next sequential key in the database, as
	// defined by order of iteration.
	//
	// Returns (nil, nil, errors.ErrIteratorDone) if there is no more data
	Next() (key, value []byte, err error)

	// Release releases the Iterator, allowing it to do any needed cleanup.
	Release()
}
