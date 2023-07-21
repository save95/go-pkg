package xss

// Policy XSS 策略
type Policy uint8

const (
	// PolicyNone 无过滤
	PolicyNone Policy = iota
	// PolicyStrict 过滤所有HTML元素及其属性
	PolicyStrict
	// PolicyUGC 过滤不安全的HTML元素和属性，如：iframes, object, embed, styles, script
	PolicyUGC
)
