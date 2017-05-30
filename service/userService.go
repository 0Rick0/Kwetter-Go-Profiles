package service

import (
	"../types"
)

func (s *Service) RemoveUser(user types.User) bool {
	var count int
	s.db.Model(&types.User{}).Where("username = ?", user.Username).Count(&count)
	if count < 1 {
		return false
	}
	s.db.Delete(&user)
	return true
}

func (s *Service) GetUserByUsername(username string) *types.User {
	var user types.User
	s.db.FirstOrInit(&user, types.User{Username: username})
	return &user
}

func (s *Service) SetProfilePicture(username string, mime string) bool {
	user := s.GetUserByUsername(username)
	user.MimeType = mime
	s.db.Save(&user)
	return true
}
