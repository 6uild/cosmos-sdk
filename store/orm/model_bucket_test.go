package orm

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/errors"
	dbm "github.com/tendermint/tm-db"
)

func TestModelBucket(t *testing.T) {
	db := Wrap(dbm.NewMemDB())

	cdc := codec.New()
	cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)

	b := NewModelBucket(cdc, "cnts", &DemoCounter{})

	if _, err := b.Put(db, []byte("c1"), &DemoCounter{Count: 1}); err != nil {
		t.Fatalf("cannot save DemoCounter instance: %s", err)
	}

	var c1 DemoCounter
	if err := b.One(db, []byte("c1"), &c1); err != nil {
		t.Fatalf("cannot get c1 DemoCounter: %s", err)
	}
	if c1.Count != 1 {
		t.Fatalf("unexpected DemoCounter state: %d", c1)
	}

	if err := b.Delete(db, []byte("c1")); err != nil {
		t.Fatalf("cannot delete c1 DemoCounter: %s", err)
	}
	if err := b.Delete(db, []byte("unknown")); !errors.ErrNotFound.Is(err) {
		t.Fatalf("unexpected error when deleting unexisting instance: %s", err)
	}
	if err := b.One(db, []byte("c1"), &c1); !errors.ErrNotFound.Is(err) {
		t.Fatalf("unexpected error for an unknown model get: %s", err)
	}
}

func TestModelBucketPutSequence(t *testing.T) {
	db := Wrap(dbm.NewMemDB())

	cdc := codec.New()
	cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)
	b := NewModelBucket(cdc, "cnts", &DemoCounter{})

	// Using a nil key should cause the sequence ID to be used.
	key, err := b.Put(db, nil, &DemoCounter{Count: 111})
	if err != nil {
		t.Fatalf("cannot save DemoCounter instance: %s", err)
	}
	if !bytes.Equal(key, itob(1)) {
		t.Fatalf("first sequence key should be 1, instead got %d", key)
	}

	// Inserting an entity with a key provided must not modify the ID
	// generation DemoCounter.
	if _, err := b.Put(db, []byte("mycnt"), &DemoCounter{Count: 12345}); err != nil {
		t.Fatalf("cannot save DemoCounter instance: %s", err)
	}

	key, err = b.Put(db, nil, &DemoCounter{Count: 222})
	if err != nil {
		t.Fatalf("cannot save DemoCounter instance: %s", err)
	}
	if !bytes.Equal(key, itob(2)) {
		t.Fatalf("second sequence key should be 2, instead got %d", key)
	}

	var c1 DemoCounter
	if err := b.One(db, itob(1), &c1); err != nil {
		t.Fatalf("cannot get first DemoCounter: %s", err)
	}
	if c1.Count != 111 {
		t.Fatalf("unexpected DemoCounter state: %d", c1)
	}

	var c2 DemoCounter
	if err := b.One(db, itob(2), &c2); err != nil {
		t.Fatalf("cannot get first DemoCounter: %s", err)
	}
	if c2.Count != 222 {
		t.Fatalf("unexpected DemoCounter state: %d", c2)
	}
}

func TestModelBucketByIndex(t *testing.T) {
	cases := map[string]struct {
		IndexName  string
		QueryKey   string
		DestFn     func() ModelSlicePtr
		WantErr    *errors.Error
		WantResPtr []*DemoCounter
		WantRes    []DemoCounter
		WantKeys   [][]byte
	}{
		"find none": {
			IndexName:  "value",
			QueryKey:   "124089710947120",
			WantErr:    nil,
			WantResPtr: nil,
			WantRes:    nil,
			WantKeys:   nil,
		},
		"find one": {
			IndexName: "value",
			QueryKey:  "1",
			WantErr:   nil,
			WantResPtr: []*DemoCounter{
				{Count: 1001},
			},
			WantRes: []DemoCounter{
				{Count: 1001},
			},
			WantKeys: [][]byte{
				itob(1),
			},
		},
		"find two": {
			IndexName: "value",
			QueryKey:  "4",
			WantErr:   nil,
			WantResPtr: []*DemoCounter{
				{Count: 4001},
				{Count: 4002},
			},
			WantRes: []DemoCounter{
				{Count: 4001},
				{Count: 4002},
			},
			WantKeys: [][]byte{
				itob(3),
				itob(4),
			},
		},
		"non existing index name": {
			IndexName: "xyz",
			WantErr:   errors.ErrInvalidIndex,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			db := Wrap(dbm.NewMemDB())
			cdc := codec.New()
			cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)

			indexByBigValue := func(obj Object) ([]byte, error) {
				c, ok := obj.Value().(*DemoCounter)
				if !ok {
					return nil, errors.Wrapf(errors.ErrType, "%T", obj.Value())
				}
				// Index by the value, ignoring anything below 1k.
				raw := strconv.FormatInt(c.Count/1000, 10)
				return []byte(raw), nil
			}
			b := NewModelBucket(cdc, "cnts", &DemoCounter{}, WithIndex("value", indexByBigValue, false))

			if _, err := b.Put(db, nil, &DemoCounter{Count: 1001}); err != nil {
				t.Fatalf("cannot save DemoCounter instance: %s", err)
			}
			if _, err := b.Put(db, nil, &DemoCounter{Count: 2001}); err != nil {
				t.Fatalf("cannot save DemoCounter instance: %s", err)
			}
			if _, err := b.Put(db, nil, &DemoCounter{Count: 4001}); err != nil {
				t.Fatalf("cannot save DemoCounter instance: %s", err)
			}
			if _, err := b.Put(db, nil, &DemoCounter{Count: 4002}); err != nil {
				t.Fatalf("cannot save DemoCounter instance: %s", err)
			}

			var dest []DemoCounter
			keys, err := b.ByIndex(db, tc.IndexName, []byte(tc.QueryKey), &dest)
			if !tc.WantErr.Is(err) {
				t.Fatalf("unexpected error: %s", err)
			}
			if exp, got := tc.WantKeys, keys; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected %v but got %v", exp, got)
			}
			if exp, got := tc.WantRes, dest; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected %v but got %v", exp, got)
			}

			var destPtr []*DemoCounter
			keys, err = b.ByIndex(db, tc.IndexName, []byte(tc.QueryKey), &destPtr)
			if !tc.WantErr.Is(err) {
				t.Fatalf("unexpected error: %s", err)
			}
			if exp, got := tc.WantKeys, keys; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected %v but got %v", exp, got)
			}
			if exp, got := tc.WantResPtr, destPtr; !reflect.DeepEqual(exp, got) {
				t.Errorf("expected %v but got %v", exp, got)
			}
		})
	}
}

func TestModelBucketPutWrongModelType(t *testing.T) {
	db := Wrap(dbm.NewMemDB())
	cdc := codec.New()
	cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)

	b := NewModelBucket(cdc, "cnts", &DemoCounter{})

	if _, err := b.Put(db, nil, &MultiRef{Refs: [][]byte{[]byte("foo")}}); !errors.ErrType.Is(err) {
		t.Fatalf("unexpected error when trying to store wrong model type value: %s", err)
	}
}

func TestModelBucketOneWrongModelType(t *testing.T) {
	db := Wrap(dbm.NewMemDB())
	cdc := codec.New()
	cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)

	b := NewModelBucket(cdc, "cnts", &DemoCounter{})

	if _, err := b.Put(db, []byte("DemoCounter"), &DemoCounter{Count: 1}); err != nil {
		t.Fatalf("cannot save DemoCounter instance: %s", err)
	}

	var ref MultiRef
	if err := b.One(db, []byte("DemoCounter"), &ref); !errors.ErrType.Is(err) {
		t.Fatalf("unexpected error when trying to get wrong model type value: %s", err)
	}
}

func TestModelBucketByIndexWrongModelType(t *testing.T) {
	db := Wrap(dbm.NewMemDB())
	cdc := codec.New()
	cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)

	b := NewModelBucket(cdc, "cnts", &DemoCounter{},
		WithIndex("x", func(o Object) ([]byte, error) { return []byte("x"), nil }, false))

	if _, err := b.Put(db, []byte("DemoCounter"), &DemoCounter{Count: 1}); err != nil {
		t.Fatalf("cannot save DemoCounter instance: %s", err)
	}

	var refs []MultiRef
	if _, err := b.ByIndex(db, "x", []byte("x"), &refs); !errors.ErrType.Is(err) {
		t.Fatalf("unexpected error when trying to find wrong model type value: %s: %v", err, refs)
	}

	var refsPtr []*MultiRef
	if _, err := b.ByIndex(db, "x", []byte("x"), &refsPtr); !errors.ErrType.Is(err) {
		t.Fatalf("unexpected error when trying to find wrong model type value: %s: %v", err, refs)
	}

	var refsPtrPtr []**MultiRef
	if _, err := b.ByIndex(db, "x", []byte("x"), &refsPtrPtr); !errors.ErrType.Is(err) {
		t.Fatalf("unexpected error when trying to find wrong model type value: %s: %v", err, refs)
	}
}

func TestModelBucketHas(t *testing.T) {
	db := Wrap(dbm.NewMemDB())
	cdc := codec.New()
	cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)

	b := NewModelBucket(cdc, "cnts", &DemoCounter{})

	if _, err := b.Put(db, []byte("DemoCounter"), &DemoCounter{Count: 1}); err != nil {
		t.Fatalf("cannot save DemoCounter instance: %s", err)
	}

	if err := b.Has(db, []byte("DemoCounter")); err != nil {
		t.Fatalf("an existing entity must return no error: %s", err)
	}

	if err := b.Has(db, nil); !errors.ErrNotFound.Is(err) {
		t.Fatalf("a nil key must return ErrNotFound: %s", err)
	}

	if err := b.Has(db, []byte("does-not-exist")); !errors.ErrNotFound.Is(err) {
		t.Fatalf("a non exists entity must return ErrNotFound: %s", err)
	}
}
