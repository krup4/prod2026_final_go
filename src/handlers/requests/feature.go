package requests

type GetFeatureRequest struct {
	Flag         string `json:"feature_name" binding:"required"`
	UserID       string `json:"person_id" binding:"required"`
	DefaultValue string `json:"fallback_value" binding:"required"`
}
