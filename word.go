package sftag

import (
	"reflect"
	"sort"
	"strings"
)

type Word string

// SkipChar 忽略开头字符
func (w Word) SkipChar(chars string) Word {
	i := 0
	for i < len(w) && strings.ContainsRune(chars, rune(w[i])) {
		i++
	}
	return w[i:]
}

// MatchChar 寻找字符第一次出现的位置，类似IndexAny
func (w Word) MatchChar(chars string) (string, Word) {
	i := 0
	for i < len(w) && w[i] >= 0x20 && w[i] < 0x7f && w[i] != chars[0] {
		if w[i] == '\\' { //不要匹配同名转义字符，例如目标是"就要跳过\"
			i++
		}
		i++
	}
	if i == 0 || i >= len(w) || !strings.HasPrefix(string(w[i:]), chars) {
		return "", ""
	}
	return string(w[:i]), w[i+len(chars):]
}

// GetIndirectType 获取对象（指针）的实际类型
func GetIndirectType(v interface{}) (rt reflect.Type) {
	var ok bool
	if rt, ok = v.(reflect.Type); !ok {
		rt = reflect.TypeOf(v)
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return
}

// SortedKeys 对map的key按字母排序
func SortedKeys(data interface{}) (keys []string) {
	rt := GetIndirectType(data)
	if rt.Kind() != reflect.Map || rt.Key().Kind() != reflect.String {
		return // data 必须是 map[string]xxxx 类型
	}
	for _, rv := range reflect.ValueOf(data).MapKeys() {
		keys = append(keys, rv.String())
	}
	sort.Strings(keys)
	return
}
