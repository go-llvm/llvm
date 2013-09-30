package llvm

/*
#include <llvm-c/Core.h>

extern LLVMValueRef getDbgDeclare(LLVMModuleRef);
*/
import "C"

import "fmt"

func (b Builder) InsertDeclare(module Module, storage Value, md Value) Value {
	nf := Value{C.getDbgDeclare(module.C)}
	if nf.IsAFunction().IsNil() || nf.Name() != "llvm.dbg.declare" {
		panic(fmt.Sprintf("Wanted llvm.dbg.declare but got: %s", nf.Name()))
	}
	return b.CreateCall(nf, []Value{storage, md}, "")
}
