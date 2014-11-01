GoLLVM
------

**This project has moved to [llvm.org](http://llvm.org). Any contributions or bug reports should be sent there.
Please refer to the [Go bindings readme](http://llvm.org/svn/llvm-project/llvm/trunk/bindings/go/README.txt)
for more information.**

[LLVM](http://llvm.org) bindings for [The Go Programming Language](http://golang.org).

[![Build Status](https://travis-ci.org/go-llvm/llvm.svg)](https://travis-ci.org/go-llvm/llvm)

Prerequisites
-------------

* Go 1.2+.
* CMake (to build LLVM).

The author has only built and tested with Linux, but there is no particular
reason why GoLLVM should not work with other operating systems.

Building LLVM
-------------

The development version of GoLLVM is tied to a specific revision of LLVM SVN
trunk. We also make releases which are synchronized with LLVM releases. Those
releases are available in branches named vMN, where M is the major number of
the LLVM release and N is the minor number.

The script `update\_llvm.sh` in this directory can be used to build the
required revision of LLVM and prepare it to be used by GoLLVM. If you receive
an error message from `go build` like this:

    ./analysis.go:4:84: fatal error: llvm-c/Analysis.h: No such file or directory
     #include <llvm-c/Analysis.h> // If you are getting an error here read README.md

or like this:

    ./llvm_dep.go:5: undefined: run_update_llvm_sh_to_get_revision_NNNNNN

it means that LLVM needs to be built or updated by running the script.

    $ ./update_llvm.sh

Any command line arguments supplied to the script are passed to LLVM's CMake
build system. A good set of arguments to use during development are:

    $ ./update_llvm.sh -DCMAKE_BUILD_TYPE=Debug -DLLVM_TARGETS_TO_BUILD=host -DBUILD_SHARED_LIBS=ON

Note that CMake keeps a cache of build settings so once you have built
LLVM there is no need to pass these arguments again after updating.

Alternatively, you can build LLVM yourself or use a version provided by
your distribution, but you must then set the CGO\_CFLAGS and CGO\_LDFLAGS
environment variables:

    $ export CGO_CFLAGS=`llvm-config --cflags`
    $ export CGO_LDFLAGS="`llvm-config --ldflags` -Wl,-L`llvm-config --libdir` -lLLVM-`llvm-config --version`"
    $ go get -tags byollvm github.com/go-llvm/llvm

