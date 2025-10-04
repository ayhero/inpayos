package protocol

import (
	"fmt"
	"inpayos/internal/utils"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/shopspring/decimal"
)

var (
	NUMBER_SUFFIX = []string{"_amount", "_count"}
	NUMBER_FIELDS = []string{"amount"}
)

// 表达式
//
// ["amount","GE",50]
type Express []string

func (t Express) ToString() string {
	return utils.ToJsonString(t)
}

func (t Express) IsNumber() bool {
	field := t[0]
	for _, s := range NUMBER_SUFFIX {
		if strings.Contains(field, s) {
			return true
		}
	}
	return utils.StringInArray(field, NUMBER_FIELDS)
}

func (t Express) ToGovaluateExpress() string {
	field_name := t[0]
	symbol := Symbol(t[1])
	match_val := t[2]
	express_template := GOVALUATE_STRING_EXP_TMP
	if t.IsNumber() {
		express_template = GOVALUATE_NUM_EXP_TMP
	}
	govaluate_symbol := symbol.ToGovaluateSymbol()
	express_ctn := fmt.Sprintf(express_template, field_name, govaluate_symbol, match_val)
	if symbol == NOT_IN_MATCH {
		express_ctn = fmt.Sprintf(GOVALUATE_NOT_EXP_TMP, express_ctn)
	}
	return express_ctn
}
func (t Express) ToGovaluate() *govaluate.EvaluableExpression {
	express_ctn := t.ToGovaluateExpress()
	exp, err := govaluate.NewEvaluableExpression(express_ctn)
	if err != nil {
		return nil
	}
	return exp
}

func NewEqExpress(target string, match string) (exp Express) {
	exp, _ = NewExpress(target, EQ_MATCH, match)
	return
}
func NewExpress(target string, symbol Symbol, match string) (exp Express, ok bool) {
	ok = false
	if target == "" || symbol == "" || match == "" || symbol.ToString() == "" {
		return
	}
	ok = true
	exp = Express([]string{target, symbol.ToString(), match})
	return
}
func NewExpressFromArrays(exps []string) (exp Express, ok bool) {
	ok = false
	if len(exps) != 3 {
		return
	}
	exp, ok = NewExpress(exps[0], Symbol(exps[1]), exps[2])
	return
}
func NewExpressFromString(content string) (exp Express, ok bool) {
	ok = false
	if content == "" {
		return
	}
	exps := strings.Split(content, "|")
	if len(exps) == 0 {
		return
	}
	return NewExpressFromArrays(exps)
}
func (t Express) Match(target string) (rs bool) {
	rs = true
	if t.IsNumber() {
		rs = t.MatchNumberVal(target)
	} else {
		rs = t.MatchStringVal(target)
	}
	return
}

func (t Express) MatchNumberVal(target string) (result bool) {
	result = true
	symbol := Symbol(t[1])
	match := t[2]
	target_val, err := decimal.NewFromString(target)
	if err != nil {
		target_val = decimal.Zero
	}
	match_val, err := decimal.NewFromString(match)
	if err != nil {
		match_val = decimal.Zero
	}
	switch symbol {
	case GT_MATCH:
		result = target_val.GreaterThan(match_val)
	case GE_MATCH:
		result = target_val.GreaterThanOrEqual(match_val)
	case EQ_MATCH:
		result = target_val.Equal(match_val)
	case LE_MATCH:
		result = target_val.LessThanOrEqual(match_val)
	case LT_MATCH:
		result = target_val.LessThan(match_val)
	}
	return result
}

func (t Express) MatchStringVal(target string) (result bool) {
	result = true
	symbol := Symbol(t[1])
	match := t[2]
	switch symbol {
	case EQ_MATCH:
		result = target == match
	case IN_MATCH:
		tg_list := strings.Split(match, ",")
		for _, _v := range tg_list {
			if _v == target {
				return true
			}
		}
	case NOT_IN_MATCH:
		tg_list := strings.Split(match, ",")
		for _, _v := range tg_list {
			if _v == target {
				return false
			}
		}
	}
	return result
}
