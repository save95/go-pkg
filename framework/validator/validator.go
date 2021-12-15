package validator

// IValidator 验证器约定
type IValidator interface {
	Validate() error
}
