# oslo

CLI tool for the [OpenSLO spec](https://github.com/OpenSLO/OpenSLO). For more
information also check the website: [openslo.com](https://openslo.com/).

## Prerequisites

- [Go](https://golang.org/)

## Installation

### From source

1. Checkout this repository
1. Install oslo with `go get github.com/OpenSLO/oslo`

### Homebrew

1. `brew tap openslo/openslo`
1. `brew install oslo`

### From Docker

1. `docker run -v "$(pwd):/manifests" ghcr.io/openslo/oslo:latest <command> /manifests/<file>.yaml`

For example, if you had an OpenSLO spec file in the current directory called `myservice.yaml`,
and you wanted to validate it, the full command would be:

```bash
# docker run -v "$(pwd):/manifests" ghcr.io/openslo/oslo:latest validate /manifests/myservice.yaml
Valid!
```

## Usage

Right now, the only function is `validate`, which you can call with `oslo validate`

## Testing

To test out the features of oslo, from the root of the project run
`oslo validate test/valid-service.yaml`
That will validate against a valid yaml file.  There are other files in that
directory to test out the functionality of `oslo`
