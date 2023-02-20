package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"plugin_center/pkg/config"
)

type Handler func(engine *gin.Engine)

var handlers = make([]Handler, 0)

func RegisterHandler(h Handler) {
	handlers = append(handlers, h)
}

func Start(ctx context.Context) error {
	engine := gin.Default()
	for _, handler := range handlers {
		handler(engine)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.WebPort),
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		select {
		case <-ctx.Done():
			logrus.Debugf("[web]gin server exiting...")
			if err := srv.Shutdown(ctx); err != nil {
				logrus.Fatal("Server Shutdown:", err)
			}
		}
	}()
	return nil
}

func errHandle(err error, c *gin.Context) {
	c.JSON(http.StatusInternalServerError, err)
}
