package render

import (
	"fmt"
	"html/template"
	"io"
)

func RenderTempate(page string, out io.Writer, data any) error {
	path := fmt.Sprintf("./tmpl/%s.tmpl.html", page)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	err = tmpl.Execute(out, data)
	return err
}
