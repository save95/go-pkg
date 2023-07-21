package xss

import "github.com/microcosm-cc/bluemonday"

type Option func(*handler)

type xssRuleItem struct {
	policy     *bluemonday.Policy
	fieldRules map[string]*bluemonday.Policy
	skipField  map[string]struct{}
}
