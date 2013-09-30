#include <llvm/DebugInfo.h>
#include <llvm/Intrinsics.h>

typedef void* LLVMValueRef;
extern "C" LLVMValueRef getDbgDeclare(llvm::Module* module) {
	llvm::Intrinsic::ID id = llvm::Intrinsic::dbg_declare;
	// TODO: why on earth is the +1 needed??
	id++;
	llvm::Function* DeclareFn = llvm::Intrinsic::getDeclaration(module, id);
	return DeclareFn;
}
