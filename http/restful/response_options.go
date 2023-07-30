package restful

import "github.com/save95/xerror/xcode"

// WithMsgHandle 指定错误消息处理器。
// 主要应用于多国语言展示错误信息，需要配合 WithLanguageKey 使用。
// Deprecated. 使用 WithErrorMsgHandle 替代
func WithMsgHandle(handle func(code int, language string) string) func(*response) {
	return func(r *response) {
		r.msgHandler = handle
	}
}

// WithLanguageKey 指定语言 header key
// Deprecated. 使用 WithErrorMsgHandle 替代
func WithLanguageKey(key string) func(*response) {
	return func(r *response) {
		r.languageHeaderKey = key
	}
}

// WithErrorMsgHandle 指定错误消息处理器。
// 主要应用于多国语言展示错误信息。
// 其中，`code` 为错误码；`language` 为标准的 i18n 标识
func WithErrorMsgHandle(languageHeaderKey string, handle func(code int, language string) string) func(*response) {
	return func(r *response) {
		r.languageHeaderKey = languageHeaderKey
		r.msgHandler = handle
	}
}

// WithShowXCode 设置需要对客户端展示的错误码。
// 默认情况下，该值为空，则表示所有错误码均向用户展示（设置为空，亦如此）；
func WithShowXCode(xcodes ...xcode.XCode) func(*response) {
	return func(r *response) {
		if len(xcodes) == 0 {
			r.showErrorCodes = make([]int, 0)
			return
		}

		for _, err := range xcodes {
			if r.showErrorCodes == nil {
				r.showErrorCodes = make([]int, 0)
			}
			r.showErrorCodes = append(r.showErrorCodes, err.Code())
		}
	}
}
