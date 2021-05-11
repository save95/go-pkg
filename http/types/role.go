package types

type IRole interface {
	String() string
}

// ToRole 转换成角色 IRole
type ToRole func(role string) (IRole, error)
