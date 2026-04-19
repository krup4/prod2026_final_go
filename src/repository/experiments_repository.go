package repository

import (
	"context"
	"errors"

	"backend/exceptions"

	"gorm.io/gorm"
)

type ExperimentsRepository interface {
	Create(ctx context.Context, experiment *Experiment) error
	GetByIdentifier(ctx context.Context, identifier string) (*Experiment, error)
	GetByFlag(ctx context.Context, flag string) ([]Experiment, error)
	UpdateStatus(ctx context.Context, experiment *Experiment, status string) (*Experiment, error)
}

type experimentsRepository struct {
	db *gorm.DB
}

func NewExperimentsRepository(db *gorm.DB) ExperimentsRepository {
	return &experimentsRepository{db: db}
}

func (r *experimentsRepository) Create(ctx context.Context, experiment *Experiment) error {
	return r.db.WithContext(ctx).Create(experiment).Error
}

func (r *experimentsRepository) GetByIdentifier(ctx context.Context, identifier string) (*Experiment, error) {
	var experiment Experiment
	err := r.db.WithContext(ctx).Preload("Variants").Where("identifier = ?", identifier).First(&experiment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exceptions.ErrExperimentNotFound
	}
	if err != nil {
		return nil, err
	}

	return &experiment, nil
}

func (r *experimentsRepository) GetByFlag(ctx context.Context, flag string) ([]Experiment, error) {
	var experiments []Experiment
	err := r.db.WithContext(ctx).Preload("Variants").Where("flag = ?", flag).Find(&experiments).Error
	if err != nil {
		return nil, err
	}
	if len(experiments) == 0 {
		return nil, exceptions.ErrExperimentNotFound
	}

	return experiments, nil
}

func (r *experimentsRepository) UpdateStatus(ctx context.Context, experiment *Experiment, status string) (*Experiment, error) {
	if err := r.db.WithContext(ctx).Model(experiment).Update("status", status).Error; err != nil {
		return nil, err
	}

	experiment.Status = status
	return experiment, nil
}
