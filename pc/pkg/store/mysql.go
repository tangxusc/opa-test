package store

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"plugin_center/ent"
)

var Client *ent.Client

func ConnectMysql(ctx context.Context) error {
	var err error
	Client, err = ent.Open("mysql", mysqldsn)
	if err != nil {
		return err
	}
	err = Client.Schema.Create(ctx)
	if err != nil {
		return err
	}
	go func() {
		select {
		case <-ctx.Done():
			_ = Client.Close()
		}
	}()
	return nil
}
