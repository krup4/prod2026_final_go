package handlers

import (
	"backend/exceptions"
	"backend/handlers/requests"
	"backend/repository"
	"backend/services"
	"backend/services/inputs"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ExperimentsHandler struct {
	experimentsService *services.ExperimentsService
}

func NewExperimentsHandler(experimentsService *services.ExperimentsService) *ExperimentsHandler {
	return &ExperimentsHandler{experimentsService: experimentsService}
}

func (h *ExperimentsHandler) Create(c *gin.Context) {
	var req requests.CreateExperimentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json payload"})
		return
	}

	exp := inputs.CreateExperimentInput{
		Identifier: req.Identifier,
		Flag:       req.Flag,
		Name:       req.Name,
		Status:     req.Status,
		Version:    req.Version,
		Part:       req.Part,
		Variants:   make([]repository.Variant, 0, len(req.Variants)),
	}

	for _, v := range req.Variants {
		exp.Variants = append(exp.Variants, repository.Variant{
			Name:      v.Name,
			Value:     v.Value,
			Part:      v.Part,
			IsControl: v.IsControl,
		})
	}

	experiment, err := h.experimentsService.CreateExperiment(c.Request.Context(), exp)

	if err != nil {
		if errors.Is(err, exceptions.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inputs"})
			return
		}

		if errors.Is(err, exceptions.ErrExperimentExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "experiment already exists"})
			return
		}

		if errors.Is(err, exceptions.ErrActiveExperiment) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "active experiment with this flag already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create experiment"})
		return
	}

	c.JSON(http.StatusCreated, experiment)
}

func (h *ExperimentsHandler) ChangeStatus(c *gin.Context) {
	var req requests.ChangeStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json payload"})
		return
	}

	experiment, err := h.experimentsService.ChangeStatus(c.Request.Context(), inputs.ChangeStatusInput{
		Identifier: req.Identifier,
		Status:     req.Status,
	})

	if err != nil {
		if errors.Is(err, exceptions.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid inputs"})
			return
		}

		if errors.Is(err, exceptions.ErrExperimentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "experiment not found"})
			return
		}

		if errors.Is(err, exceptions.ErrActiveExperiment) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "active experiment with this flag already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update experiment's status"})
		return
	}

	c.JSON(http.StatusOK, experiment)
}
