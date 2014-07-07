// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#include "llvm-c/Core.h"
#include "llvm-c/Support.h"

#ifdef __cplusplus
extern "C" {
#endif

LLVMMemoryBufferRef gollvm_WriteBitcodeToMemoryBuffer(LLVMModuleRef M);

#ifdef __cplusplus
}
#endif
