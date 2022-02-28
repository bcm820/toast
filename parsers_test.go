package toast

import (
	"fmt"
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

	f := FileFromAST(file, WithCUEPackageName("mockcue"))
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
