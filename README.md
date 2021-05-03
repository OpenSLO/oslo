# oslo

CLI tool for the OpenSLO spec

## Prerequisites

- [Go](https://golang.org/)

## Installation

1. Checkout this repository
1. Install olso with `go install oslo`

## Usage

Right now, the only function is `validate`, which you can call with `oslo validate`

## Testing

To test out the features of oslo, from the root of the project run
`oslo validate test/valid.yaml`
That will validate against a valid yaml file.  There are two other files in that
directory, one for invalid data, and one for missing data.
