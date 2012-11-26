package llvm

import "fmt"

func (t TypeKind) String() string {
	switch t {
	case VoidTypeKind:
		return "VoidTypeKind"
	case FloatTypeKind:
		return "FloatTypeKind"
	case DoubleTypeKind:
		return "DoubleTypeKind"
	case X86_FP80TypeKind:
		return "X86_FP80TypeKind"
	case FP128TypeKind:
		return "FP128TypeKind"
	case PPC_FP128TypeKind:
		return "PPC_FP128TypeKind"
	case LabelTypeKind:
		return "LabelTypeKind"
	case IntegerTypeKind:
		return "IntegerTypeKind"
	case FunctionTypeKind:
		return "FunctionTypeKind"
	case StructTypeKind:
		return "StructTypeKind"
	case ArrayTypeKind:
		return "ArrayTypeKind"
	case PointerTypeKind:
		return "PointerTypeKind"
	case VectorTypeKind:
		return "VectorTypeKind"
	case MetadataTypeKind:
		return "MetadataTypeKind"
	}
	panic("unreachable")
}

func (t Type) String() string {
	k := t.TypeKind()
	s := k.String()
	s = s[:len(s)-len("Kind")]

	switch k {
	case ArrayTypeKind:
		s += fmt.Sprintf("(%v[%v])", t.ElementType(), t.ArrayLength())
	case PointerTypeKind:
		s += fmt.Sprintf("(%v)", t.ElementType())
	case StructTypeKind:
		etypes := t.StructElementTypes()
		s += "("
		if n := len(etypes); n > 0 {
			s += fmt.Sprint(etypes[0])
			for i := 1; i < n; i++ {
				s += fmt.Sprintf(", %v", etypes[i])
			}
		}
		s += ")"
	}

	return s
}

// vim: set ft=go :
