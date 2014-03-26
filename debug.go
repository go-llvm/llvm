/*
Copyright (c) 2011, 2012 Andrew Wilkins <axwalk@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package llvm

import (
	"path"
	"reflect"
)

///////////////////////////////////////////////////////////////////////////////
// Common types and constants.

const (
	LLVMDebugVersion = (12 << 16)
)

type DwarfTag uint32

const (
	DW_TAG_lexical_block   DwarfTag = 0x0b
	DW_TAG_compile_unit    DwarfTag = 0x11
	DW_TAG_variable        DwarfTag = 0x34
	DW_TAG_base_type       DwarfTag = 0x24
	DW_TAG_pointer_type    DwarfTag = 0x0F
	DW_TAG_structure_type  DwarfTag = 0x13
	DW_TAG_subroutine_type DwarfTag = 0x15
	DW_TAG_file_type       DwarfTag = 0x29
	DW_TAG_subprogram      DwarfTag = 0x2E
	DW_TAG_auto_variable   DwarfTag = 0x100
	DW_TAG_arg_variable    DwarfTag = 0x101
)

const (
	FlagPrivate = 1 << iota
	FlagProtected
	FlagFwdDecl
	FlagAppleBlock
	FlagBlockByrefStruct
	FlagVirtual
	FlagArtificial
	FlagExplicit
	FlagPrototyped
	FlagObjcClassComplete
	FlagObjectPointer
	FlagVector
	FlagStaticMember
	FlagIndirectVariable
)

type DwarfLang uint32

const (
	// http://dwarfstd.org/ShowIssue.php?issue=101014.1&type=open
	DW_LANG_Go DwarfLang = 0x0016
)

type DwarfTypeEncoding uint32

const (
	DW_ATE_address         DwarfTypeEncoding = 0x01
	DW_ATE_boolean         DwarfTypeEncoding = 0x02
	DW_ATE_complex_float   DwarfTypeEncoding = 0x03
	DW_ATE_float           DwarfTypeEncoding = 0x04
	DW_ATE_signed          DwarfTypeEncoding = 0x05
	DW_ATE_signed_char     DwarfTypeEncoding = 0x06
	DW_ATE_unsigned        DwarfTypeEncoding = 0x07
	DW_ATE_unsigned_char   DwarfTypeEncoding = 0x08
	DW_ATE_imaginary_float DwarfTypeEncoding = 0x09
	DW_ATE_packed_decimal  DwarfTypeEncoding = 0x0a
	DW_ATE_numeric_string  DwarfTypeEncoding = 0x0b
	DW_ATE_edited          DwarfTypeEncoding = 0x0c
	DW_ATE_signed_fixed    DwarfTypeEncoding = 0x0d
	DW_ATE_unsigned_fixed  DwarfTypeEncoding = 0x0e
	DW_ATE_decimal_float   DwarfTypeEncoding = 0x0f
	DW_ATE_UTF             DwarfTypeEncoding = 0x10
	DW_ATE_lo_user         DwarfTypeEncoding = 0x80
	DW_ATE_hi_user         DwarfTypeEncoding = 0xff
)

type DebugInfo struct {
	cache map[DebugDescriptor]Value
}

type DebugDescriptor interface {
	// Tag returns the DWARF tag for this descriptor.
	Tag() DwarfTag

	// mdNode creates an LLVM metadata node.
	mdNode(i *DebugInfo) Value
}

///////////////////////////////////////////////////////////////////////////////
// Utility functions.

func constInt1(v bool) Value {
	if v {
		return ConstAllOnes(Int1Type())
	}
	return ConstNull(Int1Type())
}

func (info *DebugInfo) MDNode(d DebugDescriptor) Value {
	// A nil pointer assigned to an interface does not result in a nil
	// interface. Instead, we must check the innards.
	if d == nil || reflect.ValueOf(d).IsNil() {
		return Value{nil}
	}

	if info.cache == nil {
		info.cache = make(map[DebugDescriptor]Value)
	}
	value, exists := info.cache[d]
	if !exists {
		value = d.mdNode(info)
		info.cache[d] = value
	}
	return value
}

func (info *DebugInfo) MDNodes(d []DebugDescriptor) []Value {
	if n := len(d); n > 0 {
		v := make([]Value, n)
		for i := 0; i < n; i++ {
			v[i] = info.MDNode(d[i])
		}
		return v
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// Basic Types

type BasicTypeDescriptor struct {
	Context      DebugDescriptor
	Name         string
	File         *FileDescriptor
	Line         uint32
	Size         uint64 // Size in bits.
	Alignment    uint64 // Alignment in bits.
	Offset       uint64 // Offset in bits
	Flags        uint32
	TypeEncoding DwarfTypeEncoding
}

func (d *BasicTypeDescriptor) Tag() DwarfTag {
	return DW_TAG_base_type
}

func (d *BasicTypeDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), LLVMDebugVersion+uint64(d.Tag()), false),
		info.MDNode(d.File),
		info.MDNode(d.Context),
		MDString(d.Name),
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int64Type(), d.Size, false),
		ConstInt(Int64Type(), d.Alignment, false),
		ConstInt(Int64Type(), d.Offset, false),
		ConstInt(Int32Type(), uint64(d.Flags), false),
		ConstInt(Int32Type(), uint64(d.TypeEncoding), false)})
}

///////////////////////////////////////////////////////////////////////////////
// Composite Types

type CompositeTypeDescriptor struct {
	tag       DwarfTag
	Context   DebugDescriptor
	Name      string
	File      *FileDescriptor
	Line      uint32
	Size      uint64 // Size in bits.
	Alignment uint64 // Alignment in bits.
	Offset    uint64 // Offset in bits
	Flags     uint32
	Members   []DebugDescriptor
}

func (d *CompositeTypeDescriptor) Tag() DwarfTag {
	return d.tag
}

func (d *CompositeTypeDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), LLVMDebugVersion+uint64(d.Tag()), false),
		info.MDNode(d.File),
		info.MDNode(d.Context),
		MDString(d.Name),
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int64Type(), d.Size, false),
		ConstInt(Int64Type(), d.Alignment, false),
		ConstInt(Int64Type(), d.Offset, false),
		ConstInt(Int32Type(), uint64(d.Flags), false),
		info.MDNode(nil), // reference type derived from
		MDNode(info.MDNodes(d.Members)),
		ConstInt(Int32Type(), uint64(0), false), // Runtime language
		ConstInt(Int32Type(), uint64(0), false), // Base type containing the vtable pointer for this type
	})
}

func NewStructCompositeType(
	Members []DebugDescriptor) *CompositeTypeDescriptor {
	d := new(CompositeTypeDescriptor)
	d.tag = DW_TAG_structure_type
	d.Members = Members // XXX Take a copy?
	return d
}

func NewSubroutineCompositeType(
	Result DebugDescriptor,
	Params []DebugDescriptor) *CompositeTypeDescriptor {
	d := new(CompositeTypeDescriptor)
	d.tag = DW_TAG_subroutine_type
	d.Members = make([]DebugDescriptor, len(Params)+1)
	d.Members[0] = Result
	copy(d.Members[1:], Params)
	return d
}

///////////////////////////////////////////////////////////////////////////////
// Compilation Unit

type CompileUnitDescriptor struct {
	Path            FileDescriptor // Path to file being compiled.
	Language        DwarfLang
	Producer        string
	Optimized       bool
	CompilerFlags   string
	Runtime         int32
	EnumTypes       []DebugDescriptor
	RetainedTypes   []DebugDescriptor
	Subprograms     []DebugDescriptor
	GlobalVariables []DebugDescriptor
}

func (d *CompileUnitDescriptor) Tag() DwarfTag {
	return DW_TAG_compile_unit
}

func (d *CompileUnitDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), uint64(d.Tag())+LLVMDebugVersion, false),
		d.Path.mdNode(nil),
		ConstInt(Int32Type(), uint64(d.Language), false),
		MDString(d.Producer),
		constInt1(d.Optimized),
		MDString(d.CompilerFlags),
		ConstInt(Int32Type(), uint64(d.Runtime), false),
		MDNode(info.MDNodes(d.EnumTypes)),
		MDNode(info.MDNodes(d.RetainedTypes)),
		MDNode(info.MDNodes(d.Subprograms)),
		MDNode(info.MDNodes(d.GlobalVariables)),
		MDNode(nil),  // List of imported entities
		MDString(""), // Split debug filename
	})
}

///////////////////////////////////////////////////////////////////////////////
// Derived Types

type DerivedTypeDescriptor struct {
	tag       DwarfTag
	Context   DebugDescriptor
	Name      string
	File      *FileDescriptor
	Line      uint32
	Size      uint64 // Size in bits.
	Alignment uint64 // Alignment in bits.
	Offset    uint64 // Offset in bits
	Flags     uint32
	Base      DebugDescriptor
}

func (d *DerivedTypeDescriptor) Tag() DwarfTag {
	return d.tag
}

func (d *DerivedTypeDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), LLVMDebugVersion+uint64(d.Tag()), false),
		info.MDNode(d.File),
		info.MDNode(d.Context),
		MDString(d.Name),
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int64Type(), d.Size, false),
		ConstInt(Int64Type(), d.Alignment, false),
		ConstInt(Int64Type(), d.Offset, false),
		ConstInt(Int32Type(), uint64(d.Flags), false),
		info.MDNode(d.Base)})
}

func NewPointerDerivedType(Base DebugDescriptor) *DerivedTypeDescriptor {
	d := new(DerivedTypeDescriptor)
	d.tag = DW_TAG_pointer_type
	d.Base = Base
	return d
}

///////////////////////////////////////////////////////////////////////////////
// Subprograms.

type SubprogramDescriptor struct {
	Context     DebugDescriptor
	Name        string
	DisplayName string
	Type        DebugDescriptor
	Line        uint32
	Function    Value
	Path        FileDescriptor
	ScopeLine   uint32
	// Function declaration descriptor
	// Function variables
}

func (d *SubprogramDescriptor) Tag() DwarfTag {
	return DW_TAG_subprogram
}

func (d *SubprogramDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), LLVMDebugVersion+uint64(d.Tag()), false),
		d.Path.mdNode(nil),
		info.MDNode(d.Context),
		MDString(d.Name),
		MDString(d.DisplayName),
		MDString(""), // mips linkage name
		ConstInt(Int32Type(), uint64(d.Line), false),
		info.MDNode(d.Type),
		ConstNull(Int1Type()),                        // not static
		ConstAllOnes(Int1Type()),                     // locally defined (not extern)
		ConstNull(Int32Type()),                       // virtuality
		ConstNull(Int32Type()),                       // index into a virtual function
		info.MDNode(nil),                             // basetype containing the vtable pointer
		ConstInt(Int32Type(), FlagPrototyped, false), // flags
		ConstNull(Int1Type()),                        // not optimised
		d.Function,
		info.MDNode(nil), // Template parameters
		info.MDNode(nil), // function declaration descriptor
		MDNode(nil),      // function variables
		ConstInt(Int32Type(), uint64(d.ScopeLine), false), // Line number where the scope of the subprogram begins
	})
}

///////////////////////////////////////////////////////////////////////////////
// Global Variables.

type GlobalVariableDescriptor struct {
	Context     DebugDescriptor
	Name        string
	DisplayName string
	File        *FileDescriptor
	Line        uint32
	Type        DebugDescriptor
	Local       bool
	External    bool
	Value       Value
}

func (d *GlobalVariableDescriptor) Tag() DwarfTag {
	return DW_TAG_variable
}

func (d *GlobalVariableDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), uint64(d.Tag())+LLVMDebugVersion, false),
		ConstNull(Int32Type()),
		info.MDNode(d.Context),
		MDString(d.Name),
		MDString(d.DisplayName),
		MDNode(nil),
		info.MDNode(d.File),
		ConstInt(Int32Type(), uint64(d.Line), false),
		info.MDNode(d.Type),
		constInt1(d.Local),
		constInt1(!d.External),
		d.Value})
}

///////////////////////////////////////////////////////////////////////////////
// Local Variables.

type LocalVariableDescriptor struct {
	tag      DwarfTag
	Context  DebugDescriptor
	Name     string
	File     DebugDescriptor
	Line     uint32
	Argument uint32
	Type     DebugDescriptor
}

func (d *LocalVariableDescriptor) Tag() DwarfTag {
	return d.tag
}

func (d *LocalVariableDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), uint64(d.Tag())+LLVMDebugVersion, false),
		info.MDNode(d.Context),
		MDString(d.Name),
		info.MDNode(d.File),
		ConstInt(Int32Type(), uint64(d.Line)|(uint64(d.Argument)<<24), false),
		info.MDNode(d.Type),
		ConstNull(Int32Type()), // flags
		ConstNull(Int32Type()), // optional reference to inline location
	})
}

func NewLocalVariableDescriptor(tag DwarfTag) *LocalVariableDescriptor {
	return &LocalVariableDescriptor{tag: tag}
}

///////////////////////////////////////////////////////////////////////////////
// Files.

type FileDescriptor string

func (d *FileDescriptor) Tag() DwarfTag {
	return DW_TAG_file_type
}

func (d *FileDescriptor) mdNode(info *DebugInfo) Value {
	dirname, filename := path.Split(string(*d))
	if l := len(dirname); l > 0 && dirname[l-1] == '/' {
		dirname = dirname[:l-1]
	}
	return MDNode([]Value{MDString(filename), MDString(dirname)})
}

///////////////////////////////////////////////////////////////////////////////
// Line.

type LineDescriptor struct {
	Line    uint32
	Column  uint32
	Context DebugDescriptor
}

func (d *LineDescriptor) Tag() DwarfTag {
	panic("LineDescriptor.Tag should never be called")
}

func (d *LineDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int32Type(), uint64(d.Column), false),
		info.MDNode(d.Context),
		info.MDNode(nil),
	})
}

///////////////////////////////////////////////////////////////////////////////
// Context.

type ContextDescriptor struct{ FileDescriptor }

func (d *ContextDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{ConstInt(Int32Type(), uint64(d.Tag())+LLVMDebugVersion, false), d.FileDescriptor.mdNode(info)})
}

///////////////////////////////////////////////////////////////////////////////
// Block.

type BlockDescriptor struct {
	File    *FileDescriptor
	Context DebugDescriptor
	Line    uint32
	Column  uint32
	Id      uint32
}

func (d *BlockDescriptor) Tag() DwarfTag {
	return DW_TAG_lexical_block
}

func (d *BlockDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), uint64(d.Tag())+LLVMDebugVersion, false),
		info.MDNode(d.File),
		info.MDNode(d.Context),
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int32Type(), uint64(d.Column), false),
		ConstInt(Int32Type(), uint64(d.Id), false),
	})
}

// vim: set ft=go :
