package repository

import (
	"crypto/rand"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const FallbackVariantName = "__fallback__"
const (
	AssignmentTypeParticipating    = "participating"
	AssignmentTypeNotParticipating = "not_participating"
)

type User struct {
	ID                    string                     `gorm:"type:uuid;primaryKey"`
	Login                 string                     `gorm:"size:64;uniqueIndex;not null"`
	Role                  string                     `gorm:"size:32;not null"`
	ExperimentAssignments []UserExperimentAssignment `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"experiment_assignments"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type UserExperimentAssignment struct {
	ID             string     `gorm:"type:uuid;primaryKey"`
	UserID         string     `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_experiment_assignment"`
	ExperimentID   string     `gorm:"type:uuid;not null;index;uniqueIndex:idx_user_experiment_assignment"`
	VariantID      string     `gorm:"type:uuid;not null"`
	AssignmentType string     `gorm:"size:32;not null;default:participating;index" json:"assignment_type"`
	Experiment     Experiment `gorm:"foreignKey:ExperimentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"experiment"`
	Variant        Variant    `gorm:"foreignKey:VariantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"variant"`
}

type Variant struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	ExperimentID string `gorm:"type:uuid;not null;index"`
	Name         string `gorm:"size:255;not null"`
	Value        string `gorm:"size:255;not null"`
	Part         int    `gorm:"not null"`
	IsControl    *bool  `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Experiment struct {
	ID         string    `gorm:"type:uuid;primaryKey"`
	Identifier string    `gorm:"size:255;not null"`
	Flag       string    `gorm:"size:255;not null"`
	Name       string    `gorm:"size:255;not null"`
	Status     string    `gorm:"size:255;not null"`
	Version    int       `gorm:"not null"`
	Part       int       `gorm:"not null"`
	Variants   []Variant `gorm:"foreignKey:ExperimentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (m *User) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		id, err := newUUID()
		if err != nil {
			return err
		}
		m.ID = id
	}
	return nil
}

func (m *UserExperimentAssignment) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		id, err := newUUID()
		if err != nil {
			return err
		}
		m.ID = id
	}
	if m.AssignmentType == "" {
		m.AssignmentType = AssignmentTypeParticipating
	}
	return nil
}

func (m *Variant) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		id, err := newUUID()
		if err != nil {
			return err
		}
		m.ID = id
	}
	return nil
}

func (m *Experiment) BeforeCreate(_ *gorm.DB) error {
	if m.ID == "" {
		id, err := newUUID()
		if err != nil {
			return err
		}
		m.ID = id
	}
	return nil
}

func newUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}
