package login

import (
	"crypto/sha1"
	"encoding/base64"
)

type Loginer interface {
	GetToken(user, pass string) (string, error)
	AddUser(user, pass string) error
	Auth(user, token string) bool
}

type ServiceLogin struct {
	userStorage map[string]string
}

type UserExistError struct{}

func (u UserExistError) Error() string {
	return "User exist !!!"
}

type UserNotFindError struct{}

func (u UserNotFindError) Error() string {
	return "User user not find !!!"
}

type NotAuthError struct{}

func (u NotAuthError) Error() string {
	return "Not Auth error"
}

func (s *ServiceLogin) Auth(user, token string) bool {
	if t, ok := s.userStorage[user]; ok {
		if t == token {
			return true
		}
	}
	return false
}

func (s *ServiceLogin) AddUser(user, pass string) error {
	if s.userStorage == nil {
		s.userStorage = make(map[string]string, 100)
	}
	if _, ok := s.userStorage[user]; ok {
		return UserExistError{}
	}
	sha := sha1.New()
	sum := sha.Sum([]byte(user + pass))
	h := base64.URLEncoding.EncodeToString(sum)

	s.userStorage[user] = h

	return nil
}

func (s *ServiceLogin) GetToken(user, pass string) (string, error) {
	if s.userStorage == nil {
		s.userStorage = make(map[string]string, 100)
		return "", UserNotFindError{}
	}
	if _, ok := s.userStorage[user]; !ok {
		return "", UserNotFindError{}
	}

	sha := sha1.New()
	sum := sha.Sum([]byte(user + pass))
	h := base64.URLEncoding.EncodeToString(sum)
	return h, nil
}
