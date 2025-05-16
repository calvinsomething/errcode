package main

import (
	"go/ast"
	"go/format"
	"log"
	"math/rand"
	"os"

	"golang.org/x/sync/errgroup"
)

const (
	newIdent           = "New"
	exportedErrorIdent = "Error"
	length             = 5
)

var prefix string

// newCode returns a random string of alternating capital letters and numbers
func newCode() string {
	code := make([]rune, length)

	for i := range code {
		if i%2 == 0 {
			code[i] = 'A' + rune(rand.Uint32()%26)
		} else {
			code[i] = '1' + rune(rand.Uint32()%9)
		}
	}

	return prefix + string(code)
}

// updateCalls rewrites any errcode.New calls to errcode.<generated code>
//
// updates AST, but does not write to file
func updateCalls() {
	var didErr bool

	for _, f := range files {
		ast.Inspect(f.ast, func(n ast.Node) bool {
			if x, ok := n.(*ast.CallExpr); ok {
				if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
					obj := f.pkg.TypesInfo.ObjectOf(sel.Sel)
					if obj == nil {
						didErr = true

						p := f.pkg.Fset.Position(sel.Pos())
						log.Printf("%s:%d:%d: undefined identifier '%s'", p.Filename, p.Line, p.Column, sel.Sel.Name)

						return false
					}

					if obj.Pkg().Scope() == errcodeScope {
						if sel.Sel.Name == newIdent {
							code := newCode()
							for decls[code] != nil {
								code = newCode()
							}

							sel.Sel.Name = code

							f.didUpdate = true

							decls[code] = &decl{name: sel.Sel.Name}
						} else if _, ok := originalExportedIdents[sel.Sel.Name]; !ok {
							if decl, ok := decls[sel.Sel.Name]; ok {
								decl.uses = append(decl.uses, funcUse{
									pos:  f.pkg.TypesInfo.ObjectOf(sel.Sel).Pos(),
									fset: f.pkg.Fset,
								})
							} else {
								// should not be possible
								log.Fatalf("unknown errcode identifier '%s'", sel.Sel.Name)
							}
						}
					}

					return false
				}
			}

			return true
		})
	}

	if didErr {
		log.Fatal("errcode exited early")
	}
}

// writeFiles writes the updated ASTs to their respective source files
func writeFiles() {
	var async errgroup.Group

	for _, f := range files {
		if f.didUpdate {
			async.Go(func() error {
				dst, err := os.OpenFile(f.fileName, os.O_WRONLY, 0)
				if err != nil {
					return err
				}

				if err = format.Node(dst, f.pkg.Fset, f.ast); err != nil {
					return err
				}

				if err = dst.Sync(); err != nil {
					return err
				}

				log.Println("updated errcode function calls in " + f.fileName)

				return nil
			})
		}
	}

	if err := async.Wait(); err != nil {
		log.Fatal(err)
	}
}
