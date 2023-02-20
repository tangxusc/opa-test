package store

import (
	"github.com/spf13/cobra"
	"plugin_center/pkg/config"
)

var mysqldsn string

func InitFlags() {
	config.RegisterFlags(func(c *cobra.Command) {
		//characterEncoding=utf-8&serverTimezone=Asia/Shanghai
		c.PersistentFlags().StringVar(&mysqldsn, "mysqldsn", "root:123456@tcp(127.0.0.1:3306)/auth?parseTime=True", "mysqlsdn")
	})
}
