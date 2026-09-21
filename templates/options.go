package templates

import "fmt"

// Format is the wire format the generated client uses to talk gRPC-Web.
type Format string

const (
	// FormatText base64-encodes the messages (application/grpc-web-text).
	// It is what older browsers needed for server streaming over XHR, and
	// what this plugin has always generated.
	FormatText Format = "text"
	// FormatBinary sends the messages as they are
	// (application/grpc-web+proto). It is smaller and cheaper, and the only
	// format some servers accept - connect-go among them.
	FormatBinary Format = "binary"
)

// ParseFormat turns an option value into a Format.
func ParseFormat(value string) (Format, error) {
	switch Format(value) {
	case FormatText, FormatBinary:
		return Format(value), nil
	}
	return "", fmt.Errorf("unknown format %q (known: %s, %s)", value, FormatText, FormatBinary)
}

// Options control the generated code.
type Options struct {
	Format Format
}

// DefaultOptions keeps the output of earlier versions: text format.
func DefaultOptions() Options {
	return Options{Format: FormatText}
}
