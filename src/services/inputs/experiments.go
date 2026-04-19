package inputs

import "backend/repository"

type CreateExperimentInput struct {
	Identifier string
	Flag       string
	Name       string
	Status     string
	Version    int
	Part       int
	Variants   []repository.Variant
}

type ChangeStatusInput struct {
	Identifier string
	Status     string
}
