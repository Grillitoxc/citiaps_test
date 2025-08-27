package dto

type CreatePostDTO struct {
	Title     string   `json:"title"   binding:"required,min=5,max=140"`
	Author    string   `json:"author"  binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

type UpdatePostDTO struct {
	Title     string   `json:"title"   binding:"required,min=5,max=140"`
	Author    string   `json:"author"  binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}