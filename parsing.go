package main

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/packages"
)

type (
	fileSyntax struct {
		fileName  string
		pkg       *packages.Package
		ast       *ast.File
		didUpdate bool
	}

	funcUse struct {
		pos  token.Pos
		fset *token.FileSet
	}

	decl struct {
		name string
		uses []funcUse
		decl *ast.FuncDecl
	}
)

var (
	errcodePath  string
	errcodeScope *types.Scope

	files []*fileSyntax
	decls = map[string]*decl{}
)

// loadExistingCodes parses an existing errcode_gen.go file and maps existing
// function declaration identifiers to their call positions
func loadExistingCodes(p *packages.Package, s *ast.File) {
	ast.Inspect(s, func(n ast.Node) bool {
		switch d := n.(type) {
		case *ast.FuncDecl:
			if d.Name.Name == newIdent {
				errcodeScope = p.TypesInfo.ObjectOf(d.Name).Pkg().Scope()
			} else if ast.IsExported(d.Name.Name) {
				// add func name to decls if result is a single Error struct
				if d.Type != nil && d.Type.Results != nil {
					if len(d.Type.Results.List) == 1 && types.ExprString(d.Type.Results.List[0].Type) == exportedErrorIdent {
						decls[d.Name.Name] = &decl{name: d.Name.Name, decl: d}
					}
				}
			}

			return false
		}

		return true
	})
}

// storeFileThatImportsErrcode checks imports of file, and saves AST data if errcodePath is found
func storeFileThatImportsErrcode(p *packages.Package, f *ast.File, fileName string) {
	for _, impt := range f.Imports {
		// remove quotes from import and compare it to packagePath
		if len(impt.Path.Value) == (len(errcodePath)+2) && impt.Path.Value[1:len(errcodePath)+1] == errcodePath {
			files = append(files, &fileSyntax{
				fileName: fileName,
				pkg:      p,
				ast:      f,
			})
			return
		}
	}
}

// loadTargetModule parses all packages in the module, storing any AST data for files that import errcode
func loadTargetModule() {
	cfg := packages.Config{Mode: packages.NeedModule | packages.NeedName | packages.NeedCompiledGoFiles |
		packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo}

	pkgs, err := packages.Load(&cfg, "./...")
	if err != nil {
		log.Fatal(err)
	} else if len(pkgs) == 0 {
		log.Fatal("failed to load module")
	}

	for _, p := range pkgs {
		if p.Module != nil && p.Module.Main {
			errcodePath = p.Module.Path + "/errcode"
			break
		}
	}

	if errcodePath == "" {
		log.Fatal("could not get module path")
	}

	for _, p := range pkgs {
		for i, s := range p.Syntax {
			fname := p.CompiledGoFiles[i]
			if p.PkgPath == errcodePath {
				if len(fname) > 7 && fname[len(fname)-7:] == "_gen.go" {
					loadExistingCodes(p, s)
				}
			} else {
				storeFileThatImportsErrcode(p, s, fname)
			}
		}
	}
}
