#!/usr/bin/env node

// commander >= 5 exports the program instead of being it
const { program } = require('commander');
const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');

function spawnAsync(cmd, args) {
  return new Promise((resolve, reject) => {
    console.info(`executing: ${cmd} ${args.join(' ')}`)
    const p = spawn(cmd, args)
    p.stdout.pipe(process.stdout);
    p.stderr.pipe(process.stderr);
    // protoc is not installed, or not on PATH
    p.on('error', reject);
    // A failed protoc run has to fail this command too: otherwise
    // "npm run codegen" reports success and the old files stay in place.
    p.on('exit', (code, signal) => {
      if (code === 0) {
        resolve();
      } else if (signal) {
        reject(new Error(`${cmd} was killed by ${signal}`));
      } else {
        reject(new Error(`${cmd} exited with code ${code}`));
      }
    });
  });
}

function platform() {
  switch (process.platform) {
    case 'win32':
      return 'windows';
    default:
      return process.platform;
  }
}

function arch() {
  switch (process.arch) {
    case 'x64':
      return 'amd64';
    case 'arm64':
      return 'arm64';
    default:
      return process.arch;
  }
}

const pluginPath = path.join(__dirname, `bin/protoc-gen-grpc-ts-web-${platform()}-${arch()}`);

if (!fs.existsSync(pluginPath)) {
  console.error(
    `grpc-ts-web: no plugin binary for ${process.platform}/${process.arch} ` +
      `(looked for ${pluginPath}).\n` +
      'Build one from source with `make build` in ' +
      'https://github.com/freifunkMUC/protoc-gen-grpc-ts-web and put it there, ' +
      'or open an issue asking for the platform to be added to the release.',
  );
  process.exit(1);
}

program.arguments('<protos...>')
  .requiredOption('-o, --out <directory>', 'a directory to write the generated code to')
  .option(
    '--format <format>',
    'the gRPC-Web wire format of the generated client: "text" (default) or "binary"',
  )
  .action((protos, options) => {
    if (!fs.existsSync(options.out)) {
      fs.mkdirSync(options.out);
    }
    const includes = protos
      .map(p => path.dirname(p))
      .map(p => `--proto_path=${p}`)
      .filter((value, index, self) => self.indexOf(value) === index);
    // left to the plugin to validate, so the CLI and protoc agree on the rules
    const pluginOptions = options.format ? [`--grpc-ts-web_opt=format=${options.format}`] : [];
    return spawnAsync('protoc', [
      '--grpc-ts-web_out',
      options.out,
      ...pluginOptions,
      `--plugin=protoc-gen-grpc-ts-web=${pluginPath}`,
      ...includes,
      ...protos,
    ]).catch((err) => {
      console.error(`grpc-ts-web: ${err.message}`);
      process.exitCode = 1;
    });
  });

program.parse(process.argv);
