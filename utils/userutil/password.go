package userutil

import (
	util "github.com/save95/go-utils/userutil"
)

type hasher struct {
}

// NewHasher 密码加密器
// Deprecated
func NewHasher() *hasher {
	return &hasher{}
}

// Sum 加密密码
func (h hasher) Sum(password string) (string, error) {
	return util.NewHasher().Sum(password)
}

// Check 检查加密密码
func (h hasher) Check(input, password string) bool {
	return util.NewHasher().Check(input, password)
}
