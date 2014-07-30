// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#ifndef _GOLLVM_LLVM_CORE_H
#define _GOLLVM_LLVM_CORE_H

#include "llvm-c/Core.h"
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

void gollvm_AddFunctionAttr(LLVMValueRef Fn, uint64_t PA);
uint64_t gollvm_GetFunctionAttr(LLVMValueRef Fn);
void gollvm_RemoveFunctionAttr(LLVMValueRef Fn, uint64_t PA);

#ifdef __cplusplus
}
#endif

#endif
