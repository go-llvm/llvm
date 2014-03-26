GoLLVM
------

[LLVM](http://llvm.org) bindings for [The Go Programming Language](http://golang.org).

Prerequisites
-------------

* LLVM 3.3+. LLVM must have been built with shared libraries enabled.
* Go 1.2+.

The author has only built and tested with Linux, but there is no particular
reason why GoLLVM should not work with other operating systems.

Installation
------------

To install, run the following (assuming you have curl and Go installed):

    curl https://raw.github.com/go-llvm/master/install.sh | sh

Alternatively, you can use `go get` directly, but you must then set the
CGO\_CFLAGS and CGO\_LDFLAGS environment variables:

    $ export CGO_CFLAGS=`llvm-config --cflags`
    $ export CGO_LDFLAGS="`llvm-config --ldflags` -Wl,-L`llvm-config --libdir` -lLLVM-`llvm-config --version`"
    $ go get github.com/go-llvm/llvm

