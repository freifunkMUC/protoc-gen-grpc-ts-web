OUT=bin
BINARY=protoc-gen-grpc-ts-web
# platform/arch pairs shipped inside the npm package.
# Keep in sync with platform()/arch() in npm/command.js.
TARGETS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

build:
	echo "running build"
	go build -o $(OUT)/$(BINARY) .

release: clean
	@echo "running release"
	@mkdir -p $(OUT)
	@for target in $(TARGETS); do \
		goos=$${target%/*}; goarch=$${target#*/}; \
		echo "  building $$goos/$$goarch"; \
		CGO_ENABLED=0 GOOS=$$goos GOARCH=$$goarch \
			go build -o $(OUT)/$(BINARY)-$$goos-$$goarch . || exit 1; \
	done
	@rm -rf npm/$(OUT)
	@cp -r $(OUT)/ npm/$(OUT)/

test: build
	echo "running test"
	protoc --plugin=$(OUT)/protoc-gen-grpc-ts-web --grpc-ts-web_out=./ ./example/example.proto

test-js:
	rm -r out || true
	mkdir -p out
	protoc -I ./example --js_out=out --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:out ./example/example.proto

clean:
	echo "running clean"
	rm -r $(OUT) || true
	rm -rf npm/$(OUT) || true

publish: release
	echo "running publish"
	cd npm && npm publish

dump:
	# install with the below go get command
	# go get -u sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-dump
	protoc --dump_out=out=code-generator-request.bin:./ ./example/example.proto
