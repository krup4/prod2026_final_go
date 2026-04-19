package repository

import (
	"backend/exceptions"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	CreateExperimentAssignment(ctx context.Context, assignment *UserExperimentAssignment) error
	GetByLogin(ctx context.Context, login string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) CreateExperimentAssignment(ctx context.Context, assignment *UserExperimentAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("login = ?", login).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exceptions.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Preload("ExperimentAssignments.Experiment").
		Preload("ExperimentAssignments.Variant").
		Where("id = ?", id).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exceptions.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) List(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).
		Preload("ExperimentAssignments.Experiment").
		Preload("ExperimentAssignments.Variant").
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
