# oslo

CLI tool for the [OpenSLO spec](https://github.com/OpenSLO/OpenSLO). For more
information also check the website: [openslo.com](https://openslo.com/).

## Prerequisites

- [Go](https://golang.org/)

## Installation

### Go install

```sh
go install github.com/OpenSLO/oslo/cmd/oslo@latest
```

### From source

1. Checkout this repository
1. From the root of the project, run `make install`.  This will build and install
the binary into your `GOPATH`

### Homebrew

```sh
brew install openslo/openslo/oslo
```

### From Docker

1. `docker run -v "$(pwd):/manifests" ghcr.io/openslo/oslo:latest <command> /manifests/<file>.yaml`

For example, if you had an OpenSLO spec file in the current directory called `myservice.yaml`,
and you wanted to validate it, the full command would be:

```bash
# docker run -v "$(pwd):/manifests" ghcr.io/openslo/oslo:latest validate /manifests/myservice.yaml
Valid!
```

## Usage

### Validate

`oslo validate` will validate the provided OpenSLO YAML document

### Convert

`oslo convert` will convert the given OpenSLO YAML document to the provided
format.

example:

```bash
oslo convert -f file1.yaml -f file2.yaml -o nobl9
```

That will take the provided yaml files, convert them to Nobl9 formatted config
format, and output to stdout.

*NOTE:* Currently only Nobl9 is supported for output. Additionally, deeply nested
metric sources are not supported. For metric sources that might have a deeply
nested structure, we support a flattened structure, e.g.

```yaml
metricSource:
  type: Instana
  spec:
    infrastructure.query: "myQuery"
    infrastructure.metricRetrievalMethod: "myMetricRetrievalMethod"
```

## Testing

To test out the features of oslo, from the root of the project run
`oslo validate test/valid-service.yaml`
That will validate against a valid yaml file.  There are other files in that
directory to test out the functionality of `oslo`
