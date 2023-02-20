package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"os/signal"
	"plugin_center/pkg/config"
	"plugin_center/pkg/discovery"
	"plugin_center/pkg/store"
	"plugin_center/pkg/web"
	"time"
)

func NewCommand() (*cobra.Command, context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.TODO())
	command := &cobra.Command{
		Use:   ``,
		Short: ``,
		Long:  ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Kill)
				<-c
				cancelFunc()
			}()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.SetLevel(logrus.TraceLevel)
			rand.Seed(time.Now().UnixNano())
			err := discovery.Connect(ctx)
			if err != nil {
				return err
			}
			if err := store.ConnectMysql(ctx); err != nil {
				return err
			}
			err = web.Start(ctx)
			if err != nil {
				return err
			}

			<-ctx.Done()
			return nil
		},
	}

	discovery.InitFlags()
	web.InitFlags()
	store.InitFlags()
	config.InitFlags()
	viper.AutomaticEnv()
	viper.AddConfigPath(`.`)
	config.BuildFlags(command)
	_ = viper.BindPFlags(command.Flags())

	return command, ctx, cancelFunc
}

func main() {
	command, _, _ := NewCommand()
	if err := command.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
