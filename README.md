# Extlib Terraform Provider

This terraform provider is an extension library (hence the name "extlib") that provides some additional functionality to the
built in functions in OpenTofu and Terraform.

At the moment this consists of functions. But we may add some basic resources or data sources in the future.

## Requirements

- [OpenTofu](https://opentofu.org/) >= 1.7  or [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.25

## Building the Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the Provider

Add the provider with:

```hcl
terraform {
    required_providers {
        extlib = {
            source = "opentofu.org/tmccombs/extlib"
        }
    }
}
```

Then you can use the following functions:

- `provider::extlib::boolnum()`
- `provider::extlib::boolset()`

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
