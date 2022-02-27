package toast

type Transform interface {
	isTransform()
}

type ExcludeImport struct {
	Match func(Import) bool
}

type ExcludeType struct {
	Match func(Type) bool
}

type CopyIntoStruct struct {
	StructName     string
	FieldToReplace string
	FromStructs    map[string]struct{}
	with           []*Field
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

type AddFieldTransform struct {
	Generate func(*StructType, *Field) Transform
}

func (ei *ExcludeImport) isTransform()    {}
func (et *ExcludeType) isTransform()      {}
func (ef *ExcludeField) isTransform()     {}
func (ci *CopyIntoStruct) isTransform()   {}
func (mt *ModifyType) isTransform()       {}
func (mf *ModifyField) isTransform()      {}
func (a *AddFieldTransform) isTransform() {}
