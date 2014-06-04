// Copyright 2014 The go-llvm Authors.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

#ifndef _GOLLVM_LLVM_SUPPORT_H
#define _GOLLVM_LLVM_SUPPORT_H

#ifdef __cplusplus
extern "C" {
#endif

void gollvm_LoadLibraryPermanently(const char *Filename, char **ErrMsg);

void gollvm_ParseCommandLineOptions(int argc, const char *const *argv,
                                    const char *Overview);

#ifdef __cplusplus
}
#endif

#endif
