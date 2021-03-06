omg-transform
=============

[![wercker status](https://app.wercker.com/status/81ecedb0a515cab33d8e0554dab73271/s/master "wercker status")](https://app.wercker.com/project/byKey/81ecedb0a515cab33d8e0554dab73271)

An [`enaml`](https://github.com/enaml-ops/enaml) based tool
that allows you to perform transformations on bosh manifests.

## Development Setup

Make sure you have a working Go toolchain and have installed [glide](http://glide.sh/).

```sh
go get -d -u github.com/enaml-ops/omg-transform
cd $GOPATH/src/enaml-ops/omg-transform
glide install
go build
go test `glide nv`
```

## Usage

By default, `omg-transform` attempts to read a manifest from standard in
and writes the transformed manifest to standard out.

For example, to run a transformation on a manifest produced by
[`omg-cli`](https://github.com/enaml-ops/omg-cli):

```sh
omg-cli deploy-product --print-manifest cloudfoundry-plugin-linux | omg-transform <TRANSFORM> [flags...]
```

## Transformations

 - `change-network`: change an instance group's network
 - `clone`: clone an instance group
 - `change-az`: change an instance group's AZs
 - `add-vm-extension`: add a vm extension to an existing instance group
 - `add-tags`: add key-value pairs for VM tagging

## Adding a new transformation

Implementing a transformation is straightforward.

 1. Create a new *.go file for your transformation.
 2. Add a type that implements `Transformation`
 3. Implement a function that builds your new type from a set of arguments.
    This function should have the signature `func(args []string) (Transformation, error)`.

    The `args` slice provided to your function will be everything after the transformation
    on the command line: `omg-transform <transformation> **[args]**`.

    Go's built-in `flag.FlagSet` is a great way to parse these arguments.
    Take a look at the [clone transformation](clone_instance_group.go)
    for an example.
 4. Register your builder function in package main's `init()` function.
    The name that you provide to `RegisterTransformationBuilder()` is the
    command that users will use to invoke your transformation.

 Don't forget to add tests!
