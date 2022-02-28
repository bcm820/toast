package toast

import (
	"fmt"
	"testing"
)

var file = File{
	Package: "mock",
	Imports: map[string]Import{
		"fmt":    {Path: "fmt"},
		"loader": {Path: "github.com/bcmendoza/envoy-cue/pkg/loader"},
		"m":      {Name: "m", Path: "github.com/bcmendoza/envoy-cue/pkg/toast/mock"},
	},
	Code: []Type{
		&PlainType{
			Name: "myint",
			Type: "int",
		},
		&PlainType{
			Docs: "// a string\n",
			Name: "mystr",
			Type: "string",
		},
		&PlainType{
			Docs: "// multi-line\n// comment\n",
			Name: "myinterface",
			Type: "interface{}",
		},
		&PlainType{
			Docs: "/* another multi-line\n comment */\n",
			Name: "mystr",
			Type: "string",
		},
		&ArrayType{
			Docs: "// a slice\n",
			Name: "myslice",
			Type: "int",
		},
		&ArrayType{
			Docs:   "// a fixed-length array\n",
			Name:   "myarr",
			Type:   "int",
			Length: 5,
		},
		&MapType{
			Docs:      "// a string map\n",
			Name:      "mymap",
			KeyType:   "string",
			ValueType: "int64",
		},
		&StructType{
			Docs: "// a struct\n",
			Name: "mystruct",
			Fields: []*Field{
				{
					Type: &PlainType{
						Docs: "// field1\n",
						Name: "Field1",
						Type: "int32",
					},
					Tags: map[string][]string{
						"json": {"field1"},
					},
				},
				{
					Type: &ArrayType{
						Docs: "// field2\n",
						Name: "Field2",
						Type: "bool",
					},
					Tags: map[string][]string{
						"json": {"field2"},
					},
				},
				{
					Type: &MapType{
						Docs:      "// field3\n",
						Name:      "Field3",
						KeyType:   "int",
						ValueType: "struct{}",
					},
					Tags: map[string][]string{
						"json": {"field3"},
					},
				},
				{
					Type: &StructType{
						Docs: "// field4\n",
						Name: "Field4",
						Fields: []*Field{
							{
								Type: &PlainType{
									Docs: "nestedfield",
									Name: "NestedField",
									Type: "int64",
								},
								Tags: map[string][]string{
									"json": {"nestedfield"},
								},
							},
						},
					},
					Tags: map[string][]string{
						"json": {"field4"},
					},
				},
			},
		},
		&EnumType{
			Name:   "MyEnum",
			Values: []string{"MyEnum_A", "MyEnum_B", "MyEnum_C"},
		},
	},
}

func TestReflect(t *testing.T) {
	printJSON(file.Reflect())
}

func TestGo(t *testing.T) {
	fmt.Println(file.Go())
}

func TestCUE(t *testing.T) {
	fmt.Println(file.CUE())
}
