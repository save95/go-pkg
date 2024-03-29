package types

// User 用户基本信息
type User struct {
	ID      uint
	Account string
	Name    string
	Roles   []IRole

	IP     string            // 用户登录ID
	Extend map[string]string // 扩展信息
}

func (u *User) GetID() uint {
	if nil == u {
		return 0
	}

	return u.ID
}

func (u *User) GetAccount() string {
	if nil == u {
		return ""
	}

	return u.Account
}

func (u *User) GetName() string {
	if nil == u {
		return ""
	}

	return u.Name
}

func (u *User) GetRoles() []IRole {
	if nil == u {
		return []IRole{}
	}

	return u.Roles
}

func (u *User) Is(role IRole, others ...IRole) bool {
	if u == nil || u.Roles == nil || len(u.Roles) == 0 {
		return false
	}

	// 已存在的角色
	bes := make(map[IRole]struct{}, 0)
	for _, role := range u.Roles {
		bes[role] = struct{}{}
	}

	roles := []IRole{role}
	roles = append(roles, others...)
	for _, role := range roles {
		if _, ok := bes[role]; ok {
			return true
		}
	}

	return false
}

func (u *User) GetIP() string {
	if nil == u {
		return ""
	}

	return u.IP
}

func (u *User) GetExtend() map[string]string {
	if u == nil || u.Extend == nil {
		return make(map[string]string, 0)
	}

	return u.Extend
}

func (u *User) GetExtendValue(field string) string {
	extend := u.GetExtend()
	return extend[field]
}
