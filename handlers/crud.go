package handlers

import (
	"net/http"

	"github.com/dotping-me/baastard/services"
	"github.com/gin-gonic/gin"
)

type CrudHandler struct {
	crud *services.CrudService
}

func NewCrudHandler(crud *services.CrudService) *CrudHandler {
	return &CrudHandler{crud: crud}
}

// Returns all rows from a table bearing name parsed through URL Param
// @ GET /api/data/:table
func (h *CrudHandler) GetRows(ctx *gin.Context) {
	table := ctx.Param("table")

	rows, err := h.crud.GetRows(table)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, rows)
}

// Returns a specific row by id from a table bearing name parsed through URL Param
// @ GET /api/data/:table/:id
func (h *CrudHandler) GetRow(ctx *gin.Context) {
	table := ctx.Param("table")
	id := ctx.Param("id")

	row, err := h.crud.GetRow(table, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, row)
}

// Inserting data from a JSON payload into a table. Table is parsed through URL
// Param @ POST /api/data/:table
func (h *CrudHandler) Insert(ctx *gin.Context) {
	table := ctx.Param("table")

	var body map[string]any
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = h.crud.Insert(table, body) // Tries to insert data into table
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Row inserted",
	})
}

// Updating data by id in a given table: table name parsed through URL Param
// @ PATCH /api/data/:table/:id
func (h *CrudHandler) Update(ctx *gin.Context) {
	table := ctx.Param("table")
	id := ctx.Param("id")

	var body map[string]any
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = h.crud.Update(table, id, body) // Tries to update data in table
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Row updated",
	})
}

// Deleting data by id from a table @ DELETE /api/data/:table/:id
func (h *CrudHandler) Delete(ctx *gin.Context) {
	table := ctx.Param("table")
	id := ctx.Param("id")

	err := h.crud.Delete(table, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Row deleted",
	})
}
