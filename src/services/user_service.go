package services

import (
	"context"
	"errors"
	"strings"

	"backend/exceptions"
	"backend/repository"
	"backend/services/inputs"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, input inputs.CreateUserInput) (*repository.User, error) {
	login := strings.TrimSpace(input.Login)
	role := strings.TrimSpace(input.Role)
	if !isValidLogin(login) || !isValidRole(role) {
		return nil, exceptions.ErrInvalidInput
	}

	_, err := s.userRepo.GetByLogin(ctx, login)
	if err == nil {
		return nil, exceptions.ErrUserExists
	}
	if !errors.Is(err, exceptions.ErrUserNotFound) {
		return nil, err
	}

	user := &repository.User{
		Login:                 login,
		Role:                  role,
		ExperimentAssignments: make([]repository.UserExperimentAssignment, 0),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUsersList(ctx context.Context) ([]repository.User, error) {
	users, err := s.userRepo.List(ctx)

	if err != nil {
		return nil, err
	}

	for i := range users {
		if users[i].ExperimentAssignments == nil {
			users[i].ExperimentAssignments = make([]repository.UserExperimentAssignment, 0)
		}
	}

	return users, nil
}
