package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// RuleInfo holds the schema definition for the RuleInfo entity.
type RuleInfo struct {
	ent.Schema
}

// Fields of the RuleInfo.
func (RuleInfo) Fields() []ent.Field {
	return []ent.Field{
		field.String("module").Comment("模块").NotEmpty(),
		field.String("plugin_type").Comment("插件类型").NotEmpty(),
		field.String("filter_type").Comment("过滤器类型").NotEmpty(),
		field.String("rule_name").Comment("规则名称").NotEmpty(),
		field.Text("rule_body").Comment("规则内容").NotEmpty(),
		field.Time("create_time").Comment("创建时间").Default(time.Now()),
		field.Bool("enable").Comment("是否启用").Default(true),
	}
}

// Edges of the RuleInfo.
func (RuleInfo) Edges() []ent.Edge {
	return nil
}
