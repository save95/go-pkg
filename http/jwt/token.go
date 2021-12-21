package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/save95/go-pkg/http/types"
)

type token struct {
	claims *claims

	issuer   string        // 发行人
	issuedAt *time.Time    // 发行时间
	duration time.Duration // 有效时长
	secret   []byte        // 加密密钥
}

// NewToken 初始化 Token
// 默认发行人为 "go-pkg"，可以通过 SetIssuer 修改；
// 默认有效期为 24h，可以通过 SetDuration 设置有效时长
func NewToken(user types.User) *token {
	var jwtRoles []string
	for i := range user.Roles {
		jwtRoles = append(jwtRoles, user.Roles[i].String())
	}

	return newTokenWith(&claims{
		Account: user.Account,
		UserID:  user.ID,
		Name:    user.Name,
		Roles:   jwtRoles,

		IP:     user.IP,
		Extend: user.Extend,
	})
}

func newTokenWith(c *claims) *token {
	now := time.Now()

	c.Issuer = "go-pkg"
	c.IssuedAt = now.Unix()

	d := 24 * time.Hour
	c.ExpiresAt = now.Add(d).Unix()

	return &token{
		claims:   c,
		issuer:   c.Issuer,
		issuedAt: &now,
		duration: d,
		secret:   jwtSecret,
	}
}

// SetIssuer 设置 token 发行人，默认为 "go-pkg"
// Deprecated
func (t *token) SetIssuer(issuer string) {
	t.issuer = issuer
	t.claims.Issuer = issuer
}

// SetDuration 设置 token 过期时长，默认为 24h
// Deprecated
func (t *token) SetDuration(d time.Duration) {
	t.duration = d
	t.claims.ExpiresAt = t.issuedAt.Add(d).Unix()
}

// SetSecret 设置 token 加密密钥，默认 "go-pkg.JwtSecret"
// Deprecated
func (t *token) SetSecret(secret []byte) {
	if len(secret) == 0 {
		secret = jwtSecret
	}

	t.secret = secret
}

// WithIssuer 设置 token 发行人，默认为 "go-pkg"
func (t *token) WithIssuer(issuer string) *token {
	t.issuer = issuer
	t.claims.Issuer = issuer

	return t
}

// WithDuration 设置 token 过期时长，默认为 24h
func (t *token) WithDuration(d time.Duration) *token {
	t.duration = d
	t.claims.ExpiresAt = t.issuedAt.Add(d).Unix()

	return t
}

// WithSecret 设置 token 加密密钥，默认 "go-pkg.JwtSecret"
func (t *token) WithSecret(secret []byte) *token {
	if len(secret) == 0 {
		secret = jwtSecret
	}

	t.secret = secret

	return t
}

// WithData 设置 token 扩展数据
func (t *token) WithData(key string, val string) *token {
	if len(key) == 0 {
		return t
	}

	if t.claims.Extend == nil {
		t.claims.Extend = make(map[string]string, 0)
	}

	t.claims.Extend[key] = val

	return t
}

func (t *token) User(fun types.ToRole) (*types.User, error) {
	roles, err := t.rolesBy(fun)
	if nil != err {
		return nil, err
	}

	return &types.User{
		ID:    t.claims.UserID,
		Name:  t.claims.Name,
		Roles: roles,

		IP:     t.claims.IP,
		Extend: t.claims.Extend,
	}, nil
}

// rolesBy 通过 fun 函数将 jwt 的用户角色转换
func (t *token) rolesBy(fun types.ToRole) ([]types.IRole, error) {
	var roles []types.IRole
	for _, v := range t.claims.Roles {
		r, err := fun(v)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}

	return roles, nil
}

// IsExpired 是否过期
func (t *token) IsExpired() bool {
	return time.Now().Unix() > t.claims.ExpiresAt
}

// Refresh 刷新 token
func (t *token) Refresh() {
	t.claims.IssuedAt = time.Now().Unix()
	t.claims.ExpiresAt = time.Now().Add(t.duration).Unix()
}

// RefreshNear 自动刷新 token，如果当前时间临近过期时间
func (t *token) RefreshNear(d time.Duration) {
	if time.Now().Unix()+int64(d/time.Second) >= t.claims.ExpiresAt {
		t.Refresh()
	}
}

// ToString 转成 token 字符串
func (t *token) ToString() (string, error) {
	//t.Refresh()
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, t.claims)

	return tokenClaims.SignedString(t.secret)
}
