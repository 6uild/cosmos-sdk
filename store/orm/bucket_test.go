package orm

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/errors"
	dbm "github.com/tendermint/tm-db"
)

type DemoCounter struct {
	Count int64 `json:"count"`
}

func (c DemoCounter) Validate() error {
	return nil
}

func TestBasicBucketOperations(t *testing.T) {
	db := Wrap(dbm.NewMemDB())
	cdc := codec.New()
	cdc.RegisterConcrete(DemoCounter{}, "test/DemoCounter", nil)

	// define some indexer
	var evenIndexer = demoIndexer(func(c DemoCounter) ([]byte, error) {
		if c.Count%2 == 0 {
			return itob(c.Count), nil
		}
		return nil, nil
	})

	var oddIndexer = demoIndexer(func(c DemoCounter) ([]byte, error) {
		if c.Count%2 == 1 {
			return itob(c.Count), nil
		}
		return nil, nil
	})

	bucket := NewBucket(cdc, "counter", &DemoCounter{}).
		WithIndex("even", evenIndexer, true).
		WithIndex("odd", oddIndexer, true)

	// when
	seq := bucket.Sequence("counter_id")
	id, err := seq.NextVal(db)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	obj := NewSimpleObj(id, &DemoCounter{1})
	err = bucket.Save(db, obj)
	// then
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// and when loaded by id
	loadedObj, err := bucket.Get(db, id)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if exp, got := obj, loadedObj; !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %v but got %v", exp, got)
	}
	// and when loaded by index
	counterValue := itob(1)
	loadedObjs, err := bucket.GetIndexed(db, "odd", counterValue)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if exp, got := 1, len(loadedObjs); exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	if exp, got := obj, loadedObjs[0]; !reflect.DeepEqual(exp, got) {
		t.Fatalf("expected %v but got %v", exp, got)
	}
	// and it is not in the even index
	loadedObjs, err = bucket.GetIndexed(db, "even", counterValue)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if exp, got := 0, len(loadedObjs); exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}

	// and when deleted
	err = bucket.Delete(db, id)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	// then
	loadedObj, err = bucket.Get(db, id)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if got := loadedObj; got != nil {
		t.Errorf("expected nil but got %v", got)
	}
	// and also removed from the index
	// and it is not in the even index
	loadedObjs, err = bucket.GetIndexed(db, "odd", counterValue)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
	if exp, got := 0, len(loadedObjs); exp != got {
		t.Fatalf("expected %v but got %v", exp, got)
	}

}

func demoIndexer(f func(DemoCounter) ([]byte, error)) Indexer {
	return func(obj Object) ([]byte, error) {
		if obj == nil {
			return nil, errors.Wrap(errors.ErrHuman, "cannot take index of nil")
		}
		c, ok := obj.Value().(*DemoCounter)
		if !ok {
			return nil, errors.Wrap(errors.ErrHuman, "can only take index of DemoCounter")
		}
		return f(*c)
	}
}

// itob integer to bytes
func itob(i int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(i))
	return bz
}
