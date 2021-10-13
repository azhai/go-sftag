package sftag

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

//相同值的tag分隔符
const TAG_KEY_SEP = "/"

// StructFieldTag Yet another StructTag
type StructFieldTag struct {
	alias   map[string]string
	data    map[string]string
	changed bool
	lock    sync.RWMutex
	reflect.StructTag
}

func NewSfTag() *StructFieldTag {
	it := new(StructFieldTag)
	it.alias = make(map[string]string)
	it.data = make(map[string]string)
	return it
}

// Parse 解析全部tags
func (it *StructFieldTag) Parse(tag reflect.StructTag) {
	it.StructTag, it.changed = tag, true
	word, key, chunk := Word(tag), "", ""
	for word != "" {
		word = word.SkipChar(" ")
		if word == "" {
			break
		}
		key, word = word.MatchChar(":\"")
		if key == "" {
			break
		}
		chunk, word = word.MatchChar("\"")
		value, err := strconv.Unquote("\"" + chunk + "\"")
		if err != nil {
			break
		}
		it.Append(key, value)
	}
}

// Build 转为字符串格式
func (it StructFieldTag) Build(data []byte, name, value string) []byte {
	data = append(data, []byte(name)...)
	data = append(data, byte(':'), byte('"'))
	data = append(data, []byte(value)...)
	data = append(data, byte('"'), byte(' '))
	return data
}

// String 转为字符串格式
func (it StructFieldTag) String() string {
	if !it.changed {
		return string(it.StructTag)
	}
	var data []byte
	if len(it.alias) > len(it.data) {
		for _, key := range SortedKeys(it.alias) {
			data = it.Build(data, key, it.Get(key))
		}
	} else {
		for _, name := range SortedKeys(it.data) {
			data = it.Build(data, name, it.data[name])
		}
	}
	result := strings.TrimSpace(string(data))
	it.StructTag = reflect.StructTag(result)
	return result
}

// Get returns the value associated with key in the tag string
func (it StructFieldTag) Get(key string) string {
	v, _ := it.Lookup(key)
	return v
}

// Lookup Returns a tag from the tag data
func (it StructFieldTag) Lookup(key string) (string, bool) {
	it.lock.RLock()
	defer it.lock.RUnlock()
	if name, ok := it.alias[key]; ok {
		key = name
	}
	value, ok := it.data[key]
	return value, ok
}

// Append Sets a tag in the tag data map
func (it *StructFieldTag) Append(key, value string) {
	it.lock.Lock()
	defer it.lock.Unlock()
	names := strings.Split(key, TAG_KEY_SEP)
	name := names[len(names)-1]
	for _, key = range names {
		it.alias[key] = name
	}
	it.data[name] = value
	it.changed = true
}

// Delete Deletes a tag
func (it *StructFieldTag) Delete(key string) {
	it.lock.Lock()
	defer it.lock.Unlock()
	delete(it.alias, key)
	delete(it.data, key)
	it.changed = true
}

// GetStructTags Read all tags in a object
func GetStructTags(v interface{}) (tags map[string]reflect.StructTag) {
	tags = make(map[string]reflect.StructTag)
	vt := GetIndirectType(v)
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		tags[field.Name] = field.Tag
	}
	return
}
