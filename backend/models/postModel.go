// models/post.go
//
// Paquete models: definición de entidades de dominio persistidas en MongoDB.
// 
// Convenciones:
//   - Los structs llevan tags `bson` para mapearse a MongoDB y `json` para la salida HTTP.
//   - Validaciones básicas (binding) se pueden usar en controladores/DTO cuando corresponda.
//   - Los timestamps CreatedAt y PublishedAt son gestionados por la capa de servicios.
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post representa una publicación en el blog.
//
// Campos:
//   - ID: identificador único (ObjectID de MongoDB).
//   - Title: título del post, requerido, 5–140 caracteres.
//   - Author: autor del post, requerido.
//   - Content: contenido del post, requerido.
//   - Tags: etiquetas asociadas (opcional).
//   - Published: indica si el post está publicado.
//   - PublishedAt: fecha/hora en UTC en que se publicó (nil si no publicado).
//   - CreatedAt: fecha/hora en UTC en que se creó.
//
// Serialización:
//   - bson: usado por el driver de MongoDB.
//   - json: usado en las respuestas HTTP.
//
// Validaciones:
//   - binding:"required,min=5,max=140" en Title.
//   - binding:"required" en Author y Content.
//
// Notas:
//   - CreatedAt y PublishedAt son controlados por la capa de servicios, no por el cliente.
//   - PublishedAt se fija automáticamente cuando Published cambia de false→true.
type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"    json:"_id"`
	Title       string             `bson:"title"            json:"title"   binding:"required,min=5,max=140"`
	Author      string             `bson:"author"           json:"author"  binding:"required"`
	Content     string             `bson:"content"          json:"content" binding:"required"`
	Tags        []string           `bson:"tags,omitempty"   json:"tags,omitempty"`
	Published   bool               `bson:"published"        json:"published"`
	PublishedAt *time.Time         `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt"        json:"createdAt"`
}
