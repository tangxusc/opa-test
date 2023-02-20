package web

import (
	"github.com/spf13/cobra"
	"os"
	"plugin_center/pkg/config"
)

var tempdir string

func InitFlags() {
	config.RegisterFlags(func(c *cobra.Command) {
		c.PersistentFlags().StringVar(&tempdir, "tempdir", os.TempDir(), "")
	})
}
