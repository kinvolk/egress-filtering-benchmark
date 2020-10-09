package util

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

func DecodeAndApply(config *rest.Config, data string, action string) error {
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	dd, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("creating dynamic config: %w", err)
	}

	r := strings.NewReader(data)
	yamlReader := yamlutil.NewYAMLReader(bufio.NewReader(r))
	b, err := yamlReader.Read()
	if err != nil {
		return fmt.Errorf("reading yaml: %w", err)
	}
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(b), len(b))
	var rawObj runtime.RawExtension
	if err := decoder.Decode(&rawObj); err != nil {
		return fmt.Errorf("decoding to raw object: %w", err)
	}

	obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return fmt.Errorf("serializing to unstructured: %w", err)
	}

	unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}

	gr, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		return fmt.Errorf("getting APIGroupResources: %w", err)
	}

	mapper := restmapper.NewDiscoveryRESTMapper(gr)
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("rest mapping: %w", err)
	}

	var dri dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		if unstructuredObj.GetNamespace() == "" {
			unstructuredObj.SetNamespace("default")
		}
		dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
	} else {
		dri = dd.Resource(mapping.Resource)
	}

	if action == "CREATE" {
		if _, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("creating resource %s: %w", unstructuredObj.GetName(), err)
		}
	}

	if action == "DELETE" {
		if err := dri.Delete(context.Background(), unstructuredObj.GetName(), metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("deleting resource %s: %w", unstructuredObj.GetName(), err)
		}
	}

	return nil
}
