# GRPC TS Web

A protoc plugin for generating browser compatible GRPC clients in Typescript.

Why use this plugin instead of [protoc-gen-grpc-web](https://github.com/grpc/grpc-web)?

This plugin improves the ergonomics of the generated client code in a number of ways:

- emits typescript files only (no .js or .d.ts files).
- unary methods use Promises rather than callbacks to support async/await.
- methods accept and return simple JSON objects instead of proto Message instances.
- uses the same `google-protobuf` runtime library and is compatible with it's well-known types.

## Installation

```bash
npm install --save-dev @freifunkmuc/grpc-ts-web
```

This is a fork of [Place1/protoc-gen-grpc-ts-web](https://github.com/Place1/protoc-gen-grpc-ts-web).
The binary is still called `grpc-ts-web`, so existing `npm run codegen` scripts keep working after switching the dependency.

`google-protobuf` and `grpc-web` are peer dependencies. the generated client
imports them, this package does not. Supported ranges are `google-protobuf@^3.11.4 || ^4`
and `grpc-web@^1.0.7 || ^2`

## Usage

If you've installed the npm package you can run the following command to generate
a Typescript GRPC client from a set of proto files.

```bash
# using the npm command
./node_modules/.bin/grpc-ts-web -o ./out ./path/to/protos/**/*.proto
```

If you'd like to invoke protoc yourself using this plugin then you can use the following
command instead. Binaries ship for `darwin-amd64`, `darwin-arm64`, `linux-amd64`,
`linux-arm64` and `windows-amd64`.

```bash
# using protoc directly
protoc --plugin=protoc-gen-grpc-ts-web=./node_modules/@freifunkmuc/grpc-ts-web/bin/protoc-gen-grpc-ts-web-<platform>-<arch> --grpc-ts-web_out ./sdk
```

## Releasing

Bump the version in `npm/package.json`, then push a matching tag (`v0.3.0` for
version `0.3.0`). The release workflow builds every platform binary and publishes
to npm; it refuses to publish if the tag and `npm/package.json` disagree.
Requires an `NPM_TOKEN` secret with publish rights for the `@freifunkmuc` scope.

## Example Output

You can see an example of the emitted code here in the [example directory](https://github.com/Place1/protoc-gen-grpc-ts-web/tree/master/example).

```typescript
import { Metadata } from 'grpc-web';
import { UserService } from './example_pb';

// hostname of the grpc-web server
const hostname = window.location.origin;

// a callback to add default metadata
// on each grpc request
function metadata(): Metadata {
  return {};
}

async function main() {
  // create an instance of the generated
  // UserService client
  const client = new UserService(hostname, metadata)

  // unary calls return promises so that
  // you can use async/await
  // calls also take simple json objects
  // no more `new Message()` and `message.setField()`
  // everywhere!
  const user = await client.addUser({
    name: 'hello world',
  });

  // response messages are json objects
  console.log(user.id);
  console.log(user.name);
  console.log(user.createDate);
}
```
