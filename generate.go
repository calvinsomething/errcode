package main

import (
	_ "embed"
	"errors"
	"fmt"
	"iter"
	"log"
	"maps"
	"os"
	"slices"
)

var generatedFile *os.File

// initMissingPackage loads the generated package, creating the directory and/or errcode_gen.go if they do not exist
func initMissingPackage() bool {
	if _, err := os.Stat("go.mod"); errors.Is(err, os.ErrNotExist) {
		log.Println("Could not find go.mod file. Working directory must be module root.")
	} else if err != nil {
		log.Fatal(err)
	}

	if f, err := os.OpenFile("errcode/errcode_gen.go", os.O_RDWR, 0); err == nil {
		generatedFile = f
		return false
	} else if !errors.Is(err, os.ErrNotExist) {
		log.Fatal(err)
	}

	err := os.Mkdir("errcode", 0755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	generatedFile, err = os.Create("errcode/errcode_gen.go")
	if err != nil {
		log.Fatal(err)
	}

	doInitialFileWrite()

	if err = generatedFile.Sync(); err != nil {
		log.Fatal(err)
	}

	log.Println("errcode package added to project")

	return true
}

//go:embed errcode.go
var initialFileContents []byte

func doInitialFileWrite() {
	input := initialFileContents

	// skip line with build directive + following blank
	for i, b := range input {
		if b == '\n' {
			input = input[i+2:]
			break
		}
	}

	if n, err := generatedFile.Write(input); err != nil {
		log.Fatal(err)
	} else if n != len(input) {
		log.Fatalf("wrote %d bytes to generated file, attempted to write %d", n, len(input))
	}
}

const codesBeginComment = `
// Error Codes
`

const declFmt = `
// %s returns an Error constructed by NewErrorFunc.
func %s(m string, fa ...interface{}) Error {
	return NewErrorFunc("%s", m, fa...)
}
`

func generateDecls() {
	if err := generatedFile.Truncate(0); err != nil {
		log.Fatal(err)
	}

	doInitialFileWrite()

	shouldAddCodesStartComment := true

	for _, code := range slices.Sorted(maps.Keys(decls)) {
		decl := decls[code]

		if len(decl.uses) > 1 {
			log.Printf("multiple calls found for code '%s'", code)
			for _, use := range decl.uses {
				log.Println(use.fset.Position(use.pos).String())
			}
		} else if len(decl.uses) == 0 {
			log.Printf("removing unused func '%s'", code)
			continue
		}

		if shouldAddCodesStartComment {
			generatedFile.WriteString(codesBeginComment)
			shouldAddCodesStartComment = false
		}

		fmt.Fprintf(generatedFile, declFmt, code, code, code)
	}

	if err := generatedFile.Sync(); err != nil {
		log.Fatal(err)
	}
}
