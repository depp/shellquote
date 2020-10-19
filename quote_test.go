package shellquote

import (
	"bytes"
	"testing"
)

type testcase struct {
	input  string
	output string
}

func TestQuote(t *testing.T) {
	cases := []testcase{
		{"a", "a"},
		{"azAZ09_-%+,-./:=@_", "azAZ09_-%+,-./:=@_"},
		{"", "''"},
		// Special characters
		{"|", "'|'"},
		{"&", "'&'"},
		{";", "';'"},
		{"<", "'<'"},
		{">", "'>'"},
		{"$", "'$'"},
		{"`", "'`'"},
		{"\\", "'\\'"},
		{"\"", "'\"'"},
		{"'", `"'"`},
		{" ", "' '"},
		{"`'\\$\"", "\"\\`'\\\\\\$\\\"\""},
	}
	for _, c := range cases {
		if s, err := String(c.input); err != nil {
			t.Errorf("String(%q): error: %v", c.input, err)
		} else if s != c.output {
			t.Errorf("String(%q): got %q, expected %q", c.input, s, c.output)
		}
		var buf bytes.Buffer
		if err := WriteString(&buf, c.input); err != nil {
			t.Errorf("WriteString(%q): error: %v", c.input, err)
		} else {
			s := buf.String()
			if s != c.output {
				t.Errorf("WriteString(%q): got %q, expected %q",
					c.input, s, c.output)
			}
		}
	}
}

func TestLocalPath(t *testing.T) {
	cases := []testcase{
		{"", "."},
		{"..", ".."},
		{"-flag", "./-flag"},
		{"~tilde", "./~tilde"},
		{"not-flag", "not-flag"},
		{"remote:file.txt", "./remote:file.txt"},
		{"remote:dir/file.txt", "./remote:dir/file.txt"},
		{"local/file:with-colon", "local/file:with-colon"},
	}
	for _, c := range cases {
		s := LocalPath(c.input)
		if s != c.output {
			t.Errorf("LocalPath(%q): got %q, expected %q", c.input, s, c.output)
		}
	}
}
