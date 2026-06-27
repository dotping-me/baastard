package handlers

import (
	"net/http"

	"github.com/dotping-me/baastard/payload"
	"github.com/dotping-me/baastard/services"
	"github.com/gin-gonic/gin"
)

type MetadataHandler struct { // Also handles FieldMetadata
	metadata *services.MetadataService
}

func NewMetadataHandler(s *services.MetadataService) *MetadataHandler {
	return &MetadataHandler{metadata: s}
}

// Handler that creates an Entity Metdata and respective Fields Metadata
func (h *MetadataHandler) Create(ctx *gin.Context) {

	var req payload.CreateEntityMetadataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	em, err := h.metadata.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, em)
}

// Handler that returns all Entities' metadata
func (h *MetadataHandler) GetAll(ctx *gin.Context) {
	collections, err := h.metadata.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, collections)
}

// Handler that returns a single Entity Metadata bearing a given id parsed through
// URL Param @ GET /api/metadata/:id
func (h *MetadataHandler) GetById(ctx *gin.Context) {
	collection, err := h.metadata.GetById(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Collection not found!",
		})

		return
	}

	ctx.JSON(http.StatusOK, collection)
}

// Handler that deletes a collection with a given id parsed through URL Param
// @ DELETE /api/metadata/:id
func (h *MetadataHandler) Delete(ctx *gin.Context) {
	err := h.metadata.Delete(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Collection deleted!",
	})
}
