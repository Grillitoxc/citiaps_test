package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Title       string             `bson:"title" json:"title" binding:"required,min=5,max=140"`
	Author      string             `bson:"author" json:"author" binding:"required"`
	Content     string             `bson:"content" json:"content" binding:"required"`
	Tags        []string           `bson:"tags,omitempty" json:"tags,omitempty"`
	Published   bool               `bson:"published" json:"published"`
	PublishedAt *time.Time         `bson:"publishedAt,omitempty" json:"publishedAt,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}
