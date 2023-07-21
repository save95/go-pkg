package restful

func WithMsgHandle(handle func(code int, language string) string) func(*response) {
	return func(r *response) {
		r.msgHandler = handle
	}
}

func WithLanguageKey(key string) func(*response) {
	return func(r *response) {
		r.languageHeaderKey = key
	}
}
