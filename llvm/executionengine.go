package llvm

/*
#include <llvm-c/ExecutionEngine.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "errors"

func LinkInJIT()         { C.LLVMLinkInJIT() }
func LinkInInterpreter() { C.LLVMLinkInInterpreter() }

type (
	GenericValue struct {
		C C.LLVMGenericValueRef
	}
	ExecutionEngine struct {
		C C.LLVMExecutionEngineRef
	}
)

// helpers
func llvmGenericValueRefPtr(t *GenericValue) *C.LLVMGenericValueRef {
	return (*C.LLVMGenericValueRef)(unsafe.Pointer(t))
}

//-------------------------------------------------------------------------
// llvm.GenericValue
//-------------------------------------------------------------------------

func NewGenericValueFromInt(t Type, n uint64, signed bool) (g GenericValue) {
	g.C = C.LLVMCreateGenericValueOfInt(t.C, C.ulonglong(n), boolToLLVMBool(signed))
	return
}
func NewGenericValueFromPointer(p unsafe.Pointer) (g GenericValue) {
	g.C = C.LLVMCreateGenericValueOfPointer(p)
	return
}
func NewGenericValueFromFloat(t Type, n float64) (g GenericValue) {
	g.C = C.LLVMCreateGenericValueOfFloat(t.C, C.double(n))
	return
}
func (g GenericValue) IntWidth() int { return int(C.LLVMGenericValueIntWidth(g.C)) }
func (g GenericValue) Int(signed bool) uint64 {
	return uint64(C.LLVMGenericValueToInt(g.C, boolToLLVMBool(signed)))
}
func (g GenericValue) Float(t Type) float64 {
	return float64(C.LLVMGenericValueToFloat(t.C, g.C))
}
func (g GenericValue) Pointer() unsafe.Pointer {
	return C.LLVMGenericValueToPointer(g.C)
}
func (g GenericValue) Dispose() { C.LLVMDisposeGenericValue(g.C) }

//-------------------------------------------------------------------------
// llvm.ExecutionEngine
//-------------------------------------------------------------------------

func NewExecutionEngine(m Module) (ee ExecutionEngine, err error) {
	var cmsg *C.char
	fail := C.LLVMCreateExecutionEngineForModule(&ee.C, m.C, &cmsg)
	if fail != 0 {
		ee.C = nil
		err = errors.New(C.GoString(cmsg))
		C.LLVMDisposeMessage(cmsg)
	} else {
		err = nil
	}
	return
}

func NewInterpreter(m Module) (ee ExecutionEngine, err error) {
	var cmsg *C.char
	fail := C.LLVMCreateInterpreterForModule(&ee.C, m.C, &cmsg)
	if fail != 0 {
		ee.C = nil
		err = errors.New(C.GoString(cmsg))
		C.LLVMDisposeMessage(cmsg)
	} else {
		err = nil
	}
	return
}
func NewJITCompiler(m Module, optLevel int) (ee ExecutionEngine, err error) {
	var cmsg *C.char
	fail := C.LLVMCreateJITCompilerForModule(&ee.C, m.C, C.unsigned(optLevel), &cmsg)
	if fail != 0 {
		ee.C = nil
		err = errors.New(C.GoString(cmsg))
		C.LLVMDisposeMessage(cmsg)
	} else {
		err = nil
	}
	return
}

// XXX: Don't port deprecated
// Deprecated: Use LLVMCreateExecutionEngineForModule instead.
//LLVMBool LLVMCreateExecutionEngine(LLVMExecutionEngineRef *OutEE,
//                                   LLVMModuleProviderRef MP,
//                                   char **OutError);

// Deprecated: Use LLVMCreateInterpreterForModule instead.
//LLVMBool LLVMCreateInterpreter(LLVMExecutionEngineRef *OutInterp,
//                               LLVMModuleProviderRef MP,
//                               char **OutError);

// Deprecated: Use LLVMCreateJITCompilerForModule instead.
//LLVMBool LLVMCreateJITCompiler(LLVMExecutionEngineRef *OutJIT,
//                               LLVMModuleProviderRef MP,
//                               unsigned OptLevel,
//                               char **OutError);

func (ee ExecutionEngine) Dispose()               { C.LLVMDisposeExecutionEngine(ee.C) }
func (ee ExecutionEngine) RunStaticConstructors() { C.LLVMRunStaticConstructors(ee.C) }
func (ee ExecutionEngine) RunStaticDestructors()  { C.LLVMRunStaticDestructors(ee.C) }

// TODO(nsf): figure out how to convert that stuff from Go's "os.Argv"
//int LLVMRunFunctionAsMain(LLVMExecutionEngineRef EE, LLVMValueRef F,
//                          unsigned ArgC, const char * const *ArgV,
//                          const char * const *EnvP);

func (ee ExecutionEngine) RunFunction(f Value, args []GenericValue) (g GenericValue) {
	nargs := len(args)
	var argptr *GenericValue
	if nargs > 0 {
		argptr = &args[0]
	}
	g.C = C.LLVMRunFunction(ee.C, f.C,
		C.unsigned(nargs), llvmGenericValueRefPtr(argptr))
	return
}

func (ee ExecutionEngine) FreeMachineCodeForFunction(f Value) {
	C.LLVMFreeMachineCodeForFunction(ee.C, f.C)
}
func (ee ExecutionEngine) AddModule(m Module) { C.LLVMAddModule(ee.C, m.C) }

// XXX(nsf): Don't port deprecated
// Deprecated: Use LLVMAddModule instead.
//void LLVMAddModuleProvider(LLVMExecutionEngineRef EE, LLVMModuleProviderRef MP);

func (ee ExecutionEngine) RemoveModule(m Module) {
	var modtmp C.LLVMModuleRef
	C.LLVMRemoveModule(ee.C, m.C, &modtmp, nil)
}

// XXX(nsf): Don't port deprecated
// Deprecated: Use LLVMRemoveModule instead.
//LLVMBool LLVMRemoveModuleProvider(LLVMExecutionEngineRef EE,
//                                  LLVMModuleProviderRef MP,
//                                  LLVMModuleRef *OutMod, char **OutError);

func (ee ExecutionEngine) FindFunction(name string) (f Value) {
	cname := C.CString(name)
	C.LLVMFindFunction(ee.C, cname, &f.C)
	C.free(unsafe.Pointer(cname))
	return
}

func (ee ExecutionEngine) RecompileAndRelinkFunction(f Value) unsafe.Pointer {
	return C.LLVMRecompileAndRelinkFunction(ee.C, f.C)
}

func (ee ExecutionEngine) TargetData() (td TargetData) {
	td.C = C.LLVMGetExecutionEngineTargetData(ee.C)
	return
}

func (ee ExecutionEngine) AddGlobalMapping(global Value, addr unsafe.Pointer) {
	C.LLVMAddGlobalMapping(ee.C, global.C, addr)
}

func (ee ExecutionEngine) PointerToGlobal(global Value) unsafe.Pointer {
	return C.LLVMGetPointerToGlobal(ee.C, global.C)
}

// vim: set ft=go:
