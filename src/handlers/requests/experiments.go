package requests

type CreateVariantRequest struct {
	Name      string `json:"name" binding:"required"`
	Value     string `json:"value" binding:"required"`
	Part      int    `json:"part" binding:"required"`
	IsControl *bool  `json:"is_control" binding:"required"`
}

type CreateExperimentRequest struct {
	Identifier string                 `json:"identifier" binding:"required"`
	Flag       string                 `json:"flag" binding:"required"`
	Name       string                 `json:"name" binding:"required"`
	Status     string                 `json:"status" binding:"required"`
	Version    int                    `json:"version" binding:"required"`
	Part       int                    `json:"part" binding:"required"`
	Variants   []CreateVariantRequest `json:"variants" binding:"required,min=1,dive"`
}

type ChangeStatusRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Status     string `json:"status" binding:"required"`
}
