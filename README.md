# errcode
A simple codegen module for adding random error codes as function identifiers into a Go project.

### Download the package.
```
go get github.com/calvinsomething/errcode
```

### Run errcode.
Add `//go:generate go run github.com/calvinsomething/errcode` to any Go file in your project, then run `go generate` from the project root. The first time the program runs, it will add an `errcode` package.

### Use error codes.
Use `errcode.New` in your code to construct an errcode.Error. Running `go generate` now will replace any instance of `New` with an equivalent function with an error code as an identifier, and an `Error.Code` value equal to the same error code used as the function name.

### Help
Use `go run github.com/calvinsomething/errcode --help` to print the command line options.
