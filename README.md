# go-sftag
Yet another StructTag, with some features liked cache and alias.

# Usage
```go
package main
import (
	"fmt"
    
	"github.com/azhai/go-sftag"
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

func main() {
	tags := sftag.GetStructTags(ConnParams{})
	tag := sftag.NewSfTag()
	tag.Parse(tags["Options"])
	fmt.Println(tag.Get("yaml"))
	fmt.Println(tag.String())
}
```

