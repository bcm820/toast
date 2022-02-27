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
			docs: "// a string\n",
			Name: "mystr",
			Type: "string",
		},
		&PlainType{
			docs: "// multi-line\n// comment\n",
			Name: "myinterface",
			Type: "interface{}",
		},
		&PlainType{
			docs: "/* another multi-line\n comment */\n",
			Name: "mystr",
			Type: "string",
		},
		&ArrayType{
			docs: "// a slice\n",
			Name: "myslice",
			Type: "int",
		},
		&ArrayType{
			docs:   "// a fixed-length array\n",
			Name:   "myarr",
			Type:   "int",
			Length: 5,
		},
		&MapType{
			docs:      "// a string map\n",
			Name:      "mymap",
			KeyType:   "string",
			ValueType: "int64",
		},
		&StructType{
			docs: "// a struct\n",
			Name: "mystruct",
			Fields: []*Field{
				{
					Type: &PlainType{
						docs: "field1",
						Name: "Field1",
						Type: "int32",
					},
					Tags: map[string][]string{
						"json": {"field1"},
					},
				},
				{
					Type: &ArrayType{
						docs: "field2",
						Name: "Field2",
						Type: "bool",
					},
					Tags: map[string][]string{
						"json": {"field2"},
					},
				},
				{
					Type: &MapType{
						docs:      "field3",
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
						docs: "field4",
						Name: "Field4",
						Fields: []*Field{
							{
								Type: &PlainType{
									docs: "nestedfield",
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
	},
}

func TestReflect(t *testing.T) {
	printJSON(file.Reflect())
}

func TestRender(t *testing.T) {
	fmt.Println(file.Go())
}
