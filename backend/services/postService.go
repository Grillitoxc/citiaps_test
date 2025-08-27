// services/postService.go
//
// Paquete services: lógica de negocio y acceso a datos para la entidad Post.
// Convenciones:
//   - Todas las funciones que interactúan con MongoDB aplican un timeout fijo (defaultTimeout).
//   - Todos los errores del driver/infra se envuelven con sentinelas (ErrDB, ErrNotFound, ErrInvalidID, etc.).
//   - No se exponen errores del driver a capas superiores; use errors.Is(err, services.ErrX) en controladores.
package services

import (
	"context"
	"errors"
	"time"

	"blog-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// defaultTimeout define el tiempo máximo de espera por operación a MongoDB.
	defaultTimeout = 5 * time.Second
	// maxPageLimit limita la cantidad de documentos por página para listados.
	maxPageLimit = 100
)

// CreatePost inserta un nuevo Post.
//
// Reglas:
//   - Estampa CreatedAt=now.
//   - Si p.Published==true y p.PublishedAt==nil, fija PublishedAt=now.
//
// Parámetros:
//   - ctx: contexto de cancelación/timeout.
//   - p: entidad a persistir (campos de cliente validados en capas superiores).
//
// Retornos:
//   - ObjectID del documento insertado.
//   - error con sentinelas: ErrDB si el driver falla.
//
// Errores:
//   - ErrDB: error del driver o de infraestructura.
//   - (No valida campos de dominio; esas validaciones están en DTO/controlador).
func CreatePost(ctx context.Context, p models.Post) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	now := time.Now().UTC()
	p.CreatedAt = now
	if p.Published && p.PublishedAt == nil {
		p.PublishedAt = &now
	}

	res, err := DB.Collection("posts").InsertOne(ctx, p)
	if err != nil {
		return primitive.NilObjectID, Wrap(err, ErrDB, "insert post")
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, Wrap(errors.New("inserted id not an ObjectID"), ErrDB, "cast inserted id")
	}
	return oid, nil
}

// GetPostByID recupera un Post por su ObjectID (hexadecimal).
//
// Parámetros:
//   - ctx: contexto de cancelación/timeout.
//   - idHex: cadena hexadecimal de 24 caracteres correspondiente al ObjectID.
//
// Retornos:
//   - Post encontrado.
//   - error con sentinelas: ErrInvalidID si el id es inválido; ErrNotFound si no existe; ErrDB si falla el driver.
func GetPostByID(ctx context.Context, idHex string) (models.Post, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return models.Post{}, Wrap(err, ErrInvalidID, "parse objectid")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var out models.Post
	if err := DB.Collection("posts").FindOne(ctx, bson.M{"_id": oid}).Decode(&out); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Post{}, Wrap(err, ErrNotFound, "post not found")
		}
		return models.Post{}, Wrap(err, ErrDB, "find post")
	}
	return out, nil
}

// UpdatePostByID actualiza campos de un Post y retorna el documento resultante.
//
// Reglas:
//   - Si el Post actual no está publicado y el nuevo estado Published pasa a true,
//     y PublishedAt es nil, se fija PublishedAt=now.
//   - Actualiza: title, author, content, tags, published; updatedAt=now.
//   - publishedAt sólo se actualiza si viene definido o si aplica la regla anterior.
//
// Parámetros:
//   - ctx: contexto de cancelación/timeout.
//   - idHex: ObjectID en hex del documento a actualizar.
//   - p: valores a aplicar.
//
// Retornos:
//   - Post actualizado (estado After).
//   - error con sentinelas: ErrInvalidID, ErrNotFound, ErrDB.
func UpdatePostByID(ctx context.Context, idHex string, p models.Post) (models.Post, error) {
	current, err := GetPostByID(ctx, idHex)
	if err != nil {
		return models.Post{}, err
	}

	if !current.Published && p.Published && p.PublishedAt == nil {
		now := time.Now().UTC()
		p.PublishedAt = &now
	}

	set := bson.M{
		"title":     p.Title,
		"author":    p.Author,
		"content":   p.Content,
		"tags":      p.Tags,
		"published": p.Published,
		"updatedAt": time.Now().UTC(),
	}
	if p.PublishedAt != nil {
		set["publishedAt"] = p.PublishedAt
	}

	oid, _ := primitive.ObjectIDFromHex(idHex) // validado por GetPostByID
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Post
	if err := DB.Collection("posts").
		FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, opts).
		Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Post{}, Wrap(err, ErrNotFound, "post not found after update")
		}
		return models.Post{}, Wrap(err, ErrDB, "findOneAndUpdate post")
	}
	return updated, nil
}

// DeletePostByID elimina un Post por id.
//
// Parámetros:
//   - ctx: contexto de cancelación/timeout.
//   - idHex: ObjectID en hex.
//
// Retornos:
//   - nil si elimina correctamente.
//   - error con sentinelas: ErrInvalidID si el id es inválido;
//     ErrNotFound si no existía un documento con ese id; ErrDB en fallas del driver.
func DeletePostByID(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Wrap(err, ErrInvalidID, "parse objectid")
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	res, err := DB.Collection("posts").DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return Wrap(err, ErrDB, "delete post")
	}
	if res.DeletedCount == 0 {
		return Wrap(mongo.ErrNoDocuments, ErrNotFound, "post not found")
	}
	return nil
}

// TagMetric representa la métrica de cantidad de posts por etiqueta.
type TagMetric struct {
	Tag   string `bson:"_id"  json:"tag"`
	Count int64  `bson:"count" json:"count"`
}

// GetPostsMetricsByTag devuelve el top-N de etiquetas por cantidad de posts.
//
// Parámetros:
//   - ctx: contexto de cancelación/timeout.
//   - limit: máximo de filas a retornar (si <=0 se usa 10; si >100 se trunca a 100).
//   - onlyPublished: si no es nil, filtra por published==*onlyPublished.
//
// Retornos:
//   - slice ordenado descendentemente por Count.
//   - error con sentinelas: ErrDB ante errores del pipeline/cursor.
func GetPostsMetricsByTag(ctx context.Context, limit int, onlyPublished *bool) ([]TagMetric, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	match := bson.M{"tags": bson.M{"$type": "string"}}
	if onlyPublished != nil {
		match["published"] = *onlyPublished
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: "$tags"}},
		{{Key: "$match", Value: bson.M{"tags": bson.M{"$ne": ""}}}},
		{{Key: "$group", Value: bson.M{"_id": "$tags", "count": bson.M{"$sum": 1}}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$limit", Value: limit}},
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	cur, err := DB.Collection("posts").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, Wrap(err, ErrDB, "aggregate by-tag")
	}
	defer cur.Close(ctx)

	out := make([]TagMetric, 0, limit)
	for cur.Next(ctx) {
		var m TagMetric
		if err := cur.Decode(&m); err != nil {
			return nil, Wrap(err, ErrDB, "decode metric row")
		}
		out = append(out, m)
	}
	if err := cur.Err(); err != nil {
		return nil, Wrap(err, ErrDB, "cursor error")
	}
	return out, nil
}

// ListPostsParams define filtros de búsqueda/orden paginada.
type ListPostsParams struct {
	// Q aplica búsqueda de texto (requiere índice en {title: "text", content: "text"}).
	Q string
	// Tag filtra posts que contengan exactamente esa etiqueta.
	Tag string
	// Published: nil = no filtra; true/false = filtra por estado de publicación.
	Published *bool
	// Page: número de página 1-based.
	Page int
	// Limit: tamaño de página (se trunca a maxPageLimit).
	Limit int
	// SortField: "publishedAt" | "-publishedAt". Default: "-publishedAt".
	SortField string
}

// ListPostsResult contiene los ítems y metadatos de paginación.
type ListPostsResult struct {
	Items      []models.Post `json:"items"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int64         `json:"totalPages"`
}

// ListPosts lista posts con búsqueda, filtros, orden y paginación.
//
// Parámetros:
//   - ctx: contexto de cancelación/timeout.
//   - p: parámetros de filtro y paginación (ver ListPostsParams).
//
// Retornos:
//   - Listado paginado + total y totalPages.
//   - error con sentinelas: ErrDB ante fallas del driver.
//
// Notas:
//   - Índices recomendados: {title:"text", content:"text"} para Q; {published:1, publishedAt:-1} para listados.
func ListPosts(ctx context.Context, p ListPostsParams) (ListPostsResult, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > maxPageLimit {
		p.Limit = maxPageLimit
	}

	filter := bson.M{}
	if p.Q != "" {
		filter["$text"] = bson.M{"$search": p.Q}
	}
	if p.Tag != "" {
		filter["tags"] = p.Tag
	}
	if p.Published != nil {
		filter["published"] = *p.Published
	}

	var sort bson.D
	switch p.SortField {
	case "publishedAt":
		sort = bson.D{{Key: "publishedAt", Value: 1}}
	case "-publishedAt":
		sort = bson.D{{Key: "publishedAt", Value: -1}}
	default:
		sort = bson.D{{Key: "publishedAt", Value: -1}}
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	total, err := DB.Collection("posts").CountDocuments(ctx, filter)
	if err != nil {
		return ListPostsResult{}, Wrap(err, ErrDB, "count posts")
	}

	opts := options.Find().
		SetSort(sort).
		SetSkip(int64((p.Page-1)*p.Limit)).
		SetLimit(int64(p.Limit))

	cur, err := DB.Collection("posts").Find(ctx, filter, opts)
	if err != nil {
		return ListPostsResult{}, Wrap(err, ErrDB, "find posts")
	}
	defer cur.Close(ctx)

	items := make([]models.Post, 0, p.Limit)
	for cur.Next(ctx) {
		var post models.Post
		if err := cur.Decode(&post); err != nil {
			return ListPostsResult{}, Wrap(err, ErrDB, "decode post")
		}
		items = append(items, post)
	}
	if err := cur.Err(); err != nil {
		return ListPostsResult{}, Wrap(err, ErrDB, "cursor error")
	}

	totalPages := (total + int64(p.Limit) - 1) / int64(p.Limit)
	return ListPostsResult{
		Items:      items,
		Total:      total,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalPages: totalPages,
	}, nil
}
