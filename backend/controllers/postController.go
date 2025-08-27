// controllers/postController.go
//
// Paquete controllers: capa HTTP para la entidad Post.
// Convenciones:
//   - Las validaciones de entrada (DTO / query params) viven aquí.
//   - La traducción de errores a HTTP se realiza con writeError(...) usando sentinelas de services.
//   - La lógica de negocio y acceso a datos está encapsulada en services.
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

// httpError define la forma uniforme de las respuestas de error HTTP.
type httpError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// writeError traduce los errores de dominio (sentinelas de services) a códigos HTTP.
// Los controladores deben delegar aquí cualquier error retornado por services.
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

// CreatePost maneja POST /api/posts.
// - Valida el DTO de entrada.
// - Delegar en services.CreatePost.
// - Responde 201 con Location y el insertedID.
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

// GetPostByID maneja GET /api/posts/:id.
// - Valida presencia de :id.
// - Delegar en services.GetPostByID.
// - Responde 200 con el documento, 400/404/500 según corresponda.
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

// UpdatePostByID maneja PUT /api/posts/:id.
// - Valida :id y DTO.
// - Delegar en services.UpdatePostByID (que fija PublishedAt si aplica).
// - Responde 200 con el documento actualizado.
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

// DeletePostByID maneja DELETE /api/posts/:id.
// - Valida :id.
// - Delegar en services.DeletePostByID.
// - Responde 204 si elimina; 400/404/500 si falla.
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
	c.Status(http.StatusNoContent)
}

// GetPostsMetricsByTag maneja GET /api/posts/metrics/by-tag?limit=&onlyPublished=.
//
// Query params:
//   - limit (opcional, entero > 0; default 10; máx 100).
//   - onlyPublished (opcional, "true"/"false") para filtrar por estado.
//
// Respuestas: 200 con []TagMetric; 400 si parámetros inválidos; 500 si falla la agregación.
func GetPostsMetricsByTag(c *gin.Context) {
	// limit
	limit := 10
	if raw := c.Query("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			writeError(c, services.Wrap(err, services.ErrInvalidInput, "limit must be a positive integer"))
			return
		}
		limit = n
	}

	// onlyPublished
	var onlyPublished *bool
	if raw := c.Query("onlyPublished"); raw != "" {
		switch strings.ToLower(raw) {
		case "true":
			v := true
			onlyPublished = &v
		case "false":
			v := false
			onlyPublished = &v
		default:
			writeError(c, services.Wrap(nil, services.ErrInvalidInput, "onlyPublished must be true or false"))
			return
		}
	}

	metrics, err := services.GetPostsMetricsByTag(c.Request.Context(), limit, onlyPublished)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, metrics)
}

// ListPosts maneja:
//   GET /api/posts?q=&tag=&published=(true|false)&page=&limit=&sort=publishedAt|-publishedAt
//
// Query params:
//   - q: búsqueda de texto (requiere índice {title:"text", content:"text"}).
//   - tag: filtra por etiqueta exacta.
//   - published: "true" | "false" | "" (sin filtro).
//   - page, limit: enteros positivos (limit se sujeta a tope en services).
//   - sort: "publishedAt" | "-publishedAt" (default: "-publishedAt").
//
// Respuestas: 200 con ListPostsResult; 400 si parámetros inválidos; 500 si falla el driver.
func ListPosts(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	tag := strings.TrimSpace(c.Query("tag"))

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

	sort := c.Query("sort")

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
