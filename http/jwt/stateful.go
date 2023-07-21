package jwt

// StatefulStore 状态存储
type StatefulStore interface {
	// Save token 状态存储器
	Save(userID uint, token string, expireTs int64) error
	// Check token 状态检查器
	Check(userID uint, token string) error
	// Remove 删除指定 token
	Remove(userID uint, token string) error
	// Clean 清理用户的所有 token
	Clean(userID uint) error
}
