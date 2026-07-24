package types

import (
	"fmt"
)

func formatError(msg string, serviceTypeName string, cause error, s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprint(s, msg)
			if serviceTypeName != "" {
				_, _ = fmt.Fprintf(s, ": %s", serviceTypeName)
			}
			if cause != nil {
				_, _ = fmt.Fprintf(s, "\nCause: %+v", cause)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprint(s, msg)
		if serviceTypeName != "" {
			_, _ = fmt.Fprintf(s, ": %s", serviceTypeName)
		}
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", msg)
	}
}
