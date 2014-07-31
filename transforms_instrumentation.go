package llvm

/*
#include "transforms_instrumentation.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func (pm PassManager) AddAddressSanitizerFunctionPass() {
	C.gollvm_AddAddressSanitizerFunctionPass(pm.C)
}

func (pm PassManager) AddAddressSanitizerModulePass() {
	C.gollvm_AddAddressSanitizerModulePass(pm.C)
}

func (pm PassManager) AddThreadSanitizerPass() {
	C.gollvm_AddThreadSanitizerPass(pm.C)
}

func (pm PassManager) AddMemorySanitizerPass() {
	C.gollvm_AddMemorySanitizerPass(pm.C)
}

func (pm PassManager) AddDataFlowSanitizerPass(abilist string) {
	cabilist := C.CString(abilist)
	C.gollvm_AddDataFlowSanitizerPass(pm.C, cabilist)
	C.free(unsafe.Pointer(cabilist))
}
