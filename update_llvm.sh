#!/bin/sh -xe

llvm_components="\
all-targets \
analysis \
asmparser \
asmprinter \
bitreader \
bitwriter \
codegen \
core \
debuginfo \
executionengine \
interpreter \
ipo \
irreader \
jit \
linker \
mc \
scalaropts \
support \
target \
"

gollvmdir=$(dirname "$0")

llvmrev=$(grep run_update_llvm_sh_to_get_revision_ $gollvmdir/llvm_dep.go | \
          sed -E -e 's/.*run_update_llvm_sh_to_get_revision_([0-9]+.*)/\1/g')

workdir=$gollvmdir/workdir
llvmdir=$workdir/llvm
llvm_builddir=$workdir/llvm_build

mkdir -p $workdir
svn co -r $llvmrev https://llvm.org/svn/llvm-project/llvm/trunk $llvmdir
mkdir -p $llvm_builddir

cmake_flags="../llvm $@"
llvm_config="$llvm_builddir/bin/llvm-config"

if test -n "`which ninja`" ; then
  # If Ninja is available, we can speed up the build by building only the
  # required subset of LLVM.
  (cd $llvm_builddir && cmake -G Ninja $cmake_flags)
  ninja -C $llvm_builddir llvm-config
  llvm_buildtargets="$($llvm_config --libs $llvm_components | sed -e 's/-l//g')"
  ninja -C $llvm_builddir $llvm_buildtargets
else
  (cd $llvm_builddir && cmake $cmake_flags)
  make -C $llvm_builddir -j4
fi

llvm_cflags="$($llvm_config --cppflags)"
if [ $(uname) == "Darwin" ]; then
  # OS X doesn't like -rpath with cgo. See:
  # https://code.google.com/p/go/issues/detail?id=7293
  llvm_ldflags="$($llvm_config --ldflags)                                     $($llvm_config --libs $llvm_components) $($llvm_config --system-libs)"
else
  llvm_ldflags="$($llvm_config --ldflags) -Wl,-rpath,$($llvm_config --libdir) $($llvm_config --libs $llvm_components) $($llvm_config --system-libs)"
fi
sed -e "s#@LLVM_REVISION@#$llvmrev#g; s#@LLVM_CFLAGS@#$llvm_cflags#g; \
        s#@LLVM_LDFLAGS@#$llvm_ldflags#g" $gollvmdir/llvm_config.go.in > \
  $gollvmdir/llvm_config.go
