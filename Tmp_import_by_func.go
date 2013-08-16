package main

import (
	"bytes"
	"os"
	"text/template"
)

// 演示 template FuncMap 和闭包结合使用
func main() {
	tpl := template.New("")
	// 使用闭包让 funcs 可以直接读取到 tpl 对象
	tpl.Funcs(template.FuncMap{
		"import": func(filename string) string {
			tpl.New(filename).Parse("import ok")
			ts := tpl.Templates()
			bs := bytes.NewBufferString("")
			ts[len(ts)-1].Execute(bs, "")
			return bs.String()
		},
	})
	// 这里直接用字符串模拟 ParseFiles, 现实中用那个您随意
	tpl.New("import.tmpl").Parse(`{{import "foo.tmpl"}}`)
	tpl.ExecuteTemplate(os.Stdout, "import.tmpl", "")
}
