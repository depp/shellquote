# ShellQuote

ShellQuote is a Go library that quotes strings so they can be used by a POSIX shell.

## Installation

ShellQuote can be installed with `go get`.

```shell
go get github.com/depp/shellquote
```

## Usage

```go
import "github.com/depp/shellquote"
```

To quote a string, call `String()`. This adds quotes only as necessary. For example,

```go
// Prints x without quotes, because no quotes are needed.
fmt.Println(shellquote.String("x"))
// Prints '$x' with quotes.
fmt.Println(shellquote.String("$x"))
```

Another useful function is `LocalPath()`, which modifies a path so it is not misinterpreted as a flag.

```go
// Prints abc, because abc does not look like a flag.
fmt.Println(shellquote.LocalPath("abc"))
// Prints ./-h, because -t looks like a flag.
fmt.Println(shellquote.LocalPath("-h"))
```

## Limitations

Non-ASCII characters are not supported.

## License

This is licensed under the MIT license. See LICENSE.txt.
