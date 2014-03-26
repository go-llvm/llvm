package llvm

/*
#include <llvm-c/BitReader.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

// ParseBitcodeFile parses the LLVM IR (bitcode) in the file with the
// specified name, and returns a new LLVM module.
func ParseBitcodeFile(name string) (Module, error) {
	var buf C.LLVMMemoryBufferRef
	var errmsg *C.char
	var cfilename *C.char = C.CString(name)
	result := C.LLVMCreateMemoryBufferWithContentsOfFile(cfilename, &buf, &errmsg)
	C.free(unsafe.Pointer(cfilename))
	if result != 0 {
		err := errors.New(C.GoString(errmsg))
		C.free(unsafe.Pointer(errmsg))
		return Module{}, err
	}
	defer C.LLVMDisposeMemoryBuffer(buf)

	var m Module
	if C.LLVMParseBitcode(buf, &m.C, &errmsg) == 0 {
		return m, nil
	}

	err := errors.New(C.GoString(errmsg))
	C.free(unsafe.Pointer(errmsg))
	return Module{nil}, err
}
