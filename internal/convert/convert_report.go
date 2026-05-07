package convert

import (
	"fmt"
	"io"
)

// WriteSummary writes a brief summary of the conversion to w.
func WriteSummary(w io.Writer, source string, format Format, keyCount int) {
	fmt.Fprintf(w, "Source:  %s\n", source)
	fmt.Fprintf(w, "Format:  %s\n", format)
	fmt.Fprintf(w, "Keys:    %d\n", keyCount)
}

// WriteFormats lists all supported output formats to w.
func WriteFormats(w io.Writer) {
	fmt.Fprintln(w, "Supported formats:")
	for _, f := range []Format{FormatDotenv, FormatExport, FormatJSON, FormatYAML} {
		fmt.Fprintf(w, "  - %s\n", f)
	}
}
