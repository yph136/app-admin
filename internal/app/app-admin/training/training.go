package training

import (
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/pinlan/app-admin/internal/app/app-admin/applogs"
	"github.com/pinlan/app-admin/internal/app/app-admin/schema"
	"github.com/pinlan/app-admin/internal/app/app-admin/template"
)

const NAMESPACE = "training"

// Server
type Server struct {
	clientset     kubernetes.Interface
	applogsServer *applogs.Applogs
	podTemplate   *template.PodTemplate
}

// NewServer
func NewServer(clientset kubernetes.Interface, podTemplate *template.PodTemplate) *Server {
	applogsServer := applogs.NewApplogs(clientset)
	return &Server{clientset: clientset, applogsServer: applogsServer, podTemplate: podTemplate}
}

// Get
func (server *Server) Get(id string) (schema.Training, error) {
	options := metav1.GetOptions{}
	pod, err := server.clientset.CoreV1().Pods(NAMESPACE).Get(id, options)
	if err != nil {
		glog.Errorf("Get Pod Info Failed: %v", err.Error())
		return schema.Training{}, err
	}

	training := ConvertToSchema(*pod)
	return training, nil
}

// Create
func (server *Server) Create(user_id string, training_shcmea schema.TrainingSchema) (schema.Training, error) {
	podTemplate := server.podTemplate.GetPodTemplate()
	ConvertFromPodTemplate(user_id, training_shcmea, &podTemplate)
	_, err := server.clientset.CoreV1().Pods(NAMESPACE).Create(&podTemplate)
	if err != nil {
		glog.Errorf("create pod failed: %v", err.Error())
		return schema.Training{}, err
	}

	id := podTemplate.Name
	status := "Pending"
	return schema.Training{
		ID:     id,
		Name:   training_shcmea.Name,
		Status: status,
	}, nil
}

// List
func (server *Server) List(user_id string) (schema.TrainingList, error) {
	labelSelector := "user_id=" + user_id
	options := metav1.ListOptions{LabelSelector: labelSelector}
	pods, err := server.clientset.CoreV1().Pods(NAMESPACE).List(options)
	if err != nil {
		glog.Errorf("list pods failed: %v", err.Error())
		return schema.TrainingList{}, err
	}

	training_list := make([]schema.Training, 0)
	for _, pod := range pods.Items {
		training := ConvertToSchema(pod)
		training_list = append(training_list, training)
	}
	return schema.TrainingList{
		Items: training_list,
		Total: int32(len(training_list)),
	}, nil
}

// Delete
func (server *Server) Delete(id string) error {
	options := &metav1.DeleteOptions{}
	err := server.clientset.CoreV1().Pods(NAMESPACE).Delete(id, options)
	return err
}

// GetLogs
func (server *Server) GetLogs(id string) (string, error) {
	log_details, err := server.applogsServer.GetLogs(NAMESPACE, id)
	return log_details, err
}
