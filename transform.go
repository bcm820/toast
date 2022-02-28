package toast

import "go/ast"

type Transform interface {
	isTransform()
}

type ExcludeImport struct {
	Match func(Import) bool
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

func (ei *ExcludeImport) isTransform()          {}
func (et *ExcludeType) isTransform()            {}
func (ef *ExcludeField) isTransform()           {}
func (ci *CopyIntoStruct) isTransform()         {}
func (mt *ModifyType) isTransform()             {}
func (mf *ModifyField) isTransform()            {}
func (gft *GenFieldTransform) isTransform()     {}
func (gett *GenEnumTypeTransform) isTransform() {}
func (ftet *PromoteToEnumType) isTransform()    {}
