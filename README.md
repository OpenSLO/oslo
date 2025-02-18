#

<!-- markdownlint-disable MD033-->
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="images/openslo_light.png">
  <img alt="OpenSLO light theme" src="images/openslo.png">
</picture>
<!-- markdownlint-enable MD033-->

---

CLI tool for interacting with the [OpenSLO specification](https://github.com/OpenSLO/OpenSLO)!

## Installation

### Prebuilt binaries

Download prebuilt binaries from the
[published release assets](https://github.com/OpenSLO/oslo/releases/latest).

### Go install

```sh
go install github.com/OpenSLO/oslo/cmd/oslo@latest
```

### Homebrew

```sh
brew install openslo/openslo/oslo
```

### From Docker

```sh
docker run -v "$(pwd):/manifests" ghcr.io/openslo/oslo:latest <command> /manifests/<file>.yaml
```

For example, if you have an OpenSLO spec file in the current directory called `my-service.yaml`,
and you wanted to validate it, the full command would be:

```sh
docker run -v "$(pwd):/manifests" ghcr.io/openslo/oslo:latest validate -f /manifests/my-service.yaml
# Valid!
```

### From source

1. Clone this repository.
2. From the root of the project, run `make install`.
   This will build and install the binary into your `GOPATH`.

## Usage

### Validate

`oslo validate` will validate the provided OpenSLO YAML document.

Example:

```sh
oslo validate -f file1.yaml -f file2.yaml
```

### Format

`oslo fmt` will validate the provided OpenSLO YAML document.

Example:

```sh
oslo validate -f file1.yaml -f file2.yaml
```
