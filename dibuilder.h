// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#ifndef _GOLLVM_LLVM_DIBUILDER_H
#define _GOLLVM_LLVM_DIBUILDER_H

#include <llvm-c/Core.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct LLVMOpaqueDIBuilder *LLVMDIBuilderRef;
typedef struct LLVMOpaqueDICompileUnit *LLVMDICompileUnitRef;

LLVMDIBuilderRef NewDIBuilder(LLVMModuleRef m);

void DIBuilderDestroy(LLVMDIBuilderRef d);
void DIBuilderFinalize(LLVMDIBuilderRef d);

LLVMValueRef DIBuilderCreateCompileUnit(LLVMDIBuilderRef d, unsigned language,
                                        const char *file, const char *dir,
                                        const char *producer, int optimized,
                                        const char *flags,
                                        unsigned runtimeVersion);

LLVMValueRef DIBuilderCreateFile(LLVMDIBuilderRef d, const char *file,
                                 const char *dir);

LLVMValueRef DIBuilderCreateLexicalBlock(LLVMDIBuilderRef d,
                                         LLVMValueRef diScope,
                                         LLVMValueRef diFile, unsigned line,
                                         unsigned column,
                                         unsigned discriminator);

LLVMValueRef DIBuilderCreateLexicalBlockFile(LLVMDIBuilderRef d,
                                             LLVMValueRef diScope,
                                             LLVMValueRef diFile);

LLVMValueRef DIBuilderCreateFunction(LLVMDIBuilderRef d, LLVMValueRef diScope,
                                     const char *name, const char *linkageName,
                                     LLVMValueRef diFile, unsigned line,
                                     LLVMValueRef diCompositeType,
                                     int isLocalToUnit, int isDefinition,
                                     unsigned scopeLine, unsigned flags,
                                     int isOptimized, LLVMValueRef function);

LLVMValueRef DIBuilderCreateLocalVariable(LLVMDIBuilderRef d, unsigned tag,
                                          LLVMValueRef diScope,
                                          const char *name, LLVMValueRef diFile,
                                          unsigned line, LLVMValueRef diType,
                                          int alwaysPreserve, unsigned flags,
                                          unsigned argNo);

LLVMValueRef DIBuilderCreateBasicType(LLVMDIBuilderRef d, const char *name,
                                      uint64_t sizeInBits, uint64_t alignInBits,
                                      unsigned encoding);

LLVMValueRef DIBuilderCreatePointerType(LLVMDIBuilderRef d,
                                        LLVMValueRef pointeeType,
                                        uint64_t sizeInBits,
                                        uint64_t alignInBits, const char *name);

LLVMValueRef DIBuilderCreateSubroutineType(LLVMDIBuilderRef d,
                                           LLVMValueRef diFile,
                                           LLVMValueRef diParameterTypes);

LLVMValueRef DIBuilderCreateStructType(LLVMDIBuilderRef d, LLVMValueRef diScope,
                                       const char *name, LLVMValueRef diFile,
                                       unsigned line, uint64_t sizeInBits,
                                       uint64_t alignInBits, unsigned flags,
                                       LLVMValueRef diDerivedFrom,
                                       LLVMValueRef diElementTypes);

LLVMValueRef DIBuilderCreateMemberType(LLVMDIBuilderRef d, LLVMValueRef diScope,
                                       const char *name, LLVMValueRef diFile,
                                       unsigned line, uint64_t sizeInBits,
                                       uint64_t alignInBits,
                                       uint64_t offsetInBits, unsigned flags,
                                       LLVMValueRef diType);

LLVMValueRef DIBuilderCreateArrayType(LLVMDIBuilderRef d, uint64_t sizeInBits,
                                      uint64_t alignInBits,
                                      LLVMValueRef diElementType,
                                      LLVMValueRef diSubscripts);

LLVMValueRef DIBuilderCreateTypedef(LLVMDIBuilderRef d, LLVMValueRef diType,
                                    const char *name, LLVMValueRef diFile,
                                    unsigned line, LLVMValueRef diContext);

LLVMValueRef DIBuilderGetOrCreateSubrange(LLVMDIBuilderRef d, int64_t lo,
                                          int64_t count);

LLVMValueRef DIBuilderGetOrCreateArray(LLVMDIBuilderRef d, LLVMValueRef *data,
                                       size_t length);

LLVMValueRef DIBuilderGetOrCreateTypeArray(LLVMDIBuilderRef d,
                                           LLVMValueRef *data, size_t length);

LLVMValueRef DIBuilderInsertDeclareAtEnd(LLVMDIBuilderRef d,
                                         LLVMValueRef storage,
                                         LLVMValueRef varInfo,
                                         LLVMBasicBlockRef block);

#ifdef __cplusplus
} // extern "C"
#endif

#endif
