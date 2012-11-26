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
	LLVMDebugVersion = (11 << 16)
)

type DwarfTag uint32

const (
	DW_TAG_compile_unit    DwarfTag = 0x11
	DW_TAG_variable        DwarfTag = 0x34
	DW_TAG_base_type       DwarfTag = 0x24
	DW_TAG_pointer_type    DwarfTag = 0x0F
	DW_TAG_structure_type  DwarfTag = 0x13
	DW_TAG_subroutine_type DwarfTag = 0x15
	DW_TAG_file_type       DwarfTag = 0x29
	DW_TAG_subprogram      DwarfTag = 0x2E
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
		info.MDNode(d.Context),
		MDString(d.Name),
		info.MDNode(d.File),
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int32Type(), d.Size, false),
		ConstInt(Int32Type(), d.Alignment, false),
		ConstInt(Int32Type(), d.Offset, false),
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
		info.MDNode(d.Context),
		MDString(d.Name),
		info.MDNode(d.File),
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int32Type(), d.Size, false),
		ConstInt(Int32Type(), d.Alignment, false),
		ConstInt(Int32Type(), d.Offset, false),
		ConstInt(Int32Type(), uint64(d.Flags), false),
		MDNode(nil),
		MDNode(info.MDNodes(d.Members)),
		ConstInt(Int32Type(), uint64(0), false)})
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
	Language        DwarfLang
	Path            string // Path to file being compiled.
	Producer        string
	MainCompileUnit bool
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
	dirname, filename := path.Split(d.Path)
	return MDNode([]Value{
		ConstInt(Int32Type(), uint64(d.Tag())+LLVMDebugVersion, false),
		ConstNull(Int32Type()),
		ConstInt(Int32Type(), uint64(d.Language), false),
		MDString(filename),
		MDString(dirname),
		MDString(d.Producer),
		constInt1(d.MainCompileUnit),
		constInt1(d.Optimized),
		MDString(d.CompilerFlags),
		ConstInt(Int32Type(), uint64(d.Runtime), false),
		MDNode(info.MDNodes(d.EnumTypes)),
		MDNode(info.MDNodes(d.RetainedTypes)),
		MDNode(info.MDNodes(d.Subprograms)),
		MDNode(info.MDNodes(d.GlobalVariables))})
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
		info.MDNode(d.Context),
		MDString(d.Name),
		info.MDNode(d.File),
		ConstInt(Int32Type(), uint64(d.Line), false),
		ConstInt(Int32Type(), d.Size, false),
		ConstInt(Int32Type(), d.Alignment, false),
		ConstInt(Int32Type(), d.Offset, false),
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
	File        *FileDescriptor
	Type        DebugDescriptor
	Line        uint32
	Function    Value
	// Function declaration descriptor
	// Function variables
}

func (d *SubprogramDescriptor) Tag() DwarfTag {
	return DW_TAG_subprogram
}

func (d *SubprogramDescriptor) mdNode(info *DebugInfo) Value {
	return MDNode([]Value{
		ConstInt(Int32Type(), LLVMDebugVersion+uint64(d.Tag()), false),
		ConstNull(Int32Type()),
		info.MDNode(d.Context),
		MDString(d.Name),
		MDString(d.DisplayName),
		MDNode(nil),
		info.MDNode(d.File),
		ConstInt(Int32Type(), uint64(d.Line), false),
		info.MDNode(d.Type),
		ConstNull(Int1Type()),    // not static
		ConstAllOnes(Int1Type()), // locally defined (not extern)
		ConstNull(Int32Type()),
		ConstNull(Int32Type()),
		MDNode(nil),
		ConstNull(Int32Type()), // flags
		ConstNull(Int1Type()),  // not optimised
		d.Function,
		MDNode(nil),
		MDNode(nil),  // function declaration descriptor
		MDNode(nil)}) // function variables
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
// Files.

type FileDescriptor string

func (d *FileDescriptor) Tag() DwarfTag {
	return DW_TAG_file_type
}

func (d *FileDescriptor) mdNode(info *DebugInfo) Value {
	dirname, filename := path.Split(string(*d))
	return MDNode([]Value{MDString(filename), MDString(dirname), MDNode(nil)})
}

// vim: set ft=go :
