package apis

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/pinlan/app-admin/cmd/app-admin/app/options"
	trainingv1 "github.com/pinlan/app-admin/internal/app/app-admin/apis/v1/training"
	"github.com/pinlan/app-admin/internal/app/app-admin/template"
	bll "github.com/pinlan/app-admin/internal/app/app-admin/training"
)

// InitKubeClientset
func InitKubeClientset(kubeconfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

// Server
func Server(opt *options.ServerRunOptions) error {
	r := gin.Default()

	// Set up cross-domain requests
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))

	// Init kubernetes client
	clientset, err := InitKubeClientset(opt.Kubeconfig)
	if err != nil {
		panic("InitKubeClientset Failed...")
	}

	// Init training task pod template
	podTemplate := template.NewPodTemplate()
	err = podTemplate.LoadPodTemplate(opt.TemplatePath)
	if err != nil {
		panic("LoadPodTemplate Failed...")
	}

	// Init training server
	server := bll.NewServer(clientset, podTemplate)
	trainingServer := trainingv1.NewTrainingServer(server)

	// Register api route
	v1 := r.Group("/api/v1")
	{
		v1.POST("/trainings", trainingServer.Create)
		v1.GET("/trainings", trainingServer.List)
		v1.GET("/trainings/:id", trainingServer.Get)
		v1.DELETE("/trainings/:id", trainingServer.Delete)
		v1.GET("/trainings/:id/logs", trainingServer.GetLogs)
	}

	// Listen api port
	r.Run(fmt.Sprintf(":%d", opt.Port))
	return nil
}
