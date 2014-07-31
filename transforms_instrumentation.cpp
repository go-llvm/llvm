// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#include "transforms_instrumentation.h"

#include "llvm-c/Core.h"
#include "llvm/IR/Module.h"
#include "llvm/PassManager.h"
#include "llvm/Transforms/Instrumentation.h"

using namespace llvm;

void gollvm_AddAddressSanitizerFunctionPass(LLVMPassManagerRef PM) {
  unwrap(PM)->add(createAddressSanitizerFunctionPass());
}

void gollvm_AddAddressSanitizerModulePass(LLVMPassManagerRef PM) {
  unwrap(PM)->add(createAddressSanitizerModulePass());
}

void gollvm_AddThreadSanitizerPass(LLVMPassManagerRef PM) {
  unwrap(PM)->add(createThreadSanitizerPass());
}

void gollvm_AddMemorySanitizerPass(LLVMPassManagerRef PM) {
  unwrap(PM)->add(createMemorySanitizerPass());
}

void gollvm_AddDataFlowSanitizerPass(LLVMPassManagerRef PM,
                                     const char *ABIListFile) {
  unwrap(PM)->add(createDataFlowSanitizerPass(ABIListFile));
}
