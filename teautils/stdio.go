package teautils

import (
	"fmt"
	"os"
)

func Stderrf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}
