# errcode
A simple codegen module for adding random error codes as function identifiers into a Go project.

### download the package
```
go get github.com/calvinsomething/errcode
```

### run errcode
Add `//go:generate go run github.com/calvinsomething/errcode` to any Go file in your project, then run `go generate` from the project root. The first time the program runs, it will add an errcode package to your project.

### use error codes
Use `errcode.New` in your code to construct an `errcode.Error`. Running `go generate` now will replace any instance of `New` with an equivalent function with an error code as an identifier, and an `Error.Code` value equal to the same error code used as the function name.

```
  return errcode.New(err, "Invalid image format.")
```

becomes:

```
  return errcode.K9S6O(err, "Invalid image format.")
```

The functions will be declared in the generated errcode package.

```
// K9S6O returns an Error constructed by NewErrorFunc.
func K9S6O(e error, m ...string) Error {
	return NewErrorFunc(e, "K9S6O", m...)
}

```

### add a prefix
Append `--prefix <prefix>` to the run command to add a prefix to any new error code generated; helpful for differentiating between error code sources.

### use your own constructor
Set `errcode.NewErrorFunc` to override the default `Error` construction.
