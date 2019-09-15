package template

import (
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

// Podtemplate
type PodTemplate struct {
	Pod *corev1.Pod
}

// NewPodTemplate
func NewPodTemplate() *PodTemplate {
	return &PodTemplate{}
}

// LoadTemplate
func (template *PodTemplate) LoadPodTemplate(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(yamlFile), nil, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	switch o := obj.(type) {
	case *corev1.Pod:
		template.Pod = o
		return nil
	default:
		return fmt.Errorf("Failed...")
	}
}

// GetPodTemplate
func (pod *PodTemplate) GetPodTemplate() corev1.Pod {
	return *(pod.Pod)
}
