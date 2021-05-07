package types

// User 用户基本信息
type User struct {
	ID    uint
	Name  string
	Roles []IRole
}
