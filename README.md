omg-transform
=============

An [`enaml`](https://github.com/enaml-ops/enaml) based tool 
that allows you to perform transformations on bosh manifests.

## Development Setup

Make sure you have a working Go toolchain and have installed [glide](http://glide.sh/).

```sh
go get -d github.com/enaml-ops/omg-transform
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
