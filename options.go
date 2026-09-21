package main

import (
	"fmt"
	"strings"

	"github.com/place1/protoc-gen-grpc-ts-web/templates"
)

// parseOptions reads the parameter protoc hands the plugin: a comma-separated
// list of key=value pairs, as passed with --grpc-ts-web_opt or in front of the
// output directory (--grpc-ts-web_out=format=binary:./out).
func parseOptions(parameter string) (templates.Options, error) {
	options := templates.DefaultOptions()

	for _, pair := range strings.Split(parameter, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		key, value, found := strings.Cut(pair, "=")
		if !found {
			return options, fmt.Errorf("option %q is not of the form key=value", pair)
		}

		switch key {
		case "format":
			format, err := templates.ParseFormat(value)
			if err != nil {
				return options, err
			}
			options.Format = format
		default:
			return options, fmt.Errorf("unknown option %q (known: format)", key)
		}
	}

	return options, nil
}
