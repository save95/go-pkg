package types

// BodyProperty 响应正文属性
type BodyProperty string

const (
	BodyPropertyRaw  BodyProperty = "raw"
	BodyPropertyText BodyProperty = "text"
	BodyPropertyHtml BodyProperty = "html"
	BodyPropertyFull BodyProperty = "full"
)

func (bp BodyProperty) Verify() bool {
	switch bp {
	case BodyPropertyRaw, BodyPropertyText, BodyPropertyHtml, BodyPropertyFull:
		return true
	}

	return false
}
