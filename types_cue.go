package toast

import (
	"fmt"
	"strings"
)

func (f *File) CUE() string {
	var imports, code string
	for _, i := range f.Imports {
		imports += i.CUE()
	}
	if len(imports) > 0 {
		imports = fmt.Sprintf("import (\n%s)\n\n", imports)
	}
	for _, t := range f.Code {
		code += t.GetDocs() + t.CUE() + "\n"
	}
	src := []byte(fmt.Sprintf("package %s\n\n%s%s", f.cuePkgName, imports, code))
	return string(src)[:len(src)-1]
}

func (i *Import) CUE() string {
	if i.Name == "" {
		return fmt.Sprintf("  \"%s\"\n", i.Path)
	}
	return fmt.Sprintf("  %s \"%s\"\n", i.Name, i.Path)
}

func (p *PlainType) CUE() string {
	return fmt.Sprintf("#%s: %s\n", p.Name, fmtToCUE(p.Type))
}

func (a *ArrayType) CUE() string {
	if a.Type == "byte" {
		return fmt.Sprintf("#%s: bytes\n", a.Name)
	}
	return fmt.Sprintf("#%s: [...%s]\n", a.Name, fmtToCUE(a.Type))
}

func (m *MapType) CUE() string {
	keyTyp := fmtToCUE(m.KeyType)
	valTyp := fmtToCUE(m.ValueType)
	return fmt.Sprintf("#%s: [%s]: %s\n", m.Name, keyTyp, valTyp)
}

func (s *StructType) CUE() string {
	var fields string
	for _, f := range s.Fields {
		fields += f.CUE()
	}
	return fmt.Sprintf("#%s: {\n%s}\n", s.Name, fields)
}

func (et *EnumType) CUE() string {
	values := make([]string, 0, len(et.Values))
	for _, v := range et.Values {
		values = append(values, `"`+v+`"`)
	}
	str := "#" + et.Name + ": " + strings.Join(values, " | ") + "\n\n"
	for _, v := range et.Values {
		str += fmt.Sprintf("%s_%s: \"%s\"\n", et.Name, v, v)
	}
	return str
}

func (f *Field) CUE() string {
	if len(f.Tags) == 0 {
		return ""
	}
	jsonTags := f.Tags["json"]
	name := jsonTags[0]
	if len(jsonTags) > 1 && jsonTags[1] == "omitempty" {
		name += "?"
	}

	var str string
	switch ft := f.Type.(type) {
	case *PlainType:
		str = fmtToCUE(ft.Type)
	case *ArrayType:
		if ft.Type == "byte" {
			str = "bytes"
		} else {
			str = "[..." + fmtToCUE(ft.Type) + "]"
		}
	case *MapType:
		keyTyp := fmtToCUE(ft.KeyType)
		valTyp := fmtToCUE(ft.ValueType)
		str = fmt.Sprintf("[%s]: %s", keyTyp, valTyp)
	case *StructType:
		var fields string
		for _, f := range ft.Fields {
			fields += f.CUE()
		}
		str = fmt.Sprintf("{\n%s}", fields)
	}
	if docs := f.Type.GetDocs(); docs != "" {
		return fmt.Sprintf("%s%s: %s\n", docs, name, str)
	}
	return fmt.Sprintf("%s: %s\n", name, str)
}

var basicTypes = map[string]bool{
	"bool":      true,
	"string":    true,
	"int":       true,
	"int8":      true,
	"int16":     true,
	"int32":     true,
	"int64":     true,
	"uint":      true,
	"uint8":     true,
	"uint16":    true,
	"uint32":    true,
	"uint64":    true,
	"uintptr":   true,
	"byte":      true,
	"rune":      true,
	"float32":   true,
	"float64":   true,
	"complex64": true,
}

func fmtToCUE(typ string) string {
	typ = strings.Replace(typ, "*", "", 1)
	if typ == "interface{}" || typ == "error" {
		return "_"
	}
	if typ == "struct{}" {
		return "{}"
	}
	if strings.HasPrefix(typ, "[]") {
		typ = typ[2:]
		return fmt.Sprintf("[...%s]", fmtToCUE(typ))
	}
	if strings.Contains(typ, ".") {
		return strings.Replace(typ, ".", ".#", 1)
	}
	if !basicTypes[typ] && !strings.HasPrefix(typ, "#") {
		return "#" + typ
	}
	return typ
}
