// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#include "llvm-c/Core.h"

#ifdef __cplusplus
extern "C" {
#endif

void gollvm_AddAddressSanitizerFunctionPass(LLVMPassManagerRef PM);
void gollvm_AddAddressSanitizerModulePass(LLVMPassManagerRef PM);
void gollvm_AddThreadSanitizerPass(LLVMPassManagerRef PM);
void gollvm_AddMemorySanitizerPass(LLVMPassManagerRef PM);
void gollvm_AddDataFlowSanitizerPass(LLVMPassManagerRef PM,
                                     const char *ABIListFile);

#ifdef __cplusplus
}
#endif
