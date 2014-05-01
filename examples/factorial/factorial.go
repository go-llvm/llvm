// Taken from: http://npcontemplation.blogspot.com/2008/06/secret-of-llvm-c-bindings.html
package main

import "fmt"

import "github.com/go-llvm/llvm"

func test() {
	llvm.LinkInJIT()
	llvm.InitializeNativeTarget()

	mod := llvm.NewModule("fac_module")

	// don't do that, because ExecutionEngine takes ownership over module
	//defer mod.Dispose()

	fac_args := []llvm.Type{llvm.Int32Type()}
	fac_type := llvm.FunctionType(llvm.Int32Type(), fac_args, false)
	fac := llvm.AddFunction(mod, "fac", fac_type)
	fac.SetFunctionCallConv(llvm.CCallConv)
	n := fac.Param(0)

	entry := llvm.AddBasicBlock(fac, "entry")
	iftrue := llvm.AddBasicBlock(fac, "iftrue")
	iffalse := llvm.AddBasicBlock(fac, "iffalse")
	end := llvm.AddBasicBlock(fac, "end")

	builder := llvm.NewBuilder()
	defer builder.Dispose()

	builder.SetInsertPointAtEnd(entry)
	If := builder.CreateICmp(llvm.IntEQ, n, llvm.ConstInt(llvm.Int32Type(), 0, false), "cmptmp")
	builder.CreateCondBr(If, iftrue, iffalse)

	builder.SetInsertPointAtEnd(iftrue)
	res_iftrue := llvm.ConstInt(llvm.Int32Type(), 1, false)
	builder.CreateBr(end)

	builder.SetInsertPointAtEnd(iffalse)
	n_minus := builder.CreateSub(n, llvm.ConstInt(llvm.Int32Type(), 1, false), "subtmp")
	call_fac_args := []llvm.Value{n_minus}
	call_fac := builder.CreateCall(fac, call_fac_args, "calltmp")
	res_iffalse := builder.CreateMul(n, call_fac, "multmp")
	builder.CreateBr(end)

	builder.SetInsertPointAtEnd(end)
	res := builder.CreatePHI(llvm.Int32Type(), "result")
	phi_vals := []llvm.Value{res_iftrue, res_iffalse}
	phi_blocks := []llvm.BasicBlock{iftrue, iffalse}
	res.AddIncoming(phi_vals, phi_blocks)
	builder.CreateRet(res)

	err := llvm.VerifyModule(mod, llvm.ReturnStatusAction)
	if err != nil {
		fmt.Println(err)
		return
	}

	engine, err := llvm.NewJITCompiler(mod, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Dispose()

	pass := llvm.NewPassManager()
	defer pass.Dispose()

	pass.Add(engine.TargetData())
	pass.AddConstantPropagationPass()
	pass.AddInstructionCombiningPass()
	pass.AddPromoteMemoryToRegisterPass()
	pass.AddGVNPass()
	pass.AddCFGSimplificationPass()
	pass.Run(mod)

	mod.Dump()

	exec_args := []llvm.GenericValue{llvm.NewGenericValueFromInt(llvm.Int32Type(), 10, false)}
	exec_res := engine.RunFunction(fac, exec_args)
	fmt.Println("-----------------------------------------")
	fmt.Println("Running fac(10) with JIT...")
	fmt.Printf("Result: %d\n", exec_res.Int(false))
}

func main() {
	test()
	fmt.Println("DONE")
}
