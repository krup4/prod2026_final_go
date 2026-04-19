package handlers

import (
	"backend/exceptions"
	"backend/handlers/requests"
	"backend/services"
	"backend/services/inputs"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
)

type FeatureHandler struct {
	featureService *services.FeatureService
}

func NewFeatureHandler(featureService *services.FeatureService) *FeatureHandler {
	return &FeatureHandler{featureService: featureService}
}

func (h *FeatureHandler) GetFeature(c *gin.Context) {
	var req requests.GetFeatureRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json payload"})
		return
	}

	variant, err := h.featureService.GetFeature(c.Request.Context(), inputs.GetFeatureInput{
		Flag:         req.Flag,
		UserID:       req.UserID,
		DefaultValue: req.DefaultValue,
	})

	if err != nil {
		if errors.Is(err, exceptions.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inputs"})
			return
		}

		if errors.Is(err, exceptions.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create experiment"})
		return
	}

	c.JSON(http.StatusOK, variant)
}
