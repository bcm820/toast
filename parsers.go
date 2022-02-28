package toast

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
)

func FileFromAST(file *ast.File, opts ...Option) *File {

	f := &File{
		pkgName: file.Name.Name,
		Imports: make(map[string]Import),
	}

	for _, opt := range opts {
		opt(f)
	}

	for _, fileDecl := range file.Decls {
		switch decl := fileDecl.(type) {
		case *ast.GenDecl:
			docs := DocsFromCommentGroup(decl.Doc)
		SPEC_LOOP:
			for _, declSpec := range decl.Specs {
				switch ts := declSpec.(type) {
				case *ast.ImportSpec:
					for _, ei := range f.eximports {
						if ei.Match(ImportFromSpec(ts)) {
							continue SPEC_LOOP
						}
					}
					imp := ImportFromSpec(ts)
					impName := imp.Name
					if imp.Name == "" {
						impName = imp.Path[strings.LastIndex(imp.Path, "/")+1:]
					}
					f.Imports[impName] = imp
				case *ast.TypeSpec:
					if t := ParseExpr([]*ast.Ident{ts.Name}, docs, ts.Type); t != nil {
						f.Code = append(f.Code, t)
						for _, transform := range f.trans {
							if ok := evalTransform(transform, t, f); !ok {
								continue SPEC_LOOP
							}
						}
					}
				case *ast.ValueSpec:
					for _, gen := range f.genEnumTrans {
						if t := gen.Generate(docs, ts); t != nil {
							f.mkEnums = append(f.mkEnums, t)
						}
					}
				}
			}

		case *ast.FuncDecl:
		}
	}

COPIES_LOOP:
	for _, ci := range f.copies {
		structIdx := -1
	CODE_LOOP:
		for i, t := range f.Code {
			if t.GetName() == ci.StructName {
				structIdx = i
				break CODE_LOOP
			}
		}
		if structIdx == -1 {
			break COPIES_LOOP
		}
		if st, ok := f.Code[structIdx].(*StructType); ok {
		FIELD_LOOP:
			for i, field := range st.Fields {
				if field.GetName() == ci.FieldToReplace {
					var fields []*Field
					fields = append(fields, st.Fields[:i]...)
					fields = append(fields, ci.with...)
					fields = append(fields, st.Fields[i+1:]...)
					st.Fields = fields
					f.Code[structIdx] = st
					break FIELD_LOOP
				}
			}
		}
	}

	for i, t := range f.Code {
		switch tt := t.(type) {
		case *StructType:
			for _, field := range tt.Fields {
				for _, typ := range field.GetTypeNames() {
					if dot := strings.Index(typ, "."); dot > -1 {
						impName := strings.Replace(typ[:dot], "*", "", -1)
						if imp, ok := f.Imports[impName]; ok {
							imp.used = true
							f.Imports[impName] = imp
						}
					}
				}
			}
		default:
			pt, ok := tt.(*PlainType)
			if ok {
				for _, mkEnum := range f.mkEnums {
					if et := mkEnum.Apply(pt); et != nil {
						f.Code[i] = et
					}
				}
			}
			for _, typ := range t.GetTypeNames() {
				if dot := strings.Index(typ, "."); dot > -1 {
					impName := strings.Replace(typ[:dot], "*", "", -1)
					if imp, ok := f.Imports[impName]; ok {
						imp.used = true
						f.Imports[impName] = imp
					}
				}
			}
		}
	}

	for k, v := range f.Imports {
		if !v.used {
			delete(f.Imports, k)
		}
	}

	return f
}

func evalTransform(transform Transform, t Type, f *File) bool {
	switch tt := transform.(type) {
	case *GenFieldTransform:
		if st, ok := t.(*StructType); ok {
			for _, field := range st.Fields {
				if gen := tt.Generate(st, field); gen != nil {
					evalTransform(gen, t, f)
					switch gt := gen.(type) {
					case *CopyIntoStruct:
						f.copies = append(f.copies, gt)
						f.trans = append(f.trans, gen)
					default:
						f.trans = append(f.trans, gen)
					}
				}
			}
		}
	case *ExcludeType:
		if tt.Match(t) {
			f.Code = f.Code[:len(f.Code)-1]
			return false
		}
		if st, ok := t.(*StructType); ok {
			var fields []*Field
			for _, field := range st.Fields {
				if !tt.Match(field) {
					fields = append(fields, field)
				}
			}
			st.Fields = fields
		}
	case *ExcludeField:
		if st, ok := t.(*StructType); ok {
			var fields []*Field
			for _, field := range st.Fields {
				if !tt.Match(field) {
					fields = append(fields, field)
				}
			}
			st.Fields = fields
		}
	case *CopyIntoStruct:
		if st, ok := t.(*StructType); ok {
			if _, ok := tt.FromStructs[st.GetName()]; ok {
				tt.with = append(tt.with, st.Fields...)
				f.Code = f.Code[:len(f.Code)-1]
				return false
			}
		}
	case *ModifyType:
		switch mpt := t.(type) {
		case *StructType:
			for i, field := range mpt.Fields {
				mpt.Fields[i].Type = tt.Apply(field.Type)
			}
			f.Code[len(f.Code)-1] = mpt
		default:
			f.Code[len(f.Code)-1] = tt.Apply(mpt)
		}
	case *ModifyField:
		if st, ok := t.(*StructType); ok {
			for i, field := range st.Fields {
				st.Fields[i] = tt.Apply(field)
			}
			f.Code[len(f.Code)-1] = st
		}
	}
	return true
}

func ImportFromSpec(spec *ast.ImportSpec) Import {
	path := strings.Replace(spec.Path.Value, `"`, "", -1)
	imp := Import{Path: path}
	if spec.Name != nil {
		imp.Name = spec.Name.Name
	}
	return imp
}

func DocsFromCommentGroup(cg *ast.CommentGroup) string {
	if cg == nil {
		return ""
	}
	var docs []string
	for _, c := range cg.List {
		docs = append(docs, strings.TrimSpace(c.Text))
	}
	if len(docs) == 0 {
		return ""
	}
	return strings.Join(docs, "\n") + "\n"
}

func ParseExpr(names []*ast.Ident, docs string, expr ast.Expr) Type {
	var name string
	if len(names) > 0 {
		name = names[0].Name
	}
	switch expr := expr.(type) {
	case *ast.Ident:
		return PlainTypeFromIdent(name, docs, expr)
	case *ast.SelectorExpr:
		return PlainTypeFromSelectorExpr(name, docs, expr)
	case *ast.StarExpr:
		return PlainTypeFromStarExpr(name, docs, expr)
	case *ast.ArrayType:
		return ArrayTypeFromSpec(name, docs, expr)
	case *ast.MapType:
		return MapTypeFromSpec(name, docs, expr)
	case *ast.StructType:
		return StructTypeFromSpec(name, docs, expr)
	case *ast.InterfaceType:
		return &PlainType{
			Name: name,
			Type: "interface{}",
			Docs: docs,
		}
	default:
		log.Printf("ParseExpr: unhandled type %T for %s\n", expr, names)
	}

	return nil
}

func PlainTypeFromIdent(name, docs string, i *ast.Ident) *PlainType {
	return &PlainType{Docs: docs, Name: name, Type: i.String()}
}

func PlainTypeFromSelectorExpr(name, docs string, s *ast.SelectorExpr) *PlainType {
	return &PlainType{Docs: docs, Name: name, Type: fmt.Sprintf("%s.%s", s.X, s.Sel)}
}

func PlainTypeFromStarExpr(name, docs string, star *ast.StarExpr) *PlainType {
	return &PlainType{Docs: docs, Name: name, Type: "*" + stringFromExpr(star.X)}
}

func ArrayTypeFromSpec(name, docs string, a *ast.ArrayType) *ArrayType {
	return &ArrayType{Docs: docs, Name: name, Type: stringFromExpr(a.Elt)}
}

func MapTypeFromSpec(name, docs string, m *ast.MapType) *MapType {
	return &MapType{
		Docs:      docs,
		Name:      name,
		KeyType:   stringFromExpr(m.Key),
		ValueType: stringFromExpr(m.Value),
	}
}

func StructTypeFromSpec(name, docs string, s *ast.StructType) *StructType {
	st := &StructType{
		Docs: docs,
		Name: name,
	}

FIELD_LOOP:
	for _, f := range s.Fields.List {
		field := FieldFromSpec(f)
		if field.Type == nil {
			continue FIELD_LOOP
		}
		st.Fields = append(st.Fields, field)
	}

	return st
}

func FieldFromSpec(f *ast.Field) *Field {
	docs := DocsFromCommentGroup(f.Doc)
	field := &Field{
		Type: ParseExpr(f.Names, docs, f.Type),
		Tags: make(map[string][]string),
	}
	if f.Tag != nil {
		for _, tag := range strings.Split(strings.Replace(f.Tag.Value, "`", "", -1), " ") {
			split := strings.Split(tag, ":")
			split[1] = strings.Trim(split[1], "\"")
			field.Tags[split[0]] = strings.Split(split[1], ",")
		}
	}
	return field
}

func stringFromExpr(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.String()
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", t.X, t.Sel)
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", stringFromExpr(t.X))
	case *ast.ArrayType:
		return fmt.Sprintf("[]%s", stringFromExpr(t.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", stringFromExpr(t.Key), stringFromExpr(t.Value))
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.FuncType:
		return "func()"
	default:
		log.Printf("stringFromExpr: unhandled type %T for %v\n", t, e)
		return ""
	}
}
