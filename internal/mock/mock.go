package mock

import (
	"bytes"
)

// imported types
type buf bytes.Buffer

type (
	myint  int
	mybool bool
)

// MockStr is a type alias for a string
type MockStr string

// MockPtr is a type alias for a pointer to a string
type MockPtr *string

// MockPtrImport is a pointer to an imported type
type MockPtrImport *bytes.Buffer

// mockIface is an interface used for testing parsing of interfacees.
type MockIface interface {
	// mockFn is a function used for testing interface methods.
	MockFn(string, int, bool) error
}

// MockStruct is a struct used for testing parsing of structs.
type MockStruct struct {
	// mockField is a field used for testing struct fields.
	MockField int `json:"mockField"`
	// mockField2 is a field used for testing struct fields.
	MockField2 string `json:"mockField2"`
	// array
	MockField3 []string `json:"mockField3"`
	// map
	MockField4 map[string]string `json:"mockField4"`
	// ptr
	MockField5 *string `json:"mockField5"`
	// interface!
	MockInterface MockIface `json:"myIface"`
}

type Impl struct {
	NewField string `json:"newField"`
}

func (i Impl) MockFn(string, int, bool) error {
	return nil
}

// MockSlice is a slice used for testing parsing of slices.
type MockSlice []MockPtrImport

// MockEmptyIfSlice is an empty slice used for testing parsing of slices.
type MockEmptyIfaceSlice []interface{}

// MockEmptyStructSlice
type MockEmptyStructSlice []struct{}

// MockStructSlice
type MockStructSlice []MockStruct

// MockImportedStructSlice
type MockImportedStructSlice []bytes.Buffer

// MockMap is a map used for testing parsing of maps.
type MockMap map[string]interface{}

// MockMapSlice is a map slice used for testing parsing of maps of slices.
type MockMapSlice map[string][]MockStr

// MockSlicePointer is a pointer slice used for testing parsing of slice pointers.
type MockSlicePointer []*int

// UNGENERATED //

// MockNonStringMap does not get generated because it has no string key.
type MockNonStringMap map[int]string
