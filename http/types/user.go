package types

// User 用户基本信息
type User struct {
	ID      uint
	Account string
	Name    string
	Roles   []IRole
}
