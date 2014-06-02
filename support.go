// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package llvm

// #include "support.h"
// #include <stdlib.h>
import "C"

import (
	"errors"
	"unsafe"
)

// Loads a dynamic library such that it may be used as an LLVM plugin.
// See llvm::sys::DynamicLibrary::LoadLibraryPermanently.
func LoadLibraryPermanently(lib string) error {
	var errstr *C.char
	libstr := C.CString(lib)
	C.gollvm_LoadLibraryPermanently(libstr, &errstr)
	C.free(unsafe.Pointer(libstr))
	if errstr != nil {
		err := errors.New(C.GoString(errstr))
		C.free(unsafe.Pointer(errstr))
		return err
	}
	return nil
}

// Parse the given arguments using the LLVM command line parser.
// See llvm::cl::ParseCommandLineOptions.
func ParseCommandLineOptions(args []string, overview string) {
	argstrs := make([]*C.char, len(args))
	for i, arg := range args {
		argstrs[i] = C.CString(arg)
	}
	overviewstr := C.CString(overview)
	C.gollvm_ParseCommandLineOptions(C.int(len(args)), &argstrs[0], overviewstr)
	for _, argstr := range argstrs {
		C.free(unsafe.Pointer(argstr))
	}
	C.free(unsafe.Pointer(overviewstr))
}
