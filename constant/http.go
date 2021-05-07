package constant

const (
	DefaultRequestLimit = 20
	DefaultPageSize     = 20

	HttpCustomContextKey = "_custom_ctx"
	HttpTokenClaimsKey   = "_access_jwt_claims"
	HttpTokenHeaderKey   = "X-Token"

	HttpCustomRawRequestBodyKey = "_ctx_request_body"
)
