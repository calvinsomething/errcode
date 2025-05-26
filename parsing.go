package main

import (
	"go/ast"
	"go/parser"
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

// initialExportedIdents is used to compare against error code constructor declarations
// as well as potential undefined function calls
var initialExportedIdents = map[string]struct{}{}

func loadInitialExportedIdents() {
	fs := token.NewFileSet()

	f, err := parser.ParseFile(fs, "", initialFileContents, 0)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		switch d := n.(type) {
		case *ast.FuncDecl:
			if ast.IsExported(d.Name.Name) {
				initialExportedIdents[d.Name.Name] = struct{}{}
			}
		case *ast.GenDecl:
			for _, s := range d.Specs {
				switch s := s.(type) {
				case *ast.ValueSpec:
					for _, i := range s.Names {
						if ast.IsExported(i.Name) {
							initialExportedIdents[i.Name] = struct{}{}
						}
					}
				case *ast.TypeSpec:
					if ast.IsExported(s.Name.Name) {
						initialExportedIdents[s.Name.Name] = struct{}{}
					}
				}
			}
		}

		return true
	})
}

// loadExistingCodes parses an existing errcode_gen.go file and maps existing
// function declaration identifiers to their call positions
func loadExistingCodes(p *packages.Package, s *ast.File) {
	ast.Inspect(s, func(n ast.Node) bool {
		switch d := n.(type) {
		case *ast.FuncDecl:
			if _, ok := initialExportedIdents[d.Name.Name]; ok {
				if errcodeScope == nil {
					errcodeScope = p.TypesInfo.ObjectOf(d.Name).Pkg().Scope()
				}
			} else if ast.IsExported(d.Name.Name) {
				// add func name to decls if result is a single Error struct
				if d.Type != nil && d.Type.Results != nil {
					if len(d.Type.Results.List) == 1 && types.ExprString(d.Type.Results.List[0].Type) == exportedErrorIdent {
						decls[d.Name.Name] = &decl{name: d.Name.Name, decl: d}
					}
				}
			}
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
