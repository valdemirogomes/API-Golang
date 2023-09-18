# go_dev_base_3_v2

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

This is a basic Go application created by Fury to be used as a starting point for your project.

## First steps

### Go Runtime Version

Specify the Go runtime version tag you desire in your `Dockerfile`. If in doubt, it's completely safe to always use the
latest one given the [Go 1 compatibility guarantees](https://golang.org/doc/go1compat).

```docker
FROM hub.furycloud.io/mercadolibre/go:1.17-mini
```

> You can find all available image tags for your Dockerfile
> [here](https://github.com/mercadolibre/fury_go-mini#supported-tags).

### Dependency Management

Run `go mod tidy` to ensure that the `go.mod` file matches the source code in the module.

For more information refer to the
[`go-mini` docs](https://github.com/mercadolibre/fury_go-mini#dependency-management-support).

## Questions

* [Fury Issue Tracker](https://github.com/mercadolibre/fury/issues)