// Package shellquote provides functions for quoting strings so they can be
// passed to a shell.
//
// See "Shell Command Language",
// http://pubs.opengroup.org/onlinepubs/009695399/utilities/xcu_chap02.html
package shellquote

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	maskPlain = 1 << iota
	maskSingle
	maskDouble
	maskEscape
)

var charMask [128]uint8 = func() (r [128]uint8) {
	const maskAll = maskPlain | maskSingle | maskDouble
	for c := 32; c <= 126; c++ {
		r[c] = maskSingle | maskDouble
	}
	for c := 'A'; c <= 'Z'; c++ {
		r[c] = maskAll
	}
	for c := 'a'; c <= 'z'; c++ {
		r[c] = maskAll
	}
	for c := '0'; c <= '9'; c++ {
		r[c] = maskAll
	}
	for _, c := range "%+,-./:=@_" {
		r[c] = maskAll
	}
	r['\''] = maskDouble
	for _, c := range "$\"\\`" {
		r[c] = maskSingle | maskDouble | maskEscape
	}
	return
}()

func stringMask(s string) (int, error) {
	m := maskPlain | maskSingle | maskDouble
	for _, c := range s {
		if c >= 128 {
			return 0, fmt.Errorf(
				"string contains non-ASCII character U+%04X", c)
		}
		cm := charMask[c]
		m &= int(cm)
		if cm == 0 {
			return 0, fmt.Errorf(
				"string contains control character U+%04X", c)
		}
	}
	return m, nil
}

func writeDoubleQuoted(w *bytes.Buffer, s string) {
	w.WriteByte('"')
	for _, c := range s {
		if charMask[c]&maskEscape != 0 {
			w.WriteByte('\\')
		}
		w.WriteByte(byte(c))
	}
	w.WriteByte('"')
}

// String quotes a string so it can be passed to the shell.
func String(s string) (string, error) {
	if s == "" {
		return "''", nil
	}
	switch m, err := stringMask(s); {
	case err != nil:
		return "", err
	case m&maskPlain != 0:
		return s, nil
	case m&maskSingle != 0:
		return "'" + s + "'", nil
	default:
		var buf bytes.Buffer
		writeDoubleQuoted(&buf, s)
		return buf.String(), nil
	}
}

// BareQuotedString quotes a string so that it can be passed to the shell
// between double quotes. The quotes are not included.
func BareQuotedString(s string) (string, error) {
	switch m, err := stringMask(s); {
	case err != nil:
		return "", err
	case m&maskPlain != 0:
		return s, nil
	default:
		var buf bytes.Buffer
		for _, c := range s {
			if charMask[c]&maskEscape != 0 {
				buf.WriteByte('\\')
			}
			buf.WriteByte(byte(c))
		}
		return buf.String(), nil
	}
}

// Command quotes a command so it can be passed to the shell.
func Command(cmd []string) (string, error) {
	var buf bytes.Buffer
	if err := WriteCommand(&buf, cmd); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WriteString writes a string to a buffer so it can be passed to the shell.
func WriteString(w *bytes.Buffer, s string) error {
	if s == "" {
		w.WriteString("''")
		return nil
	}
	switch m, err := stringMask(s); {
	case err != nil:
		return err
	case m&maskPlain != 0:
		w.WriteString(s)
	case m&maskSingle != 0:
		w.WriteByte('\'')
		w.WriteString(s)
		w.WriteByte('\'')
	default:
		writeDoubleQuoted(w, s)
	}
	return nil
}

// WriteCommand writes a shell command to a buffer.
func WriteCommand(w *bytes.Buffer, cmd []string) error {
	for i, arg := range cmd {
		if i > 0 {
			w.WriteByte(' ')
		}
		if err := WriteString(w, arg); err != nil {
			return err
		}
	}
	return nil
}

// LocalPath transforms a string referring to a local file so it unambiguously
// refers to a local file, and won't be misinterpreted as a command-line flag or
// remote file when passed to a program as an argument. Strings beginning with
// '-' can be misinterpreted as flags, and strings containing ':' before the
// first '/' can be misinterpreted as remote files. This function prepends './'
// to the string in both of these cases. If the string is the empty string, '.'
// is returned.
func LocalPath(s string) string {
	if s == "" {
		return "."
	}
	if s[0] == '-' || s[0] == '~' {
		return "./" + s
	}
	i := strings.IndexByte(s, '/')
	leading := s
	if i != -1 {
		leading = leading[:i]
	}
	j := strings.IndexByte(leading, ':')
	if j != -1 {
		return "./" + s
	}
	return s
}
