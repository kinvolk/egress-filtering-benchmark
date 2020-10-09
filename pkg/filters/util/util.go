package util

import (
	"bytes"
	"text/template"
)

func RenderTemplate(tmpl string, obj interface{}) (string, error) {
	t, err := template.New("render").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, obj); err != nil {
		return "", err
	}

	return buf.String(), nil
}
