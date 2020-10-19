# ShellQuote

ShellQuote is a Go library that quotes strings so they can be used by a Posix shell.

For example, `shellquote.String("$x")` returns `'$x'`, surrounded by single quotes, so the `$` is not interpreted by the shell. If you pass a simple string like `shellquote.String("x")`, the result is `x`.

Another useful function is `shellquote.LocalPath`, which modifies a path so it is not misinterpreted as a flag. For example, `shellquote.LocalPath("-h")` returns `"./-h"`.

## Limitations

- Non-ASCII not supported

## License

This is licensed under the MIT license. See LICENSE.txt.
