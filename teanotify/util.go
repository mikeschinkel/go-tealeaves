package teanotify

import (
	"bytes"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

// getLines splits s on newlines and returns the lines along with the
// width (in terminal cells) of the widest line.
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")
	for _, l := range lines {
		w := ansi.StringWidth(l)
		if widest < w {
			widest = w
		}
	}
	return lines, widest
}

// isANSITerminator reports whether c is a valid ANSI CSI sequence terminator.
// Terminators are in the ranges 0x40–0x5A (uppercase letters and @) and
// 0x61–0x7A (lowercase letters).
func isANSITerminator(c rune) (ok bool) {
	ok = (c >= 0x40 && c <= 0x5a) || (c >= 0x61 && c <= 0x7a)
	return ok
}

// cutLeft cuts printable characters from the left of s, removing cutWidth
// cells worth of visible content. ANSI escape sequences are preserved.
func cutLeft(s string, cutWidth int) (result string) {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)
	for _, c := range s {
		var w int
		if c == ansi.ESC || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if isANSITerminator(c) {
				isAnsi = false
				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
					ab.Reset()
				}
			}
		} else {
			w = runewidth.RuneWidth(c)
		}

		if pos >= cutWidth {
			if b.Len() == 0 {
				if ab.Len() > 0 {
					b.Write(ab.Bytes())
				}
				if pos-cutWidth > 1 {
					b.WriteByte(' ')
					continue
				}
			}
			b.WriteRune(c)
		}
		pos += w
	}
	result = b.String()
	return result
}

// cutRight keeps printable characters from the left of s, up to keepWidth
// cells. ANSI escape sequences are preserved. Complement to cutLeft.
func cutRight(s string, keepWidth int) (result string) {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)

	for _, c := range s {
		var w int
		if c == ansi.ESC || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if isANSITerminator(c) {
				isAnsi = false
				b.Write(ab.Bytes())
				ab.Reset()
			}
			continue
		}

		w = runewidth.RuneWidth(c)
		if pos+w > keepWidth {
			break
		}

		b.WriteRune(c)
		pos += w
	}

	// Reset to avoid color bleed
	if b.Len() > 0 && !bytes.HasSuffix(b.Bytes(), []byte("[0m")) {
		b.WriteByte(ansi.ESC)
		b.WriteString("[0m")
	}

	result = b.String()
	return result
}

// hangingWrap wraps msg with a prefix and hanging indent for continuation
// lines within textWidth terminal cells.
func hangingWrap(prefix, msg string, textWidth int) (result string) {
	prefix = prefix + " "
	indentW := ansi.StringWidth(prefix)
	avail := textWidth - indentW
	if avail < 1 {
		result = prefix + msg
		goto end
	}

	{
		wrapped := ansi.Wordwrap(msg, avail, " ")
		indent := strings.Repeat(" ", indentW)
		lines := strings.Split(wrapped, "\n")
		for i := 1; i < len(lines); i++ {
			lines[i] = indent + lines[i]
		}
		result = prefix + strings.Join(lines, "\n")
	}

end:
	return result
}
