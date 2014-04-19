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
