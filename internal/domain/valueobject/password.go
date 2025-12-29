package valueobject

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordToShort      = errors.New("password must be at least 8 characters long")
	ErrPasswordTooWeak      = errors.New("password is too weak")
	ErrPasswordHashFailed   = errors.New("failed to hash password")
	ErrPasswordVerifyFailed = errors.New("failed to verify password")
)

// Password 密码值对象
// 值对象封装了密码的创建和验证逻辑
// 存储的是加密后的哈希值，而不是铭文
type Password struct {
	hash string
}

func NewPassword(plaintext string) (Password, error) {
	if len(plaintext) < 8 {
		return Password{}, ErrPasswordToShort
	}

	if !isStrongPassword(plaintext) {
		return Password{}, ErrPasswordTooWeak
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, ErrPasswordHashFailed
	}

	return Password{hash: string(hashedBytes)}, nil
}

func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

func (p *Password) Hash() string {
	return p.hash
}

// isStrongPassword 用于判断密码强度，要求包含大写字母、小写字母和数字。
// 解释：
// 1. 遍历密码字符串中的每个字符：
//   - 如果是大写字母，则 hashUpper 设为 true。
//   - 如果是小写字母，则 hashLower 设为 true。
//   - 如果是数字，则 hashDigit 设为 true。
//
// 2. 最终如果三者都包含（即全为 true），返回 true，否则返回 false。
// 这样确保密码同时包含了大写、小写字母和数字，提高密码安全性。
func isStrongPassword(plaintext string) bool {
	var hashUpper, hashLower, hashDigit bool

	for _, char := range plaintext {
		if unicode.IsUpper(char) {
			hashUpper = true
		} else if unicode.IsLower(char) {
			hashLower = true
		} else if unicode.IsDigit(char) {
			hashDigit = true
		}
	}

	return hashUpper && hashLower && hashDigit
}
