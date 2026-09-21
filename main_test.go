package main

import (
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"

	"github.com/place1/protoc-gen-grpc-ts-web/templates"
)

func TestParseOptions(t *testing.T) {
	tests := []struct {
		parameter string
		want      templates.Format
		wantErr   bool
	}{
		// what earlier versions generated, so existing users see no change
		{parameter: "", want: templates.FormatText},
		{parameter: "format=text", want: templates.FormatText},
		{parameter: "format=binary", want: templates.FormatBinary},
		{parameter: " format=binary ", want: templates.FormatBinary},
		{parameter: "format=binary,", want: templates.FormatBinary},
		{parameter: "format=json", wantErr: true},
		{parameter: "binary", wantErr: true},
		{parameter: "mode=grpcweb", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.parameter, func(t *testing.T) {
			options, err := parseOptions(tt.parameter)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseOptions(%q) accepted it, want an error", tt.parameter)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseOptions(%q) = %v", tt.parameter, err)
			}
			if options.Format != tt.want {
				t.Errorf("format = %q, want %q", options.Format, tt.want)
			}
		})
	}
}

// request describes one file with one service, the way protoc would.
func request(parameter string) *plugin.CodeGeneratorRequest {
	message := func(name string) *descriptor.DescriptorProto {
		return &descriptor.DescriptorProto{Name: proto.String(name)}
	}
	return &plugin.CodeGeneratorRequest{
		Parameter:      proto.String(parameter),
		FileToGenerate: []string{"example.proto"},
		ProtoFile: []*descriptor.FileDescriptorProto{{
			Name:        proto.String("example.proto"),
			Package:     proto.String("example"),
			MessageType: []*descriptor.DescriptorProto{message("Req"), message("Res")},
			Service: []*descriptor.ServiceDescriptorProto{{
				Name: proto.String("Example"),
				Method: []*descriptor.MethodDescriptorProto{{
					Name:       proto.String("Call"),
					InputType:  proto.String(".example.Req"),
					OutputType: proto.String(".example.Res"),
				}},
			}},
		}},
	}
}

func generatedCode(t *testing.T, parameter string) string {
	t.Helper()
	res := generate(request(parameter))
	if res.Error != nil {
		t.Fatalf("generation failed: %s", res.GetError())
	}
	if len(res.File) != 1 {
		t.Fatalf("got %d files, want 1", len(res.File))
	}
	return res.File[0].GetContent()
}

func TestGeneratedClientUsesTheFormat(t *testing.T) {
	for parameter, want := range map[string]string{
		"":              "format: 'text',",
		"format=text":   "format: 'text',",
		"format=binary": "format: 'binary',",
	} {
		t.Run(parameter, func(t *testing.T) {
			if code := generatedCode(t, parameter); !strings.Contains(code, want) {
				t.Errorf("generated client does not contain %q", want)
			}
		})
	}
}

// protoc prints the error and fails; the plugin must not crash or emit files.
func TestBadOptionIsReportedToProtoc(t *testing.T) {
	res := generate(request("format=json"))
	if res.Error == nil {
		t.Fatal("no error for an unknown format")
	}
	if !strings.Contains(res.GetError(), "json") {
		t.Errorf("error %q does not name the bad value", res.GetError())
	}
	if len(res.File) != 0 {
		t.Errorf("emitted %d files alongside the error", len(res.File))
	}
}
