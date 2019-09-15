package options

import (
	"github.com/spf13/pflag"
)

// ServerRunOptions
type ServerRunOptions struct {
	Port         int
	Kubeconfig   string
	Config       string
	TemplatePath string
}

// NewServerRunOptions
func NewServerRunOptions() *ServerRunOptions {
	opt := &ServerRunOptions{}
	return opt
}

// ServerRunOptions 添加命令行参数
func (s *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {
	fs.IntVar(&s.Port, "port", 8080, "api port.")
	fs.StringVar(&s.Config, "config", "/tmp/config.toml", "server config.")
	fs.StringVar(&s.Kubeconfig, "kubeconfig", "/etc/kubernetes/config.yaml", "kubernetes client config.")
	fs.StringVar(&s.TemplatePath, "template_path", "/etc/template/pod.template.yam", "pod template.")
}
