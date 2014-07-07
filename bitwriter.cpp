// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#include "bitwriter.h"

#include "llvm-c/Core.h"
#include "llvm-c/Support.h"
#include "llvm/Bitcode/ReaderWriter.h"
#include "llvm/IR/Module.h"
#include "llvm/Support/MemoryBuffer.h"
#include "llvm/Support/raw_ostream.h"

using namespace llvm;

LLVMMemoryBufferRef gollvm_WriteBitcodeToMemoryBuffer(LLVMModuleRef M) {
  std::string Data;
  raw_string_ostream OS(Data);

  WriteBitcodeToFile(unwrap(M), OS);
  return wrap(MemoryBuffer::getMemBufferCopy(OS.str()));
}
