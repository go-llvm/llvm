// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#include "core.h"

#include "llvm/IR/Attributes.h"
#include "llvm/IR/Function.h"

using namespace llvm;

void gollvm_AddFunctionAttr(LLVMValueRef Fn, uint64_t PA) {
  Function *Func = unwrap<Function>(Fn);
  const AttributeSet PAL = Func->getAttributes();
  AttrBuilder B(PA);
  const AttributeSet PALnew =
    PAL.addAttributes(Func->getContext(), AttributeSet::FunctionIndex,
                      AttributeSet::get(Func->getContext(),
                                        AttributeSet::FunctionIndex, B));
  Func->setAttributes(PALnew);
}

uint64_t gollvm_GetFunctionAttr(LLVMValueRef Fn) {
  Function *Func = unwrap<Function>(Fn);
  const AttributeSet PAL = Func->getAttributes();
  return PAL.Raw(AttributeSet::FunctionIndex);
}

void gollvm_RemoveFunctionAttr(LLVMValueRef Fn, uint64_t PA) {
  Function *Func = unwrap<Function>(Fn);
  const AttributeSet PAL = Func->getAttributes();
  AttrBuilder B(PA);
  const AttributeSet PALnew =
    PAL.removeAttributes(Func->getContext(), AttributeSet::FunctionIndex,
                         AttributeSet::get(Func->getContext(),
                                           AttributeSet::FunctionIndex, B));
  Func->setAttributes(PALnew);
}
