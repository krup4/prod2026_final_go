package services

import (
	"backend/exceptions"
	"backend/repository"
	"backend/services/inputs"
	"context"
	"errors"
	"strings"
)

type ExperimentsService struct {
	experimentsRepo repository.ExperimentsRepository
}

func NewExperimentsService(experimentsRepo repository.ExperimentsRepository) *ExperimentsService {
	return &ExperimentsService{experimentsRepo: experimentsRepo}
}

func (s *ExperimentsService) CreateExperiment(ctx context.Context, input inputs.CreateExperimentInput) (*repository.Experiment, error) {
	identifier := strings.TrimSpace(input.Identifier)
	flag := strings.TrimSpace(input.Flag)
	name := strings.TrimSpace(input.Name)
	status := strings.TrimSpace(input.Status)

	if !isValidIdentifier(identifier) || !isValidFlag(flag) || !isValidName(name) || !isValidStatus(status) || input.Version <= 0 || input.Part <= 0 || len(input.Variants) == 0 {
		return nil, exceptions.ErrInvalidInput
	}

	experiment := &repository.Experiment{
		Identifier: identifier,
		Flag:       flag,
		Name:       name,
		Status:     status,
		Version:    input.Version,
		Part:       input.Part,
		Variants:   make([]repository.Variant, 0, len(input.Variants)),
	}

	sumParts := 0
	cntControl := 0
	for _, v := range input.Variants {
		variantName := strings.TrimSpace(v.Name)
		variantValue := strings.TrimSpace(v.Value)
		if !isValidName(variantName) || !isValidVariantValue(variantValue) || v.Part <= 0 || v.Part > 100 || v.IsControl == nil {
			return nil, exceptions.ErrInvalidInput
		}

		if *v.IsControl {
			cntControl++
		}

		sumParts += v.Part

		experiment.Variants = append(experiment.Variants, repository.Variant{
			Name:      variantName,
			Value:     variantValue,
			Part:      v.Part,
			IsControl: v.IsControl,
		})
	}

	if sumParts != input.Part || cntControl != 1 {
		return nil, exceptions.ErrInvalidInput
	}

	_, err := s.experimentsRepo.GetByIdentifier(ctx, identifier)
	if err == nil {
		return nil, exceptions.ErrExperimentExists
	}

	if !errors.Is(err, exceptions.ErrExperimentNotFound) {
		return nil, err
	}

	experimentsByFlag, err := s.experimentsRepo.GetByFlag(ctx, flag)
	if err != nil && !errors.Is(err, exceptions.ErrExperimentNotFound) {
		return nil, err
	}
	for _, exp := range experimentsByFlag {
		if (status == "active" || status == "pause") && (exp.Status == "active" || exp.Status == "pause") {
			return nil, exceptions.ErrActiveExperiment
		}
	}

	if sumParts != input.Part || cntControl != 1 {
		return nil, exceptions.ErrInvalidInput
	}

	if err := s.experimentsRepo.Create(ctx, experiment); err != nil {
		return nil, err
	}

	return experiment, nil
}

func (s *ExperimentsService) ChangeStatus(ctx context.Context, input inputs.ChangeStatusInput) (*repository.Experiment, error) {
	identifier := strings.TrimSpace(input.Identifier)
	status := strings.TrimSpace(input.Status)

	if !isValidIdentifier(identifier) || !isValidStatus(status) {
		return nil, exceptions.ErrInvalidInput
	}

	experiment, err := s.experimentsRepo.GetByIdentifier(ctx, identifier)
	if err != nil {
		return nil, err
	}

	experiments, err := s.experimentsRepo.GetByFlag(ctx, experiment.Flag)
	if err != nil && !errors.Is(err, exceptions.ErrExperimentNotFound) {
		return nil, err
	}

	for _, exp := range experiments {
		if exp.Identifier != experiment.Identifier && ((status == "active" || status == "pause") && (exp.Status == "active" || exp.Status == "pause")) {
			return nil, exceptions.ErrActiveExperiment
		}
	}

	experiment, err = s.experimentsRepo.UpdateStatus(ctx, experiment, status)

	if err != nil {
		return nil, err
	}

	return experiment, nil
}
