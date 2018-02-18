# repro-datastore-array-serialization-error

Google Go Datastore serialization library fails to serialize omitempty elements in an array correctly.

Here's a simple reproduction case for the problem.

We create a struct with datastore serialization tags:

    type OmitemptyElem struct {
    	A string `datastore:"A,noindex,omitempty"`
    	B string `datastore:"B,noindex,omitempty"`
    }

    type OmitemptyArrayElem struct {
    	Array []OmitemptyElem `datastore:"Array,noindex"`
    }

We then define an instance:

    &OmitemptyArrayElem{
    	Array: []OmitemptyElem{{A: "A0"}, {A: "A1", B: "B1"}},
    }

When we roundtrip this through Datastore's SaveStruct/LoadStruct, we get:

    &OmitemptyArrayElem{
    	Array: []OmitemptyElem{{A: "A0", B: "B1"}, {A: "A1"}},
    }

To reproduce, set up an environment with the (current as of 2/17/2018) downloadable appengine SDK (https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_darwin_amd64-1.9.62.zip) and a GOPATH, and run:

    $ goapp get google.golang.org/appengine/datastore
    $ goapp test .

Which produces:

    2018/02/17 16:55:40 appengine: not running under devappserver2; using some default configuration
    --- FAIL: TestDataSerialization (0.00s)
    	data_test.go:57: got unexpected datastore roundtrip.
    			got &{[{A0 B1} {A1 }]}
    			expected &{[{A0 } {A1 B1}]}
    FAIL
    FAIL	_/Users/morgan/repro	0.011s
