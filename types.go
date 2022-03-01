package toast

import (
	"bytes"
	"encoding/json"
	"fmt"

	"strings"
)

type Type interface {
	Node
	GetName() string
	GetTypeNames() []string
	SetTypeNames([]string)
	GetDocs() string
}

type Node interface {
	// Renders JSON representation of the node.
	Reflect() json.RawMessage
	// Renders Go type definitions.
	Go() string
	// Renders CUE definitions.
	CUE() string
}

type File struct {
	pkgName    string
	cuePkgName string

	Imports map[string]Import
	Code    []Type

	trans        []Transform
	copies       []*CopyIntoStruct
	eximports    []*ExcludeImport
	modimports   []*ModifyImport
	genEnumTrans []*GenEnumTypeTransform
	mkEnums      []*PromoteToEnumType

	debug bool
}

func (f *File) Reflect() json.RawMessage {
	imports := make([]string, 0, len(f.Imports))
	for _, imp := range f.Imports {
		imports = append(imports, string(imp.Reflect()))
	}
	code := make([]string, len(f.Code))
	for i, t := range f.Code {
		code[i] = string(t.Reflect())
	}
	raw := fmt.Sprintf(
		`{"package":"%s","imports":[%s],"code":[%s]}`,
		f.pkgName, strings.Join(imports, ","), strings.Join(code, ","),
	)
	if f.debug {
		fmt.Println(raw)
		return json.RawMessage{}
	}
	return json.RawMessage(raw)
}

type Import struct {
	Name    string `json:"name,omitempty"`
	Path    string `json:"path"`
	oldPath string
	used    bool
}

func (i *Import) Reflect() json.RawMessage {
	raw, _ := json.Marshal(i)
	return injectKind(string(raw), "import")
}

func (i *Import) GetOldPath() string {
	return i.oldPath
}

type PlainType struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Docs string `json:"-"`
}

func (p *PlainType) Reflect() json.RawMessage {
	raw, _ := json.Marshal(p)
	return injectKind(string(raw), "plain")
}

type ArrayType struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Length int    `json:"length,omitempty"`
	Docs   string `json:"-"`
}

func (a *ArrayType) Reflect() json.RawMessage {
	raw, _ := json.Marshal(a)
	return injectKind(string(raw), "array")
}

type MapType struct {
	Name      string `json:"name"`
	KeyType   string `json:"key_type"`
	ValueType string `json:"value_type"`
	Docs      string `json:"-"`
}

func (m *MapType) Reflect() json.RawMessage {
	raw, _ := json.Marshal(m)
	return injectKind(string(raw), "map")
}

type StructType struct {
	Name   string   `json:"name"`
	Fields []*Field `json:"fields"`
	Docs   string   `json:"-"`
}

func (s *StructType) Reflect() json.RawMessage {
	fields := make([]string, len(s.Fields))
	for i, f := range s.Fields {
		fields[i] = string(f.Reflect())
	}
	return json.RawMessage(
		fmt.Sprintf(
			`{"kind":"struct","name":"%s","fields":[%s]}`,
			s.Name, strings.Join(fields, ","),
		),
	)
}

type EnumType struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
	Docs   string   `json:"-"`
}

func (et *EnumType) Reflect() json.RawMessage {
	raw, _ := json.Marshal(et)
	return injectKind(string(raw), "enum")
}

func (p *PlainType) GetName() string  { return p.Name }
func (a *ArrayType) GetName() string  { return a.Name }
func (m *MapType) GetName() string    { return m.Name }
func (s *StructType) GetName() string { return s.Name }
func (et *EnumType) GetName() string  { return et.Name }

func (p *PlainType) GetTypeNames() []string  { return []string{p.Type} }
func (a *ArrayType) GetTypeNames() []string  { return []string{a.Type} }
func (m *MapType) GetTypeNames() []string    { return []string{m.KeyType, m.ValueType} }
func (s *StructType) GetTypeNames() []string { panic("not implemented") }
func (et *EnumType) GetTypeNames() []string  { return []string{et.Name} }

func (p *PlainType) SetTypeNames(tt []string)   { p.Type = tt[0] }
func (a *ArrayType) SetTypeNames(tt []string)   { a.Type = tt[0] }
func (m *MapType) SetTypeNames(tt []string)     { m.KeyType = tt[0]; m.ValueType = tt[1] }
func (s *StructType) SetTypeNames(tt []string)  { panic("not implemented") }
func (et *EnumType) SetTypeNames(typs []string) { et.Name = typs[0] }

func (p *PlainType) GetDocs() string  { return p.Docs }
func (a *ArrayType) GetDocs() string  { return a.Docs }
func (m *MapType) GetDocs() string    { return m.Docs }
func (s *StructType) GetDocs() string { return s.Docs }
func (et *EnumType) GetDocs() string  { return et.Docs }

type Field struct {
	Type
	Tags map[string][]string
}

func (f *Field) Reflect() json.RawMessage {
	raw := f.Type.Reflect()
	var tags []string
	for k, v := range f.Tags {
		tags = append(tags, fmt.Sprintf(`"%s":"%s"`, k, strings.Join(v, ",")))
	}
	return json.RawMessage(fmt.Sprintf(`%s,"tags":{%s}}`, raw[:len(raw)-1], strings.Join(tags, ",")))
}

func injectKind(raw string, kind string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{"kind":"%s",%s`, kind, raw[1:]))
}

func printJSON(raw json.RawMessage) {
	b := new(bytes.Buffer)
	json.Indent(b, raw, "", "  ")
	fmt.Println(b.String())
}
