pdrest
======

Installation
------------

pdrest binaries can be built from the go source or downloaded for specific architectures from
[releases](https://github.com/triddell/pdrest/releases).

To clone and compile a binary from source:

```bash
$ git clone https://github.com/triddell/pdrest.git $GOPATH/src/github.com/triddell/pdrest
$ cd $GOPATH/src/github.com/triddell/pdrest
$ go install
```

The `go install` command will create the binary at `$GOBIN/pdrest`.

Or, you can build a binary and copy it wherever you'd like with `go build`:

```bash
$ git clone https://github.com/triddell/pdrest.git $GOPATH/src/github.com/triddell/pdrest
$ cd $GOPATH/src/github.com/triddell/pdrest
$ go build pdrest.go
$ mv pdrest /to/some/directory
```
