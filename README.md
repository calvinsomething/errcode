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

<img src="https://github.com/user-attachments/assets/f9b1e7ba-e720-404a-9c7c-7786dd2a6dc1" alt="Before" height="50">

**_becomes_**

<img src="https://github.com/user-attachments/assets/2595e0df-6d86-4b63-b865-225f29068b3c" alt="After" height="50" style="margin-bottom: 10px">
<br></br>
The functions will be declared in the generated errcode package.

<img src="https://github.com/user-attachments/assets/b93267e3-3a3c-4f6d-9c48-6c6b7688bf8a" alt="Declaration" height="160">

### add a prefix
Append `--prefix <prefix>` to the run command to add a prefix to any new error code generated; helpful for differentiating between error code sources.

### use your own constructor
Set `errcode.NewErrorFunc` to override the default `Error` construction.
