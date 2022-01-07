package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"order-validation-v2/internal/entity"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetUserbyID(id string) (*entity.User, error) {
	u, err := s.repo.GetbyID(id)

	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("not found")
	}

	return u, nil
}

func (s *Service) CreateUser(username string, email string, password string, role string) (string, error) {
	u := entity.NewUser(email, username, password, role)
	return s.repo.Create(u)

}

func (s *Service) GetUserbyUsername(username string) (*entity.User, error) {
	u, err := s.repo.GetbyUsername(username)

	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("not found")
	}

	return u, nil
}

func (s *Service) SearchUser(query string) ([]*entity.User, error) {
	users, err := s.repo.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("not found")
	}
	return users, nil
}

func (s *Service) ListUsers() ([]*entity.User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("not found")
	}
	return users, nil
}

func (s *Service) DeleteUser(id string) error {
	_, err := s.GetUserbyID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) UpdateUser(u *entity.User) error {
	return s.repo.Update(u)
}

func (s *Service) Login(username string, password string) (string, string, bool, error) {
	u, err := s.repo.GetbyUsername(username)
	if err != nil {
		return username, "", false, err
	}
	if u == nil {
		return username, "", false, errors.New("Username/Password wrong")
	}
	sum := sha256.Sum256([]byte(password))
	if u.Password == fmt.Sprintf("\\x%x", sum) {
		return u.ID, u.UserRole, true, nil
	}
	return username, u.UserRole, false, errors.New("Username/Password wrong")

}

func (s *Service) ValidateAndRetrieveUser(userID string, password string) (bool, *entity.User, error) {
	u, err := s.repo.GetbyID(userID)
	if err != nil {
		return false, nil, err
	}
	if u == nil {
		return false, nil, errors.New("UserID not found")
	}
	sum := sha256.Sum256([]byte(password))
	fmt.Println(sum)
	fmt.Println(fmt.Sprintf("\\x%x", sum))
	if u.Password == fmt.Sprintf("\\x%x", sum) {
		return true, u, nil
	}
	return false, nil, nil
}

func (s *Service) ValidateUsername(username string) (bool, error) {
	return s.repo.CheckUsername(username)
}
