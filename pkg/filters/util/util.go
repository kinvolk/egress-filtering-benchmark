package util

import (
	"bytes"
	"text/template"
)

const (
	// This number is restricted by the number of entries that can be present in an ipset
	// While that number could be `1048576`, the actual number is far less because of the
	// maximum size that can be POST'ed to the Kubernetes API server before getting error
	// `etcdserver: Request entity too large`.
	// This number can certainly be improved upon which would result in lesser number of
	// GlobalNetworkSet batch thereby resulting in quicker turnaround time of applying the
	// manifests. However the expected improvement is small.
	RulesPerManifest = 50000
	// Timout in milliseconds.
	Timeout = 600000
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
