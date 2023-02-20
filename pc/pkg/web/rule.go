package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"plugin_center/ent"
	"plugin_center/ent/ruleinfo"
	"plugin_center/pkg/store"
	"time"
)

type Page struct {
	//每页大小
	Size uint `json:"size"`
	//总记录数量
	Total uint `json:"total"`
	//总页数
	Pages uint `json:"pages"`
	//当前页
	Current uint `json:"current"`
	//数据
	Records []*ent.RuleInfo `json:"records"`
}

func (p *Page) SetTotal(count int) {
	p.Total = uint(count)
	if count < 1 {
		return
	}
	i := uint(count) / p.Size
	if uint(count)%p.Size > 0 {
		i = i + 1
	}
	p.Pages = i
}

type ruleListRequest struct {
	ent.RuleInfo `json:",inline"`
	*Page        `json:",inline"`
}

func (r *ruleListRequest) execute(c *gin.Context, query *ent.RuleInfoQuery) {
	//注入查询条件
	if len(r.Module) > 0 {
		query.Where(ruleinfo.ModuleEQ(r.Module))
	}
	if len(r.PluginType) > 0 {
		query.Where(ruleinfo.PluginTypeEQ(r.PluginType))
	}
	if len(r.FilterType) > 0 {
		query.Where(ruleinfo.FilterTypeEQ(r.FilterType))
	}
	if len(r.RuleName) > 0 {
		query.Where(ruleinfo.RuleNameEQ(r.RuleName))
	}
	query.Order(ent.Desc(ruleinfo.FieldCreateTime))

	count := query.CountX(c)
	r.Page.SetTotal(count)
	all := query.AllX(c)
	r.Page.Records = all
}

func init() {
	RegisterHandler(func(engine *gin.Engine) {
		//store rule
		engine.POST("/rule/list", func(c *gin.Context) {
			listRequest := &ruleListRequest{Page: &Page{
				Size:    10,
				Current: 1,
			}}
			err := c.ShouldBind(listRequest)
			if err != nil {
				errHandle(err, c)
				return
			}
			listRule(c, listRequest)
		})

		//store rule
		engine.POST("/rule/save", func(c *gin.Context) {
			saveRequest := &ent.RuleInfo{}
			err := c.ShouldBind(saveRequest)
			if err != nil {
				errHandle(err, c)
				return
			}
			deleteRule(c, saveRequest)
			saveRule(c, saveRequest)

		})
		engine.POST("/rule/delete", func(c *gin.Context) {
			request := &ent.RuleInfo{}
			err := c.ShouldBind(request)
			if err != nil {
				errHandle(err, c)
				return
			}
			deleteRule(c, request)
			c.JSON(http.StatusOK, "")
		})
	})
}

func saveRule(c *gin.Context, request *ent.RuleInfo) {
	x := store.Client.RuleInfo.Create().
		SetModule(request.Module).
		SetPluginType(request.PluginType).
		SetFilterType(request.FilterType).
		SetRuleName(request.RuleName).
		SetRuleBody(request.RuleBody).
		SetCreateTime(time.Now()).
		SetEnable(request.Enable).SaveX(c)
	c.JSON(http.StatusOK, x)
}

func deleteRule(c *gin.Context, request *ent.RuleInfo) {
	if request.ID > 0 {
		store.Client.RuleInfo.DeleteOneID(request.ID).ExecX(c)
	}
}

func listRule(c *gin.Context, request *ruleListRequest) {
	query := store.Client.RuleInfo.Query()
	request.execute(c, query)

	c.JSON(http.StatusOK, request.Page)
}
