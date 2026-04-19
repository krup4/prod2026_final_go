package services

import (
	"backend/exceptions"
	"backend/repository"
	"backend/services/inputs"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"strings"
)

type FeatureService struct {
	userRepo        repository.UserRepository
	experimentsRepo repository.ExperimentsRepository
}

func NewFeatureService(userRepo repository.UserRepository, experimentsRepo repository.ExperimentsRepository) *FeatureService {
	return &FeatureService{userRepo: userRepo, experimentsRepo: experimentsRepo}
}

func Normalize(value, min, max float64) float64 {
	return (value - min) / (max - min)
}

func hashToUint64(s string) uint64 {
	sum := sha256.Sum256([]byte(s))
	return binary.BigEndian.Uint64(sum[:8])
}

func hashToPercent(s string) int {
	return int(hashToUint64(s) % 100)
}

func pickVariantByPart(variants []repository.Variant, seed string, all int) *repository.Variant {
	score := int(hashToUint64(seed) % uint64(all))
	cumulative := 0
	for i := range variants {
		cumulative += variants[i].Part
		if score <= cumulative {
			return &variants[i]
		}
	}

	return &variants[len(variants)-1]
}

func (s *FeatureService) GetFeature(ctx context.Context, input inputs.GetFeatureInput) (*repository.Variant, error) {
	flag := strings.TrimSpace(input.Flag)
	userID := strings.TrimSpace(input.UserID)
	defaultValue := strings.TrimSpace(input.DefaultValue)

	if !isValidFlag(flag) || !isValidUUID(userID) || !isValidDefaultValue(defaultValue) {
		return nil, exceptions.ErrInvalidInput
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	participatingCount := 0
	for _, assignment := range user.ExperimentAssignments {
		if assignment.AssignmentType == repository.AssignmentTypeParticipating {
			participatingCount++
		}
		if assignment.Experiment.Flag != flag || assignment.Experiment.Status != "active" {
			continue
		}
		if assignment.AssignmentType == repository.AssignmentTypeParticipating {
			return &assignment.Variant, nil
		}
		if assignment.AssignmentType == repository.AssignmentTypeNotParticipating {
			return &repository.Variant{Name: "fallback", Value: defaultValue}, nil
		}
	}

	experiments, err := s.experimentsRepo.GetByFlag(ctx, flag)
	if err != nil && !errors.Is(err, exceptions.ErrExperimentNotFound) {
		return nil, err
	}

	totalAssignments := len(user.ExperimentAssignments)
	rel := 0.5
	if totalAssignments > 0 {
		rel = float64(participatingCount) / float64(totalAssignments)
	}

	for _, exp := range experiments {
		if exp.Status != "active" {
			continue
		}
		p := float64(exp.Part) / 100
		resP := Normalize(p-0.8*(rel-0.5), -0.8*0.5, 1+0.8*0.5)
		if resP*100 < float64(hashToPercent(userID+":"+exp.ID)) {
			assignment := repository.UserExperimentAssignment{
				UserID:         user.ID,
				ExperimentID:   exp.ID,
				VariantID:      exp.Variants[0].ID,
				AssignmentType: repository.AssignmentTypeNotParticipating,
				Experiment:     exp,
				Variant:        exp.Variants[0],
			}

			if err := s.userRepo.CreateExperimentAssignment(ctx, &assignment); err != nil {
				return nil, err
			}

			return &repository.Variant{Name: "fallback", Value: defaultValue}, nil
		}

		selectedVariant := pickVariantByPart(exp.Variants, userID+":"+exp.ID, exp.Part)
		if selectedVariant == nil {
			continue
		}

		assignment := repository.UserExperimentAssignment{
			UserID:         user.ID,
			ExperimentID:   exp.ID,
			VariantID:      selectedVariant.ID,
			AssignmentType: repository.AssignmentTypeParticipating,
			Experiment:     exp,
			Variant:        *selectedVariant,
		}

		if err := s.userRepo.CreateExperimentAssignment(ctx, &assignment); err != nil {
			return nil, err
		}

		return selectedVariant, nil
	}

	return &repository.Variant{Name: "fallback", Value: defaultValue}, nil
}
