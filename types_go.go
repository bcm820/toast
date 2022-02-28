package toast

import (
	"fmt"
	"go/format"
	"strings"

	"golang.org/x/tools/imports"
)

func (f *File) Go() string {
	var imports, code string
	for _, i := range f.Imports {
		imports += i.Go()
	}
	if len(imports) > 0 {
		imports = fmt.Sprintf("import (\n%s)\n\n", imports)
	}
	for _, t := range f.Code {
		code += t.GetDocs() + t.Go() + "\n"
	}
	src := []byte(fmt.Sprintf("package %s\n\n%s%s", f.Package, imports, code))
	if f.debug {
		return string(src)
	}
	return string(goFormat(src))
}

func (i *Import) Go() string {
	if i.Name == "" {
		return fmt.Sprintf("  \"%s\"\n", i.Path)
	}
	return fmt.Sprintf("  %s \"%s\"\n", i.Name, i.Path)
}

func (p *PlainType) Go() string {
	return fmt.Sprintf("type %s %s\n", p.Name, p.Type)
}

func (a *ArrayType) Go() string {
	brackets := "[]"
	if a.Length > 0 {
		brackets = fmt.Sprintf("[%d]", a.Length)
	}
	return fmt.Sprintf("type %s %s%s\n", a.Name, brackets, a.Type)
}

func (m *MapType) Go() string {
	return fmt.Sprintf("type %s map[%s]%s\n", m.Name, m.KeyType, m.ValueType)
}

func (s *StructType) Go() string {
	var fields string
	for _, f := range s.Fields {
		fields += f.Go()
	}
	return fmt.Sprintf("type %s struct {\n%s}\n", s.Name, fields)
}

func (et *EnumType) Go() string {
	str := fmt.Sprintf("type %s string\n\nconst (\n", et.Name)
	for _, v := range et.Values {
		str += fmt.Sprintf("  %s_%s = \"%s\"\n", et.Name, v, v)
	}
	str += ")\n"
	return str
}

func (f *Field) Go() string {
	tags := make([]string, len(f.Tags))
	for k, v := range f.Tags {
		tags = append(tags, fmt.Sprintf("%s:\"%s\"", k, strings.Join(v, ",")))
	}
	var tag string
	if len(f.Tags) > 0 {
		tag = "`" + strings.TrimSpace(strings.Join(tags, " ")) + "`"
	}
	str := f.Type.Go()
	str = strings.Replace(str[:len(str)-1], "type ", "", 1)
	if docs := f.Type.GetDocs(); docs != "" {
		return fmt.Sprintf("%s%s %s\n", docs, str, tag)
	}
	return fmt.Sprintf("%s %s\n", str, tag)
}

func goFormat(src []byte) []byte {
	formatted, err := imports.Process("", src, nil)
	if err != nil {
		panic(err)
	}
	formatted, err = format.Source(formatted)
	if err != nil {
		panic(err)
	}
	return formatted
}
