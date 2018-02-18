package data

import (
	"reflect"
	"testing"
	"google.golang.org/appengine/datastore"
)

type OmitemptyElem struct {
	A string `datastore:"A,noindex,omitempty"`
	B string `datastore:"B,noindex,omitempty"`
}

type OmitemptyArrayElem struct {
	Array []OmitemptyElem `datastore:"Array,noindex"`
}

func TestDataSerialization(t *testing.T) {
	tests := []struct{
		In interface{}
		Expected interface{}
	}{{
		In: &OmitemptyArrayElem{
			Array: []OmitemptyElem{{A:"A0"}},
		},
	}, {
		In: &OmitemptyArrayElem{
			Array: []OmitemptyElem{{A:"A0"},{A:"A1"}},
		},
	}, {
		In: &OmitemptyArrayElem{
			Array: []OmitemptyElem{{A: "A0", B: "B0"}, {A: "A1", B: "B1"}},
		},
	}, {
		In: &OmitemptyArrayElem{
			Array: []OmitemptyElem{{A: "A0", B: "B0"}, {B: "B1"}},
		},
	}, {
		In: &OmitemptyArrayElem{
			Array: []OmitemptyElem{{A: "A0"}, {A: "A1", B: "B1"}},
		},
	}}

	for _, v := range tests {
		gotSave, err := datastore.SaveStruct(v.In)
		if err != nil {
			t.Errorf("got unexpected error during SaveStruct: %v", err)
			continue
		}
		out := reflect.New(reflect.ValueOf(v.In).Elem().Type()).Interface()
		if err := datastore.LoadStruct(out, gotSave); err != nil {
			t.Errorf("got unexpected error during LoadStruct: %v", err)
			continue
		}

		if !reflect.DeepEqual(v.In, out) {
			t.Errorf("got unexpected datastore roundtrip.\n\tgot %+v\n\texpected %+v", out, v.In)
		}
	}
}
