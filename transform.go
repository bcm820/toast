package toast

import "go/ast"

type Transform interface {
	isTransform()
}

type ExcludeImport struct {
	Match func(Import) bool
}

type ModifyImport struct {
	Apply func(Import) Import
}

type ExcludeType struct {
	Match func(Type) bool
}

type ModifyType struct {
	Apply func(Type) Type
}

type ExcludeField struct {
	Match func(*Field) bool
}

type ModifyField struct {
	Apply func(*Field) *Field
}

type CopyIntoStruct struct {
	StructName     string
	FieldToReplace string
	FromStructs    map[string]struct{}
	with           []*Field
}

type GenFieldTransform struct {
	Generate func(*StructType, *Field) Transform
}

type GenEnumTypeTransform struct {
	Generate func(string, *ast.ValueSpec) *PromoteToEnumType
}

type PromoteToEnumType struct {
	Apply func(*PlainType) *EnumType
}

func (*ExcludeImport) isTransform()        {}
func (*ModifyImport) isTransform()         {}
func (*ExcludeType) isTransform()          {}
func (*ExcludeField) isTransform()         {}
func (*CopyIntoStruct) isTransform()       {}
func (*ModifyType) isTransform()           {}
func (*ModifyField) isTransform()          {}
func (*GenFieldTransform) isTransform()    {}
func (*GenEnumTypeTransform) isTransform() {}
func (*PromoteToEnumType) isTransform()    {}
