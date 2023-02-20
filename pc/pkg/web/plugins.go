package web

import (
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"plugin_center/pkg/discovery"
)

func (q *pluginQuery) filter(instance *model.Instance) bool {
	metadata := instance.Metadata
	if len(q.Module) > 0 {
		s, ok := metadata["module"]
		if !ok {
			return false
		}
		if q.Module != s {
			return false
		}
	}
	if len(q.PluginType) > 0 {
		pluginType, ok := metadata["plugin_type"]
		if !ok {
			return false
		}
		if q.PluginType != pluginType {
			return false
		}
	}
	if len(q.FilterType) > 0 {
		filterType, ok := metadata["filter_type"]
		if !ok {
			return false
		}
		if q.FilterType != filterType {
			return false
		}
	}
	return true
}

type pluginInfo struct {
	*model.Instance `json:"instance"`
	PluginName      string `json:"plugin_name"`
}

type pluginQuery struct {
	Module     string `form:"module"`
	PluginType string `form:"plugin_type"`
	FilterType string `form:"filter_type"`
}

func init() {
	RegisterHandler(func(engine *gin.Engine) {
		//param: module,plugin_type,filter_type
		engine.GET("/plugins", func(c *gin.Context) {
			query := &pluginQuery{}
			if err := c.ShouldBindQuery(query); err != nil {
				errHandle(err, c)
				return
			}

			client := discovery.GetDiscoveryClient()
			info, err := client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
				NameSpace: discovery.Namespace,
				GroupName: discovery.GroupName,
				PageNo:    0,
				PageSize:  math.MaxUint32 / 10,
			})
			if err != nil {
				errHandle(err, c)
				return
			}
			logrus.Debugf("[web]get all service count:%v", len(info.Doms))
			plugins := make([]*pluginInfo, 0)
			for _, dom := range info.Doms {
				instance, err := client.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
					ServiceName: dom,
					GroupName:   discovery.GroupName,
				})
				if err != nil {
					continue
				}
				if query.filter(instance) {
					plugins = append(plugins, newPlugin(instance))
				}
			}
			c.JSON(http.StatusOK, plugins)
		})
	})
}

func newPlugin(instance *model.Instance) *pluginInfo {
	return &pluginInfo{
		Instance:   instance,
		PluginName: instance.Metadata["plugin_name"],
	}
}
