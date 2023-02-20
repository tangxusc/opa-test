package discovery

import (
	"github.com/spf13/cobra"
	"plugin_center/pkg/config"
)

var servers []string
var Namespace string
var serviceName string
var GroupName string

func InitFlags() {
	config.RegisterFlags(func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringSliceVar(&servers, "servers", []string{"127.0.0.1"}, "discovery servers address,split by ,")
		cmd.PersistentFlags().StringVar(&Namespace, "Namespace", "", "discovery Namespace")
		cmd.PersistentFlags().StringVar(&serviceName, "serviceName", "plugin_center", "service name")
		cmd.PersistentFlags().StringVar(&GroupName, "GroupName", "", "discovery group name")
	})
}
