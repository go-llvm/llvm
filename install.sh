#!/bin/sh
export CGO_CFLAGS="`llvm-config --cflags`"
export CGO_LDFLAGS="`llvm-config --ldflags` -Wl,-L`llvm-config --libdir` `llvm-config --libs interpreter all-targets bitwriter ipo engine`"
go install $* github.com/axw/gollvm/llvm
