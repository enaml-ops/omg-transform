omg-transform
=============

An `enaml` tool that allows you to perform transformations on bosh manifests.

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

The input manifest can also be read from a file by using the `-f` flag.

_Example from stdin:_

`omg-cli deploy-product --print-manifest cloudfoundry-plugin-linux | omg-transform <TRANSFORM> [flags...]`

_From file:_

`omg-transform -f cf.yml <TRANSFORM> [flags...]`

## Transformations

### `change-network`: Change Instance Group Network

Move instance group to new network:

`omg-cli deploy-product --print-manifest cloudfoundry-plugin | omg-transform change-network --move diego_cell:apps-network`

### `clone`: Clone Instance Group

Clone the `router` instance group, naming the new instance group `router2`.

TODO: customize # instances, other ig fields

`omg-transform -f cf.yml clone -g router -o router2`

