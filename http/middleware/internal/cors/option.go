package cors

import "time"

type Option func(*handler)

func WithAllowOriginFunc(fun func(origin string) bool) Option {
	return func(ch *handler) {
		ch.allowOriginFunc = fun
	}
}

func WithAllowMethods(methods ...string) Option {
	return func(ch *handler) {
		if len(methods) > 0 {
			ch.allowMethods = append(ch.allowMethods, methods...)
		}
	}
}

func WithAllowHeaders(keys ...string) Option {
	return func(ch *handler) {
		if len(keys) > 0 {
			ch.allowHeaders = append(ch.allowHeaders, keys...)
		}
	}
}

func WithExposeHeaders(keys ...string) Option {
	return func(ch *handler) {
		if len(keys) > 0 {
			ch.exposeHeaders = append(ch.exposeHeaders, keys...)
		}
	}
}

func WithMaxAge(d time.Duration) Option {
	return func(ch *handler) {
		if d > 0 {
			ch.maxAge = d
		}
	}
}
