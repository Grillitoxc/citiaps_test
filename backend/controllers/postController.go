// controllers/postController.go
package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"blog-api/dto"
	"blog-api/models"
	"blog-api/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Estructura uniforme para respuestas de error HTTP.
type httpError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// writeError traduce sentinelas de services → códigos HTTP.
func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidInput),
		errors.Is(err, services.ErrInvalidID):
		c.JSON(http.StatusBadRequest, httpError{
			Code:    http.StatusBadRequest,
			Message: "Solicitud inválida",
			Details: err.Error(),
		})
	case errors.Is(err, services.ErrNotFound):
		c.JSON(http.StatusNotFound, httpError{
			Code:    http.StatusNotFound,
			Message: "Recurso no encontrado",
		})
	case errors.Is(err, services.ErrConflict):
		c.JSON(http.StatusConflict, httpError{
			Code:    http.StatusConflict,
			Message: "Conflicto de estado",
		})
	case errors.Is(err, services.ErrDB):
		c.JSON(http.StatusInternalServerError, httpError{
			Code:    http.StatusInternalServerError,
			Message: "Error interno",
		})
	default:
		c.JSON(http.StatusInternalServerError, httpError{
			Code:    http.StatusInternalServerError,
			Message: "Error interno",
		})
	}
}

// CreatePost valida el DTO, delega al servicio y responde consistente.
func CreatePost(c *gin.Context) {
	var in dto.CreatePostDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			fields := map[string]string{}
			for _, fe := range verrs {
				fields[fe.Field()] = fe.Tag()
			}
			writeError(c, services.Wrap(verrs, services.ErrInvalidInput, "validation"))
			return
		}
		writeError(c, services.Wrap(err, services.ErrInvalidInput, "bind json"))
		return
	}

	post := models.Post{
		Title:     in.Title,
		Author:    in.Author,
		Content:   in.Content,
		Tags:      in.Tags,
		Published: in.Published,
	}

	id, err := services.CreatePost(c.Request.Context(), post)
	if err != nil {
		writeError(c, err)
		return
	}

	c.Header("Location", fmt.Sprintf("/api/posts/%s", id.Hex()))
	c.JSON(http.StatusCreated, gin.H{"insertedID": id.Hex()})
}

// GetPostByID recibe el id, delega al servicio y usa el mapper de errores.
func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		writeError(c, services.ErrInvalidID)
		return
	}

	post, err := services.GetPostByID(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, post)
}

func UpdatePostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		writeError(c, services.ErrInvalidID)
		return
	}

	var in dto.UpdatePostDTO
	if err := c.ShouldBindJSON(&in); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			fields := map[string]string{}
			for _, fe := range verrs {
				fields[fe.Field()] = fe.Tag()
			}
			writeError(c, services.Wrap(verrs, services.ErrInvalidInput, "validation"))
			return
		}
		writeError(c, services.Wrap(err, services.ErrInvalidInput, "bind json"))
		return
	}

	post := models.Post{
		Title:     in.Title,
		Author:    in.Author,
		Content:   in.Content,
		Tags:      in.Tags,
		Published: in.Published,
	}

	updated, err := services.UpdatePostByID(c.Request.Context(), id, post)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}


func DeletePostByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		writeError(c, services.ErrInvalidID)
		return
	}

	if err := services.DeletePostByID(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}

	// Convención REST: 204 No Content al eliminar correctamente
	c.Status(http.StatusNoContent)
}


// GetPostsMetricsByTag maneja: GET /api/posts/metrics/by-tag
func GetPostsMetricsByTag(c *gin.Context) {
	metrics, err := services.GetPostsMetricsByTag(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, metrics)
}

func ListPosts(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	tag := strings.TrimSpace(c.Query("tag"))

	// published: "true" | "false" | "" (no filtrar)
	var publishedPtr *bool
	if raw := c.Query("published"); raw != "" {
		switch strings.ToLower(raw) {
		case "true":
			v := true
			publishedPtr = &v
		case "false":
			v := false
			publishedPtr = &v
		default:
			writeError(c, services.Wrap(nil, services.ErrInvalidInput, "published must be true or false"))
			return
		}
	}

	// page y limit
	page := 1
	if raw := c.Query("page"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			writeError(c, services.Wrap(err, services.ErrInvalidInput, "page must be a positive integer"))
			return
		}
		page = n
	}
	limit := 10
	if raw := c.Query("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			writeError(c, services.Wrap(err, services.ErrInvalidInput, "limit must be a positive integer"))
			return
		}
		limit = n
	}

	sort := c.Query("sort") // "publishedAt" | "-publishedAt" | ""

	params := services.ListPostsParams{
		Q:         q,
		Tag:       tag,
		Published: publishedPtr,
		Page:      page,
		Limit:     limit,
		SortField: sort,
	}

	result, err := services.ListPosts(c.Request.Context(), params)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}