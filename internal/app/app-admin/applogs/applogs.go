package applogs

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Applogs
type Applogs struct {
	Clientset kubernetes.Interface
}

// PodEvent
type PodEvent struct {
	Reason  string
	Message string
}

// PodInfo
type PodInfo struct {
	Namespace string            `json:"namespace"`
	PodName   string            `json:"pod_name"`
	Labels    map[string]string `json:"labels"`
	Status    string            `json:"status"`
}

// NewApplogs
func NewApplogs(clientset kubernetes.Interface) *Applogs {
	return &Applogs{Clientset: clientset}
}

// GetPodList
func (cli *Applogs) GetPodList(namespace, server_type, server_name string) ([]PodInfo, error) {
	pod_lists := make([]PodInfo, 0)
	labelSelector := "app=" + server_type + ",name=" + server_name
	opt := metav1.ListOptions{LabelSelector: labelSelector}

	pods, err := cli.Clientset.CoreV1().Pods(namespace).List(opt)
	if err != nil {
		return pod_lists, err
	}

	if len(pods.Items) == 0 {
		return pod_lists, nil
	}

	for _, pod := range pods.Items {
		pod_info := new(PodInfo)

		pod_info.Namespace = namespace
		pod_info.PodName = pod.ObjectMeta.Name
		pod_info.Labels = pod.ObjectMeta.Labels

		status := GetStatusForPod(&pod)

		pod_info.Status = string(status)
		pod_lists = append(pod_lists, *pod_info)
	}
	return pod_lists, nil
}

// GetLogs
func (cli *Applogs) GetLogs(namespace, pod_name string) (string, error) {
	pod, err := cli.Clientset.CoreV1().Pods(namespace).Get(pod_name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	status := GetStatusForPod(pod)

	var log string
	if status == "Pending" {
		log, err = cli.GetDescribe(namespace, pod_name)
		if err != nil {
			return "", err
		}
		return log, nil
	}

	response := cli.Clientset.CoreV1().Pods(namespace).GetLogs(pod_name, &v1.PodLogOptions{})
	log_info, err := response.Do().Raw()
	if err != nil {
		return "", err
	}

	return string(log_info), nil
}

// GetDescribe
func (cli *Applogs) GetDescribe(namespace, pod_name string) (string, error) {
	podevent, err := cli.Clientset.CoreV1().Events(namespace).Get(pod_name, metav1.GetOptions{})
	return podevent.Message, err
}

// GetStatusForPod
func GetStatusForPod(pod *v1.Pod) string {
	status := string(pod.Status.Phase)
	for _, containerstatus := range pod.Status.ContainerStatuses {
		if containerstatus.State.Waiting != nil {
			status = containerstatus.State.Waiting.Reason
			break
		}
		if containerstatus.State.Terminated != nil {
			status = containerstatus.State.Terminated.Reason
		}
		if containerstatus.Ready == false {
			status = "NotReady"
			break
		}
	}
	return status
}
