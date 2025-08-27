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

func CreatePost(ctx context.Context, p models.Post) (primitive.ObjectID, error) {
	now := time.Now().UTC()
	p.CreatedAt = now
	if p.Published {
		p.PublishedAt = &now
	}

	res, err := DB.Collection("posts").InsertOne(ctx, p)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid, nil
}

func GetPostByID(ctx context.Context, idHex string) (models.Post, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return models.Post{}, Wrap(err, ErrInvalidID, "parse objectid")
	}

	var out models.Post
	err = DB.Collection("posts").FindOne(ctx, bson.M{"_id": oid}).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Post{}, Wrap(err, ErrNotFound, "post not found")
		}
		return models.Post{}, Wrap(err, ErrDB, "find post")
	}
	return out, nil
}

func UpdatePostByID(ctx context.Context, idHex string, p models.Post) (models.Post, error) {
	// Reutilizamos GetPostByID para validar id y obtener el documento actual
	current, err := GetPostByID(ctx, idHex)
	if err != nil {
		return models.Post{}, err
	}

	// Si published cambia de false -> true, set publishedAt
	if !current.Published && p.Published {
		now := time.Now().UTC()
		p.PublishedAt = &now
	}

	// Update document
	update := bson.M{
		"$set": bson.M{
			"title":       p.Title,
			"author":      p.Author,
			"content":     p.Content,
			"tags":        p.Tags,
			"published":   p.Published,
			"publishedAt": p.PublishedAt,
		},
	}

	oid, _ := primitive.ObjectIDFromHex(idHex) // seguro: ya lo validó GetPostByID
	_, err = DB.Collection("posts").UpdateByID(ctx, oid, update)
	if err != nil {
		return models.Post{}, Wrap(err, ErrDB, "update post")
	}

	// Devolver post actualizado
	updated, err := GetPostByID(ctx, idHex)
	if err != nil {
		return models.Post{}, err
	}
	return updated, nil
}


func DeletePostByID(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return Wrap(err, ErrInvalidID, "parse objectid")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := DB.Collection("posts").DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return Wrap(err, ErrDB, "delete post")
	}
	if res.DeletedCount == 0 {
		// No había documento con ese _id
		return Wrap(mongo.ErrNoDocuments, ErrNotFound, "post not found")
	}
	return nil
}


type TagMetric struct {
	Tag   string `bson:"_id"  json:"tag"`
	Count int64  `bson:"count" json:"count"`
}

// GetPostsMetricsByTag ejecuta la agregación y retorna las top-N etiquetas por cantidad de posts.
func GetPostsMetricsByTag(ctx context.Context) ([]TagMetric, error) {
	limit := 10

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"tags": bson.M{"$type": "string"}}}},
		{{Key: "$unwind", Value: "$tags"}},
		{{Key: "$match", Value: bson.M{"tags": bson.M{"$ne": ""}}}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$tags",
			"count": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
		{{Key: "$limit", Value: limit}},
	}

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



type ListPostsParams struct {
	Q         string  // búsqueda de texto
	Tag       string  // una etiqueta exacta
	Published *bool   // nil = sin filtro; true/false = filtra
	Page      int     // 1-based
	Limit     int     // tamaño de página
	SortField string  // "publishedAt" o "-publishedAt"
}

// ListPosts devuelve posts paginados + total de coincidencias.
type ListPostsResult struct {
	Items      []models.Post `json:"items"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int64         `json:"totalPages"`
}

func ListPosts(ctx context.Context, p ListPostsParams) (ListPostsResult, error) {
	// Sanitización de paginación
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}

	// Construcción del filtro
	filter := bson.M{}
	if p.Q != "" {
		// Requiere índice de texto en title y content
		filter["$text"] = bson.M{"$search": p.Q}
	}
	if p.Tag != "" {
		filter["tags"] = p.Tag
	}
	if p.Published != nil {
		filter["published"] = *p.Published
	}

	// Sort
	sort := bson.D{}
	switch p.SortField {
	case "publishedAt":
		sort = bson.D{{Key: "publishedAt", Value: 1}}
	case "-publishedAt":
		sort = bson.D{{Key: "publishedAt", Value: -1}}
	default:
		// Por defecto: publicados más recientes primero
		sort = bson.D{{Key: "publishedAt", Value: -1}}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Total para paginación
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