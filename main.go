package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/place1/protoc-gen-grpc-ts-web/templates"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	plugin "google.golang.org/protobuf/types/pluginpb"
)

var input = flag.String("code-generator-request", "", "A path to a protobuf encoded CodeGeneratorRequest")

func main() {
	flag.Parse()

	reader := os.Stdin
	if *input != "" {
		f, err := os.OpenFile(*input, os.O_RDONLY, 0)
		if err != nil {
			log.Fatal(errors.Wrap(err, "failed to read code-generator-request file"))
		}
		defer f.Close()
		reader = f
	}

	data, err := ioutil.ReadAll(reader)
	// data, err := ioutil.ReadFile("example-stdin.bin")
	if err != nil {
		log.Fatal(errors.Wrap(err, "reading input"))
	}

	req := &plugin.CodeGeneratorRequest{}
	if err := proto.Unmarshal(data, req); err != nil {
		log.Fatal(errors.Wrap(err, "bad codegen request"))
	}

	res := generate(req)

	out, err := proto.Marshal(res)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to marshal codegen response"))
	}
	os.Stdout.Write(out)
}

// generate builds the response for one protoc run. A bad option is reported
// through the response's error field, which is how protoc expects a plugin to
// fail: it prints the message and exits non-zero instead of choking on output
// it cannot parse.
func generate(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	options, err := parseOptions(req.GetParameter())
	if err != nil {
		return &plugin.CodeGeneratorResponse{Error: proto.String(err.Error())}
	}

	depLookupTable := templates.NewDependencyLookupTable(req)

	res := &plugin.CodeGeneratorResponse{}
	for _, f := range req.ProtoFile {
		res.File = append(res.File, templates.NewFile(f, depLookupTable, options)...)
	}
	return res
}
