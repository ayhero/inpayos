package protocol

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type Symbol string

const (
	LE_MATCH     Symbol = "le"
	LT_MATCH     Symbol = "lt"
	EQ_MATCH     Symbol = "eq"
	NOT_EQ_MATCH Symbol = "neq"
	GT_MATCH     Symbol = "gt"
	GE_MATCH     Symbol = "ge"
	IN_MATCH     Symbol = "in"
	NOT_IN_MATCH Symbol = "nin"

	AND_LOGIC Symbol = "and"
	OR_LOGIC  Symbol = "or"

	NOT_EXP_TMP                    = " !(%v) "
	GOVALUATE_NUM_EXP_TMP          = " %v %v %v "
	GOVALUATE_NUM_IN_RANGE_EXP_TMP = " %v in (%v) "
	GOVALUATE_STRING_EXP_TMP       = " %v %v '%v' "
	GOVALUATE_STRING_IN_EXP_TMP    = " %v in ('%v') "
	GOVALUATE_NOT_EXP_TMP          = " !(%v) "
)

func (s Symbol) ToString() string {
	switch s {
	case GT_MATCH:
		fallthrough
	case GE_MATCH:
		fallthrough
	case EQ_MATCH:
		fallthrough
	case NOT_EQ_MATCH:
		fallthrough
	case LT_MATCH:
		fallthrough
	case LE_MATCH:
		fallthrough
	case IN_MATCH:
		fallthrough
	case NOT_IN_MATCH:
		fallthrough
	case AND_LOGIC:
		fallthrough
	case OR_LOGIC:
		return string(s)
	default:
		return ""
	}
}
func (s Symbol) ToGovaluateSymbol() string {
	switch s {
	case GT_MATCH:
		return govaluate.GT.String()
	case GE_MATCH:
		return govaluate.GTE.String()
	case EQ_MATCH:
		return fmt.Sprintf("%v%v", govaluate.EQ.String(), govaluate.EQ.String())
	case NOT_EQ_MATCH:
		return govaluate.NEQ.String()
	case LT_MATCH:
		return govaluate.LT.String()
	case LE_MATCH:
		return govaluate.LTE.String()
	case IN_MATCH:
		return govaluate.IN.String()
	case NOT_IN_MATCH:
		return govaluate.IN.String()
	case AND_LOGIC:
		return govaluate.AND.String()
	case OR_LOGIC:
		return govaluate.OR.String()
	default:
		return ""
	}
}
func (s Symbol) ToSqlSymbol() string {
	switch s {
	case GT_MATCH:
		return ">"
	case GE_MATCH:
		return ">="
	case EQ_MATCH:
		return "="
	case NOT_EQ_MATCH:
		return "!="
	case LT_MATCH:
		return "<"
	case LE_MATCH:
		return "<="
	case IN_MATCH:
		return "IN"
	case NOT_IN_MATCH:
		return "NOT IN"
	case AND_LOGIC:
		return "AND"
	case OR_LOGIC:
		return "OR"
	default:
		return ""
	}
}
