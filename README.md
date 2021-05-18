# oslo

CLI tool for the OpenSLO spec

## Prerequisites

- [Go](https://golang.org/)

## Installation

1. Checkout this repository
1. Install oslo with `go get github.com/OpenSLO/oslo`

## Usage

Right now, the only function is `validate`, which you can call with `oslo validate`

## Testing

To test out the features of oslo, from the root of the project run
`oslo validate test/valid-service.yaml`
That will validate against a valid yaml file.  There are other files in that
directory to test out the functionality of `oslo`
