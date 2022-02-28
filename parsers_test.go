package toast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGoFileFromAST(t *testing.T) {
	path := "./internal/mock"
	loadFiles(t, path, "go")
}

func TestCUEFileFromAST(t *testing.T) {
	path := "./internal/mock"
	loadFiles(t, path, "cue")
}

func loadFiles(t *testing.T, dirPath, output string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		filePath := fmt.Sprintf("%s/%s", dirPath, entry.Name())
		load(t, filePath, output)
	}
}

func load(t *testing.T, path, output string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	f := FileFromAST(file, "", &GenEnumTypeTransform{
		Generate: func(docs string, spec *ast.ValueSpec) *PromoteToEnumType {
			etName := strings.Replace(docs, "// Enum value maps for ", "", 1)
			if len(etName) == len(docs) ||
				len(spec.Names) < 1 ||
				!strings.HasSuffix(spec.Names[0].Name, "_name") ||
				len(spec.Values) < 1 {
				return nil
			}

			et := &EnumType{
				Name: etName[:len(etName)-2],
			}
			for _, expr := range spec.Values[0].(*ast.CompositeLit).Elts {
				v := expr.(*ast.KeyValueExpr).Value.(*ast.BasicLit).Value
				et.Values = append(et.Values, strings.Replace(v, "\"", "", -1))
			}

			return &PromoteToEnumType{
				Apply: func(pt *PlainType) *EnumType {
					if pt.Name != et.Name {
						return nil
					}
					et.Docs = pt.Docs
					return et
				},
			}
		},
	})

	// f.debug = true

	outputPath := path[strings.LastIndex(path, "/")+1:]
	var out string
	if output == "go" {
		out = f.Go()
	} else {
		out = f.CUE()
		outputPath = outputPath[:len(outputPath)-3] + ".cue"
	}

	if err := os.MkdirAll("./.testdata/", 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("./.testdata/"+outputPath, []byte(out), 0644); err != nil {
		panic(err)
	}
}
