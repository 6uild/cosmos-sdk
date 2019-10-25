package orm

import (
	"bytes"
	"sort"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMultiRef creates a MultiRef with any number of initial references
func NewMultiRef(refs ...[]byte) (*MultiRef, error) {
	m := new(MultiRef)
	for _, r := range refs {
		err := m.Add(r)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

// multiRefFromStrings is like NewMultiRef, but takes strings
// intended for test code.
func multiRefFromStrings(strs ...string) (*MultiRef, error) {
	refs := make([][]byte, len(strs))
	for i, s := range strs {
		refs[i] = []byte(s)
	}
	return NewMultiRef(refs...)
}

// Add inserts this reference in the multiref, sorted by order.
// Returns an error if already there
func (m *MultiRef) Add(ref []byte) error {
	i, found := m.findRef(ref)
	if found {
		return errors.Wrap(errors.ErrNotFound, "cannot add a ref twice")
	}
	// append to end
	if i == len(m.Refs) {
		m.Refs = append(m.Refs, ref)
		return nil
	}
	// or insert in the middle
	m.Refs = append(m.Refs, nil)
	copy(m.Refs[i+1:], m.Refs[i:])
	m.Refs[i] = ref
	return nil
}

// Remove removes this reference from the multiref.
// Returns an error if already there
func (m *MultiRef) Remove(ref []byte) error {
	i, found := m.findRef(ref)
	if !found {
		return errors.Wrap(errors.ErrNotFound, "cannot remove non-existent ref")
	}
	// splice it out
	m.Refs = append(m.Refs[:i], m.Refs[i+1:]...)
	return nil
}

// Sort will make sure everything is in order
func (m *MultiRef) Sort() {
	sort.Slice(m.Refs, func(i, j int) bool {
		return bytes.Compare(m.Refs[i], m.Refs[j]) == -1
	})
}

// returns (index, found) where found is true if
// the ref was in the set, index is where it is
// (or where it should be)
func (m *MultiRef) findRef(ref []byte) (int, bool) {
	for i, r := range m.Refs {
		switch bytes.Compare(ref, r) {
		case -1:
			return i, false
		case 0:
			return i, true
		}
	}
	// hit the end, must append
	return len(m.Refs), false
}

// Validate just returns an error if empty
func (m *MultiRef) Validate() error {
	if len(m.Refs) == 0 {
		return errors.Wrap(errors.ErrEmpty, "no references")
	}
	return nil
}

// Unmarshal amino decode from data.
func (m *MultiRef) Unmarshal(data []byte) error {
	return ModuleCdc.UnmarshalBinaryBare(data, m)
}

// Marshal amino encode into bytes.
func (m *MultiRef) Marshal() ([]byte, error) {
	return ModuleCdc.MarshalBinaryBare(m)
}

// Empty returns true when it contains no key reference elements.
func (m *MultiRef) Empty() bool {
	return m == nil || len(m.Refs) == 0
}
