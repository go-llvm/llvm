package llvm

/*
#include <llvm-c/BitWriter.h>
#include "bitwriter.h"
#include <stdlib.h>
*/
import "C"
import "os"
import "errors"

var writeBitcodeToFileErr = errors.New("Failed to write bitcode to file")

func WriteBitcodeToFile(m Module, file *os.File) error {
	fail := C.LLVMWriteBitcodeToFD(m.C, C.int(file.Fd()), C.int(0), C.int(0))
	if fail != 0 {
		return writeBitcodeToFileErr
	}
	return nil
}

func WriteBitcodeToMemoryBuffer(m Module) MemoryBuffer {
	mb := C.gollvm_WriteBitcodeToMemoryBuffer(m.C)
	return MemoryBuffer{mb}
}

// TODO(nsf): Figure out way how to make it work with io.Writer
