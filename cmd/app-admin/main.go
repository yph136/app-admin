package main

import (
	"github.com/golang/glog"
	"github.com/spf13/pflag"

	"github.com/pinlan/app-admin/cmd/app-admin/app"
	"github.com/pinlan/app-admin/cmd/app-admin/app/options"
)

// main
func main() {

	s := options.NewServerRunOptions()
	s.AddFlags(pflag.CommandLine)

	pflag.Parse()
	defer glog.Flush()

	if err := app.Run(s); err != nil {
		glog.Error("app exit...")
	}
}
