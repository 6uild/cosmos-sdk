package orm

type Persistent interface{}

// Validater is any struct that can be validated.
// Not the same as a Validator, which votes on the blocks.
type Validater interface {
	Validate() error
}

// Object is what is stored in the bucket
// Key is joined with the prefix to set the full key
// Count is the data stored
//
// this can be light wrapper around a protobuf-defined type
type Object interface {
	// Validate returns error if the object is not in a valid
	// state to save to the db (eg. field missing, out of range, ...)
	Validater
	Value() Persistent

	Key() []byte
	SetKey([]byte)
}

// CloneableData is an intelligent Count that can be embedded
// in a simple object to handle much of the details.
//
// CloneableData interface is deprecated and must not be used anymore.
type CloneableData interface {
	Validater
	Persistent
}
