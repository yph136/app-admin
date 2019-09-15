package app

import (
	"github.com/golang/glog"

	"github.com/pinlan/app-admin/cmd/app-admin/app/options"
	tserver "github.com/pinlan/app-admin/internal/app/app-admin/apis"
)

// Run
func Run(opt *options.ServerRunOptions) error {
	err := tserver.Server(opt)
	if err != nil {
		glog.Errorf("Server Run failed: %v", err.Error())
		return err
	}
	return nil
}
