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
	// Renders JSON
	Reflect() json.RawMessage
	// Renders Go
	Go() string
	// Renders CUE
	CUE() string
}

type File struct {
	Package string
	Imports map[string]Import
	Code    []Type

	transforms []Transform
	copies     []*CopyIntoStruct
	eximports  []*ExcludeImport

	debug bool
}

func (f *File) Reflect() json.RawMessage {
	imports := make([]string, len(f.Imports))
	for _, imp := range f.Imports {
		imports = append(imports, string(imp.Reflect()))
	}
	code := make([]string, len(f.Code))
	for i, t := range f.Code {
		code[i] = string(t.Reflect())
	}
	raw := fmt.Sprintf(
		`{"package":"%s","imports":[%s],"code":[%s]}`,
		f.Package, strings.Join(imports, ","), strings.Join(code, ","),
	)
	if f.debug {
		fmt.Println(raw)
		return json.RawMessage{}
	}
	return json.RawMessage(raw)
}

type Import struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path"`
	used bool
}

func (i *Import) Reflect() json.RawMessage {
	raw, _ := json.Marshal(i)
	return injectKind(string(raw), "import")
}

type PlainType struct {
	Name string `json:"name"`
	Type string `json:"type"`
	docs string
}

func (p *PlainType) Reflect() json.RawMessage {
	raw, _ := json.Marshal(p)
	return injectKind(string(raw), "plain")
}

type ArrayType struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Length int    `json:"length,omitempty"`
	docs   string
}

func (a *ArrayType) Reflect() json.RawMessage {
	raw, _ := json.Marshal(a)
	return injectKind(string(raw), "array")
}

type MapType struct {
	Name      string `json:"name"`
	KeyType   string `json:"key_type"`
	ValueType string `json:"value_type"`
	docs      string
}

func (m *MapType) Reflect() json.RawMessage {
	raw, _ := json.Marshal(m)
	return injectKind(string(raw), "map")
}

type StructType struct {
	Name   string   `json:"name"`
	Fields []*Field `json:"fields"`
	docs   string
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

func (p *PlainType) GetName() string  { return p.Name }
func (a *ArrayType) GetName() string  { return a.Name }
func (m *MapType) GetName() string    { return m.Name }
func (s *StructType) GetName() string { return s.Name }

func (p *PlainType) GetTypeNames() []string  { return []string{p.Type} }
func (a *ArrayType) GetTypeNames() []string  { return []string{a.Type} }
func (m *MapType) GetTypeNames() []string    { return []string{m.KeyType, m.ValueType} }
func (s *StructType) GetTypeNames() []string { panic("not implemented") }

func (p *PlainType) SetTypeNames(tt []string)  { p.Type = tt[0] }
func (a *ArrayType) SetTypeNames(tt []string)  { a.Type = tt[0] }
func (m *MapType) SetTypeNames(tt []string)    { m.KeyType = tt[0]; m.ValueType = tt[1] }
func (s *StructType) SetTypeNames(tt []string) { panic("not implemented") }

func (p *PlainType) GetDocs() string  { return p.docs }
func (a *ArrayType) GetDocs() string  { return a.docs }
func (m *MapType) GetDocs() string    { return m.docs }
func (s *StructType) GetDocs() string { return s.docs }

type Field struct {
	Type
	Tags map[string][]string
}

func (f *Field) Reflect() json.RawMessage {
	raw := f.Type.Reflect()
	tags := "{"
	for k, v := range f.Tags {
		tags += fmt.Sprintf(`"%s":"%s",`, k, strings.Join(v, ","))
	}
	tags += "}"
	return json.RawMessage(fmt.Sprintf(`%s,"tags":%s}`, raw[:len(raw)-1], tags))
}

func injectKind(raw string, kind string) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`%s"kind":"%s",%s`, raw[0:1], kind, raw[1:]))
}

func printJSON(raw json.RawMessage) {
	b := new(bytes.Buffer)
	json.Indent(b, raw, "", "  ")
	fmt.Println(b.String())
}
