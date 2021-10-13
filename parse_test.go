package sftag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 连接配置
type ConnParams struct {
	Host     string                 `json:"host" yaml:"host" toml:"host"`
	Port     int                    `json:"port,omitempty" yaml:"port,omitempty" toml:"port"`
	Username string                 `yaml/json:"username,omitempty" toml:"username"`
	Password string                 `toml/yaml/json:"password"`
	Database string                 `toml/yaml/json:"database"`
	Options  map[string]interface{} `yaml/json:"options,omitempty" toml:"options"`
}

// go test -run=Parse
func Test_01_Parse(t *testing.T) {
	tags := GetStructTags(ConnParams{})
	tag := NewSfTag()
	tag.Parse(tags["Username"])
	assert.Equal(t, "username,omitempty", tag.Get("yaml"))
	assert.Equal(t, `json:"username,omitempty" toml:"username" yaml:"username,omitempty"`, tag.String())
	tag = NewSfTag()
	tag.Parse(tags["Password"])
	assert.Equal(t, "password", tag.Get("yaml"))
	assert.Equal(t, `json:"password" toml:"password" yaml:"password"`, tag.String())
}
