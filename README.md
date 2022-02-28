# ToAST

ToAST is a thin wrapper for the Go standard library [ast](https://pkg.go.dev/go/ast) package that
provides a simple interface for working with Go types.

This package only supports parsing of type definitions. It was originally intended to be used for
generating type definitions for other languages -- particularly [CUE](https://cuelang.org), a
typesafe superset of JSON.

## Types

Each _go/ast_ type specification maps to one of four structs:
* `PlainType`: Basic types, type aliases, and pointers
* `ArrayType`: Array and slice types
* `MapType`: Maps
* `StructType`: Structs, which may contain fields that are themselves one of the four types.

There is also partial support for an `EnumType`, which is expressed by convention in Go as a type
declaration with a group of constants of the same type.

## Transforms

When parsing a file, ToAST can apply a number of transformations on matching objects:
* `ExcludeImport` excludes specific imports
* `ExcludeType` excludes specific types
* `ModifyType` mutates a specific type
* `ExcludeField` excludes a specific field in a `StructType`
* `ModifyField` mutates a specific field in a `StructType`
* `CopyIntoStruct` copies fields from a number of named `StructType`s into a target `StructType`,
  replacing a field
* `PromoteToEnumType` converts a `PlainType` into an `EnumType`
* `GenFieldTransform` takes a `StructType` and a `Field` and returns a `Transform` that can be
  matched on subsequent nodes in the file
* `GenEnumTypeTransform` takes a string and `ast.ValueSpec` (from _go/ast_) and returns a
  `PromoteToEnumType` transform that can be matched on a `PlainType` in the file

## Usage

First, load an `*ast.File`. For example:

```go
import "go/parser"

filePath := "path/to/file.go"
astFile, err := parser.ParseFile(token.NewFileSet(), filePath, nil, parser.ParseComments)
if err != nil {
	panic(err)
}
```

Then create a `*toast.File`:

```go
file := toast.NewFile(astFile,
  WithTransform(&toast.ExcludeImport{
    Match: func(i Import) bool {
      return i.Name == "foo"
    }
  }),
  WithTransform(&toast.ExcludeType{
    Match: func(t Type) bool {
      return t.Name == "bar"
    }
  }),
  ...
)
```
