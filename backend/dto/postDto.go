// dto/postDto.go
//
// Paquete dto: objetos de transferencia de datos (Data Transfer Objects).
// Se utilizan en la capa de controladores para validar y recibir payloads JSON
// de entrada en los endpoints de la API.
//
// Convenciones:
//   - Todos los campos llevan tags `json` para mapearse al body.
//   - Se usan tags `binding` de Gin/validator para validación automática.
//   - No contienen lógica de negocio ni persistencia.
package dto

// CreatePostDTO define el cuerpo esperado en POST /api/posts.
//
// Validaciones:
//   - Title: requerido, entre 5 y 140 caracteres.
//   - Author: requerido.
//   - Content: requerido.
//   - Tags: opcional, arreglo de strings.
//   - Published: opcional (default false si no se envía).
//
// Ejemplo JSON:
//   {
//     "title": "Introducción a Go",
//     "author": "Alice",
//     "content": "Texto del post...",
//     "tags": ["go", "backend"],
//     "published": true
//   }
type CreatePostDTO struct {
	Title     string   `json:"title"   binding:"required,min=5,max=140"`
	Author    string   `json:"author"  binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}

// UpdatePostDTO define el cuerpo esperado en PUT /api/posts/:id.
//
// Validaciones:
//   - Title: requerido, entre 5 y 140 caracteres.
//   - Author: requerido.
//   - Content: requerido.
//   - Tags: opcional, arreglo de strings.
//   - Published: opcional.
//
// Nota: PublishedAt no se controla aquí; lo fija la capa de servicio
//       cuando se cambia Published de false→true.
//
// Ejemplo JSON:
//   {
//     "title": "Título actualizado",
//     "author": "Alice",
//     "content": "Contenido nuevo",
//     "tags": ["go", "mongo"],
//     "published": false
//   }
type UpdatePostDTO struct {
	Title     string   `json:"title"   binding:"required,min=5,max=140"`
	Author    string   `json:"author"  binding:"required"`
	Content   string   `json:"content" binding:"required"`
	Tags      []string `json:"tags"`
	Published bool     `json:"published"`
}
