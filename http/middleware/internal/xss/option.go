package xss

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/save95/go-pkg/http/xss"
)

// WithGlobalPolicy 指定全局过滤策略
func WithGlobalPolicy(p xss.Policy) Option {
	return func(h *handler) {
		h.policy = h.makePolicy(p)
	}
}

// WithGlobalFieldPolicy 指定全局字段过滤策略
func WithGlobalFieldPolicy(p xss.Policy, fields ...string) Option {
	return func(h *handler) {
		if h.fieldRules == nil {
			h.fieldRules = make(map[string]*bluemonday.Policy, 0)
		}
		for _, field := range fields {
			h.fieldRules[field] = h.makePolicy(p)
		}
	}
}

// WithDebug 设置调试模式
func WithDebug() func(h *handler) {
	return func(h *handler) {
		h.debug = true
	}
}

// WithGlobalSkipFields 指定全局忽略字段
func WithGlobalSkipFields(fields ...string) Option {
	return func(h *handler) {
		h.skipField = h.makeSkipFields(fields)
	}
}

// WithRoutePolicy 指定路由策略
// routeRule 路由规则，如果路由包含该字符串则匹配成功
func WithRoutePolicy(routeRule string, policy xss.Policy, skipFields ...string) Option {
	return func(h *handler) {
		if policy == xss.PolicyNone {
			h.skipRoutes[routeRule] = struct{}{}
			return
		}

		h.routePolicies[routeRule] = &xssRuleItem{
			policy:    h.makePolicy(policy),
			skipField: h.makeSkipFields(skipFields),
		}
	}
}

// WithRouteFieldPolicy 指定路由的字段策略
// routeRule 路由规则，如果路由包含该字符串则匹配成功
func WithRouteFieldPolicy(routeRule string, policy xss.Policy, fields ...string) Option {
	return func(h *handler) {
		if policy == xss.PolicyNone {
			h.skipRoutes[routeRule] = struct{}{}
			return
		}

		rp, ok := h.routePolicies[routeRule]
		if !ok {
			rp = &xssRuleItem{
				policy:     h.policy,
				fieldRules: make(map[string]*bluemonday.Policy, 0),
				skipField:  make(map[string]struct{}, 0),
			}
		}
		for _, field := range fields {
			rp.fieldRules[field] = h.makePolicy(policy)
		}
		h.routePolicies[routeRule] = rp
	}
}
