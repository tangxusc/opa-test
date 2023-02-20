package config

import (
	"github.com/spf13/cobra"
)

var WebPort string

func InitFlags() {
	RegisterFlags(func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringVar(&WebPort, "port", "8080", "web server port")
	})
}
