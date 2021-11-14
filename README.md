# streaming-pipeline

The actual implementation of a streaming files pipeline explained [here](https://github.com/valerykalashnikov/streaming-pipeline/blob/master/design/Design.md#worker-pool-with-message-queue)

## Repository structure
The repository is a monorepo with a source code of the services located in the `./cmd` directory.It contains of
  - [file emitter service](#file-emitter-service)
  - [publisher service](#publisher-service)
  - [consumer service](#consumer-service)
  
Compiled binaries are located in `./bin`

## How to run all unit tests
`make test-unit`

## How to run an environment
`make env-start`

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

## Publisher service

## Consumer service