package training

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	//"k8s.io/apimachinery/pkg/api/resource"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/pinlan/app-admin/internal/app/app-admin/schema"
	"github.com/pinlan/app-admin/pkg/utils"
)

// ConvertToPod
func ConvertFromPodTemplate(user_id string, training schema.TrainingSchema, pod_template *v1.Pod) {
	labels := generateLabels(user_id, training.Name)
	id := generateName(user_id, training.Name)

	pod_template.Labels = labels
	pod_template.Name = id
	command := make([]string, 0)
	command = append(command, "/bin/sh")
	command = append(command, "-c")
	args := convertToArgs(training.Args)
	command = append(command, args)
	pod_template.Spec.Containers[0].Command = command
}

// ConvertToSchema
func ConvertToSchema(pod v1.Pod) schema.Training {
	training := schema.Training{}
	training.ID = pod.Name
	training.Name = pod.Labels["name"]
	training.Status = string(pod.Status.Phase)

	created_at := generateTimestamp(pod.CreationTimestamp.String())
	training.CreatedAt = created_at
	return training
}

// generateLabels
func generateLabels(user_id, name string) map[string]string {
	labels := make(map[string]string)
	labels["user_id"] = user_id
	labels["app"] = "training"
	labels["name"] = name
	return labels
}

// generateName
func generateName(user_id, name string) string {
	id := utils.GenerateID(user_id + name)
	return "training" + "-" + id
}

// generaTimestamp
func generateTimestamp(timeStr string) int64 {
	loc, _ := time.LoadLocation("Local")
	the_time, err := time.ParseInLocation("2006-01-02 15:04:05 +0100 BST", timeStr, loc)
	if err == nil {
		unix_time := the_time.Unix() //1504082441
		fmt.Println(unix_time)
		return unix_time
	}

	return 0
}

// convertToArgs
func convertToArgs(args map[string]string) string {
	args_str := ""
	for key, value := range args {
		args_str = args_str + "--" + key + "=" + value + " "
	}

	return "python3 main.py " + args_str
}
