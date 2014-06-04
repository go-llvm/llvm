// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#include "support.h"

#include "llvm/Support/CommandLine.h"
#include "llvm/Support/DynamicLibrary.h"
#include <stdlib.h>
#include <string.h>

void gollvm_LoadLibraryPermanently(const char *Filename, char **ErrMsg) {
  std::string ErrMsgStr;
  if (llvm::sys::DynamicLibrary::LoadLibraryPermanently(Filename, &ErrMsgStr)) {
    *ErrMsg = static_cast<char *>(malloc(ErrMsgStr.size() + 1));
    memcpy(static_cast<void *>(*ErrMsg),
           static_cast<const void *>(ErrMsgStr.c_str()), ErrMsgStr.size() + 1);
  }
}

void gollvm_ParseCommandLineOptions(int argc, const char *const *argv,
                                    const char *Overview) {
  llvm::cl::ParseCommandLineOptions(argc, argv, Overview);
}
