# streaming-pipeline

## Repository structure
The repository is a monorepo with a source code of the services located in the `./cmd` directory and compiled binaries located in `./bin`

## How to run all unit tests
`make test-unit`

## File emitter service
Command line service to generate a file randomly in a range from some megabytes to some gigabytes. The format of the record in a file is "123 12345".

### How to build
`make build-fileemitter`

And the binary will be put into the `./bin` directory

### How to run
`./bin/file-emitter`

### Comand line options
| Parsed argments | Resulting value | Defaults         |
| -----------     | -----------     |----------        |
| --min-size=10MB | 10MB            | 2MB              |
| --max-size=3GB  | 3GB             | 10MB             |
| --out=/tmp/files| /tmp/files      | /tmp/fileemitter |