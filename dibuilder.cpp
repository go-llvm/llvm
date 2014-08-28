// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#include "dibuilder.h"

#include <llvm/IR/Module.h>
#include <llvm/IR/DIBuilder.h>

#include <iostream>

using namespace llvm;

namespace {
template <typename T> T unwrapDI(LLVMValueRef v) {
  return v ? T(unwrap<MDNode>(v)) : T();
}
}

DEFINE_SIMPLE_CONVERSION_FUNCTIONS(DIBuilder, LLVMDIBuilderRef)

LLVMDIBuilderRef NewDIBuilder(LLVMModuleRef mref) {
  Module *m = unwrap(mref);
  return wrap(new DIBuilder(*m));
}

void DIBuilderDestroy(LLVMDIBuilderRef dref) {
  DIBuilder *d = unwrap(dref);
  delete d;
}

void DIBuilderFinalize(LLVMDIBuilderRef dref) { unwrap(dref)->finalize(); }

LLVMValueRef DIBuilderCreateCompileUnit(LLVMDIBuilderRef dref, unsigned lang,
                                        const char *file, const char *dir,
                                        const char *producer, int optimized,
                                        const char *flags,
                                        unsigned runtimeVersion) {
  DIBuilder *d = unwrap(dref);
  DICompileUnit cu = d->createCompileUnit(lang, StringRef(file), StringRef(dir),
                                          StringRef(producer), optimized,
                                          StringRef(flags), runtimeVersion);
  return wrap(cu);
}

LLVMValueRef DIBuilderCreateFile(LLVMDIBuilderRef dref, const char *file,
                                 const char *dir) {
  DIBuilder *d = unwrap(dref);
  DIFile f = d->createFile(StringRef(file), StringRef(dir));
  return wrap(f);
}

LLVMValueRef DIBuilderCreateLexicalBlock(LLVMDIBuilderRef dref,
                                         LLVMValueRef diScope,
                                         LLVMValueRef diFile, unsigned line,
                                         unsigned column) {
  DIBuilder *d = unwrap(dref);
  DILexicalBlock lb = d->createLexicalBlock(unwrapDI<DIDescriptor>(diScope),
                                            unwrapDI<DIFile>(diFile), line,
                                            column);
  return wrap(lb);
}

LLVMValueRef DIBuilderCreateLexicalBlockFile(LLVMDIBuilderRef dref,
                                             LLVMValueRef diScope,
                                             LLVMValueRef diFile,
                                             unsigned discriminator) {
  DIBuilder *d = unwrap(dref);
  DILexicalBlockFile lbf = d->createLexicalBlockFile(
      unwrapDI<DIDescriptor>(diScope), unwrapDI<DIFile>(diFile), discriminator);
  return wrap(lbf);
}

LLVMValueRef DIBuilderCreateFunction(LLVMDIBuilderRef dref,
                                     LLVMValueRef diScope, const char *name,
                                     const char *linkageName,
                                     LLVMValueRef diFile, unsigned line,
                                     LLVMValueRef diCompositeType,
                                     int isLocalToUnit, int isDefinition,
                                     unsigned scopeLine, unsigned flags,
                                     int isOptimized, LLVMValueRef function) {
  DIBuilder *d = unwrap(dref);
  DISubprogram sp = d->createFunction(
      unwrapDI<DIDescriptor>(diScope), StringRef(name), StringRef(linkageName),
      unwrapDI<DIFile>(diFile), line,
      unwrapDI<DICompositeType>(diCompositeType), isLocalToUnit, isDefinition,
      scopeLine, flags, isOptimized, unwrap<Function>(function));
  return wrap(sp);
}

LLVMValueRef DIBuilderCreateLocalVariable(LLVMDIBuilderRef dref, unsigned tag,
                                          LLVMValueRef diScope,
                                          const char *name, LLVMValueRef diFile,
                                          unsigned line, LLVMValueRef diType,
                                          int alwaysPreserve, unsigned flags,
                                          unsigned argNo) {
  DIBuilder *d = unwrap(dref);
  DIVariable v = d->createLocalVariable(
      tag, unwrapDI<DIDescriptor>(diScope), StringRef(name),
      unwrapDI<DIFile>(diFile), line, unwrapDI<DIType>(diType), alwaysPreserve,
      flags, argNo);
  return wrap(v);
}

LLVMValueRef DIBuilderCreateBasicType(LLVMDIBuilderRef dref, const char *name,
                                      uint64_t sizeInBits, uint64_t alignInBits,
                                      unsigned encoding) {
  DIBuilder *d = unwrap(dref);
  DIBasicType t =
      d->createBasicType(StringRef(name), sizeInBits, alignInBits, encoding);
  return wrap(t);
}

LLVMValueRef DIBuilderCreatePointerType(LLVMDIBuilderRef dref,
                                        LLVMValueRef pointeeType,
                                        uint64_t sizeInBits,
                                        uint64_t alignInBits,
                                        const char *name) {
  DIBuilder *d = unwrap(dref);
  DIDerivedType t = d->createPointerType(
      unwrapDI<DIType>(pointeeType), sizeInBits, alignInBits, StringRef(name));
  return wrap(t);
}

LLVMValueRef DIBuilderCreateSubroutineType(LLVMDIBuilderRef dref,
                                           LLVMValueRef diFile,
                                           LLVMValueRef diParameterTypes) {
  DIBuilder *d = unwrap(dref);
  DICompositeType ct = d->createSubroutineType(
      unwrapDI<DIFile>(diFile), unwrapDI<DITypeArray>(diParameterTypes));
  return wrap(ct);
}

LLVMValueRef DIBuilderCreateStructType(LLVMDIBuilderRef dref,
                                       LLVMValueRef diScope, const char *name,
                                       LLVMValueRef diFile, unsigned line,
                                       uint64_t sizeInBits,
                                       uint64_t alignInBits, unsigned flags,
                                       LLVMValueRef diDerivedFrom,
                                       LLVMValueRef diElementTypes) {
  DIBuilder *d = unwrap(dref);
  DICompositeType ct = d->createStructType(
      unwrapDI<DIDescriptor>(diScope), StringRef(name),
      unwrapDI<DIFile>(diFile), line, sizeInBits, alignInBits, flags,
      unwrapDI<DIType>(diDerivedFrom), unwrapDI<DIArray>(diElementTypes));
  return wrap(ct);
}

LLVMValueRef DIBuilderCreateMemberType(LLVMDIBuilderRef dref,
                                       LLVMValueRef diScope, const char *name,
                                       LLVMValueRef diFile, unsigned line,
                                       uint64_t sizeInBits,
                                       uint64_t alignInBits,
                                       uint64_t offsetInBits, unsigned flags,
                                       LLVMValueRef diType) {
  DIBuilder *d = unwrap(dref);
  DIDerivedType dt = d->createMemberType(
      unwrapDI<DIDescriptor>(diScope), StringRef(name),
      unwrapDI<DIFile>(diFile), line, sizeInBits, alignInBits, offsetInBits,
      flags, unwrapDI<DIType>(diType));
  return wrap(dt);
}

LLVMValueRef DIBuilderCreateArrayType(LLVMDIBuilderRef dref,
                                      uint64_t sizeInBits, uint64_t alignInBits,
                                      LLVMValueRef diElementType,
                                      LLVMValueRef diSubscripts) {
  DIBuilder *d = unwrap(dref);
  DICompositeType ct = d->createArrayType(sizeInBits, alignInBits,
                                          unwrapDI<DIType>(diElementType),
                                          unwrapDI<DIArray>(diSubscripts));
  return wrap(ct);
}

LLVMValueRef DIBuilderCreateTypedef(LLVMDIBuilderRef dref, LLVMValueRef diType,
                                    const char *name, LLVMValueRef diFile,
                                    unsigned line, LLVMValueRef diContext) {
  DIBuilder *d = unwrap(dref);
  DIDerivedType dt = d->createTypedef(unwrapDI<DIType>(diType), StringRef(name),
                                      unwrapDI<DIFile>(diFile), line,
                                      unwrapDI<DIDescriptor>(diContext));
  return wrap(dt);
}

LLVMValueRef DIBuilderGetOrCreateSubrange(LLVMDIBuilderRef dref, int64_t lo,
                                          int64_t count) {
  DIBuilder *d = unwrap(dref);
  DISubrange s = d->getOrCreateSubrange(lo, count);
  return wrap(s);
}

LLVMValueRef DIBuilderGetOrCreateArray(LLVMDIBuilderRef dref,
                                       LLVMValueRef *data, size_t length) {
  DIBuilder *d = unwrap(dref);
  Value **dataValue = unwrap(data);
  ArrayRef<Value *> elements(dataValue, length);
  DIArray a = d->getOrCreateArray(elements);
  return wrap(a);
}

LLVMValueRef DIBuilderGetOrCreateTypeArray(LLVMDIBuilderRef dref,
                                           LLVMValueRef *data, size_t length) {
  DIBuilder *d = unwrap(dref);
  Value **dataValue = unwrap(data);
  ArrayRef<Value *> elements(dataValue, length);
  DITypeArray a = d->getOrCreateTypeArray(elements);
  return wrap(a);
}

LLVMValueRef DIBuilderInsertDeclareAtEnd(LLVMDIBuilderRef dref,
                                         LLVMValueRef storage,
                                         LLVMValueRef varInfo,
                                         LLVMBasicBlockRef block) {
  DIBuilder *d = unwrap(dref);
  Instruction *instr = d->insertDeclare(
      unwrap(storage), unwrapDI<DIVariable>(varInfo), unwrap(block));
  return wrap(instr);
}
