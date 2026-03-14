# errcode
A simple codegen module for generating unique, stable error codes as function identifiers, making errors easily traceable.

### download the package
```
go get github.com/calvinsomething/errcode
```

### run errcode
Add `//go:generate go run github.com/calvinsomething/errcode` to any Go file in your project, then run `go generate` from the project root. The first time the program runs, it will add an errcode package to your project.

### use error codes
Use `errcode.New` in your code to construct an `errcode.Error`. Running `go generate` now will replace any instance of `New` with an equivalent function with an error code as an identifier, and an `Error.Code` value equal to the same error code used as the function name.

```
	return errcode.New("Invalid image format.")
```

becomes:

```
	return errcode.K9S6O("Invalid image format.")
```

The functions will be declared in the generated errcode package.

```
// K9S6O returns an Error constructed by NewErrorFunc.
func K9S6O(m string, fa ...interface{}) Error {
	return NewErrorFunc("K9S6O", m, fa...)
}

```

### add a prefix
Append `--prefix <prefix>` to the run command to add a prefix to any new error code generated; helpful for differentiating between error code sources.

### reassignable implementations
`NewErrorFunc`, `ErrorStringFunc`, and `WrapFunc` can be assigned for customizing the behavior of `New`, `Error.Error` and `Error.Wrap`, respectively.

You can reassign these variables from outside of the generated `errcode` package, or you can add a file within the package for access to private members.

Example:
```
package errcode

import (
	"fmt"
	"runtime"
)

func myWrapFunc(e Error, w error) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "'unknown file'"
		line = -1
	}

	e.p.wrapped = fmt.Errorf("%s:%d: %w", file, line, w)
}

func init() {
	WrapFunc = myWrapFunc
}
```
