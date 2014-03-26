// +build llvmsvn llvm3.2 llvm3.3

package llvm

/*
#include <llvm-c/Linker.h>
#include <stdlib.h>
*/
import "C"
import "errors"

type LinkerMode C.LLVMLinkerMode

const (
	LinkerDestroySource  = C.LLVMLinkerDestroySource
	LinkerPreserveSource = C.LLVMLinkerPreserveSource
)

func LinkModules(Dest, Src Module, Mode LinkerMode) error {
	var cmsg *C.char
	failed := C.LLVMLinkModules(Dest.C, Src.C, C.LLVMLinkerMode(Mode), &cmsg)
	if failed != 0 {
		err := errors.New(C.GoString(cmsg))
		C.LLVMDisposeMessage(cmsg)
		return err
	}
	return nil
}
