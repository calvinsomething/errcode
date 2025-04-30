package main

import (
	_ "embed"
	"go/token"
	"log"
	"os"
)

//go:embed help.txt
var help []byte

func printHelp(_ string) {
	_, err := os.Stdout.Write(help)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}

func parseOption(a string) string {
	if len(a) <= 2 || a[:2] != "--" {
		log.Fatal("invalid command argument: " + a)
	}

	return a[2:]
}

func setPrefix(a string) {
	if !token.IsIdentifier(a) || !token.IsExported(a) && len(a) > 10 {
		log.Fatal("prefix must be a valid exportable Go identifier and at most 10 characters long")
	}

	prefix = a
}

func getNextHandler(a string) func(string) {
	option := parseOption(a)

	switch option {
	case "help":
		return printHelp
	case "prefix":
		return setPrefix
	default:
		log.Fatal("invalid option: --" + option)
		return nil
	}
}

func handleOptions() {
	var nextHandler func(string)

	for _, a := range os.Args[1:] {
		if nextHandler == nil {
			nextHandler = getNextHandler(a)
		} else {
			nextHandler(a)
			nextHandler = nil
		}
	}

	if nextHandler != nil {
		nextHandler("")
	}
}
