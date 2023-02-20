// Code generated by ent, DO NOT EDIT.

package ruleinfo

import (
	"time"
)

const (
	// Label holds the string label denoting the ruleinfo type in the database.
	Label = "rule_info"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldModule holds the string denoting the module field in the database.
	FieldModule = "module"
	// FieldPluginType holds the string denoting the plugin_type field in the database.
	FieldPluginType = "plugin_type"
	// FieldFilterType holds the string denoting the filter_type field in the database.
	FieldFilterType = "filter_type"
	// FieldRuleName holds the string denoting the rule_name field in the database.
	FieldRuleName = "rule_name"
	// FieldRuleBody holds the string denoting the rule_body field in the database.
	FieldRuleBody = "rule_body"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldEnable holds the string denoting the enable field in the database.
	FieldEnable = "enable"
	// Table holds the table name of the ruleinfo in the database.
	Table = "rule_infos"
)

// Columns holds all SQL columns for ruleinfo fields.
var Columns = []string{
	FieldID,
	FieldModule,
	FieldPluginType,
	FieldFilterType,
	FieldRuleName,
	FieldRuleBody,
	FieldCreateTime,
	FieldEnable,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// ModuleValidator is a validator for the "module" field. It is called by the builders before save.
	ModuleValidator func(string) error
	// PluginTypeValidator is a validator for the "plugin_type" field. It is called by the builders before save.
	PluginTypeValidator func(string) error
	// FilterTypeValidator is a validator for the "filter_type" field. It is called by the builders before save.
	FilterTypeValidator func(string) error
	// RuleNameValidator is a validator for the "rule_name" field. It is called by the builders before save.
	RuleNameValidator func(string) error
	// RuleBodyValidator is a validator for the "rule_body" field. It is called by the builders before save.
	RuleBodyValidator func(string) error
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime time.Time
	// DefaultEnable holds the default value on creation for the "enable" field.
	DefaultEnable bool
)