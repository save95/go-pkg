package types

type IRole interface {
	String() string
}

// RoleConvert 转换成角色 IRole
type RoleConvert func(role string) (IRole, error)
