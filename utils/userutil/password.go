package userutil

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type hasher struct {
}

// NewHasher 密码加密器
func NewHasher() *hasher {
	return &hasher{}
}

func (h hasher) hash(str string) string {
	s := sha1.New()
	bs := s.Sum([]byte(str))
	return fmt.Sprintf("%+x", string(bs))
}

// Sum 加密密码
func (h hasher) Sum(password string) (string, error) {
	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return "", errors.New("password empty")
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if nil != err {
		return "", err
	}
	return string(bs), nil
}

// Check 检查加密密码
func (h hasher) Check(input, password string) bool {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return false
	}

	return nil == bcrypt.CompareHashAndPassword([]byte(password), []byte(input))
}
