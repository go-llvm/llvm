package llvm_test

import (
	"testing"

	"github.com/axw/gollvm/llvm"
)

func TestStringRecursiveType(t *testing.T) {
	//mod := llvm.NewModule("fac_module")
	ctx := llvm.NewContext()
	defer ctx.Dispose()
	s := ctx.StructCreateNamed("recursive")
	s.StructSetBody([]llvm.Type{s, s}, false)
	if str := s.String(); str != "%recursive: StructType(%recursive, %recursive)" {
		t.Errorf("incorrect string result %q", str)
	}
}
