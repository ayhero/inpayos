package utils

import (
	"strings"

	"github.com/spf13/cast"
)

// 将source字符串拆分后合并为sql 的in查询值
//
// 如:
//
// source: ["a","b","c"],ori: ",", tgt: "','"
//
// return: "a','b','c"
func ArrayToSqlString(source string) string {
	return ArrayToString(cast.ToString(source), ",", "','")
}

// 根据ori的分隔符，将source字符串拆解，后按tgt的分隔符，合并为新的字符串
//
// 如:
//
// source: ["a","b","c"],ori: ",", tgt: "','"
//
// return: "a','b','c"
func ArrayToString(source, ori, tgt string) string {
	_source := cast.ToString(source)
	if _source == "" {
		return ""
	}
	_ls := strings.Split(source, ori)
	return strings.Join(_ls, tgt)
}

func ArrayToInt(_ls []string) []int64 {
	n_ls := []int64{}
	for _, v := range _ls {
		n_ls = append(n_ls, cast.ToInt64(v))
	}
	return n_ls
}

// 对source按,分割后，去重复，返回不重复的数组
//
// "a,a,b" => ["a","b"]
func SplitAndUniqueDefault(source string) (rs []string) {
	rs = []string{}
	_ls := strings.Split(source, ",")
	if len(_ls) == 0 {
		return
	}
	rs = UniqueArray(_ls)
	return
}

// 对source按,分割后:去重复，返回不重复的数组
//
// "a:a:b" => ["a","b"]
func SplitAndUniqueByColon(source string) (rs []string) {
	rs = []string{}
	_ls := strings.Split(source, ":")
	if len(_ls) == 0 {
		return
	}
	rs = UniqueArray(_ls)
	return
}

// 对source按|分割后，去重复，返回不重复的数组
//
// "a|a|b" => ["a","b"]
func SplitAndUniqueByVL(source string) (rs []string) {
	rs = []string{}
	_ls := strings.Split(source, "|")
	if len(_ls) == 0 {
		return
	}
	rs = UniqueArray(_ls)
	return
}

// 根据sep对source分割后，去重复，返回不重复的数组
//
// "a,a,b" => ["a","b"]
func SplitAndUnique(source, sep string) (rs []string) {
	rs = []string{}
	_ls := strings.Split(source, sep)
	if len(_ls) == 0 {
		return
	}
	rs = UniqueArray(_ls)
	return
}

// 对_ls去重复
//
// ["a","a","b"] =>["a","b"]
func UniqueArray(_ls []string) (rs []string) {
	if len(_ls) == 0 {
		return
	}
	//过滤重复值
	rs = []string{}
	_t_list := map[string]bool{}
	for _, l := range _ls {
		_l := cast.ToString(l)
		if _l == "" {
			continue
		}
		if _, ok := _t_list[_l]; !ok {
			_t_list[_l] = true
			rs = append(rs, _l)
		}
	}
	return
}

// 对_ls去重复，根据sep合并并返回结果
//
// ["a","a","b"] => "a,b"
func UniqueAndCombind(_ls []string, sep string) (source string) {
	if len(_ls) == 0 {
		source = ""
		return
	}
	//过滤重复值
	_s_ls := []string{}
	_s_lsm := map[string]bool{}
	for _, l := range _ls {
		_l := cast.ToString(l)
		if _l != "" {
			_s_lsm[_l] = true
		}
	}
	if len(_s_lsm) > 0 {
		for k, _ := range _s_lsm {
			_s_ls = append(_s_ls, k)
		}
	}
	_ls = _s_ls
	source = strings.Join(_ls, sep)
	return
}

func StringInArray(target string, list []string) bool {
	if len(list) == 0 {
		return false
	}
	for _, _l := range list {
		if target == _l {
			return true
		}
	}
	return false
}

func StringNotInArray(target string, list []string) bool {
	return !StringInArray(target, list)
}

func InArrays(target int, list []int) bool {
	if len(list) == 0 {
		return false
	}
	for _, _l := range list {
		if target == _l {
			return true
		}
	}
	return false
}
