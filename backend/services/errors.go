// services/errors.go
//
// Paquete services: definición de errores de dominio y utilidades
// para envolver causas concretas bajo sentinelas estables.
//
// Convenciones:
//   - Los servicios nunca retornan errores "crudos" del driver o librerías externas.
//   - Siempre se debe envolver con Wrap(...) usando un sentinel de dominio.
//   - Los controladores hacen match con errors.Is(err, services.ErrX)
//     para traducir a HTTP, sin depender de strings ni detalles de implementación.
package services

import (
	"errors"
	"fmt"
)

// Sentinelas de error de dominio.
// Son valores comparables con errors.Is(...) y se mantienen estables
// aunque la implementación cambie.
var (
	// ErrInvalidInput: datos de entrada inválidos (DTO, validaciones de negocio).
	ErrInvalidInput = errors.New("invalid_input")

	// ErrInvalidID: identificador con formato incorrecto (ej. ObjectID mal formado).
	ErrInvalidID = errors.New("invalid_id")

	// ErrNotFound: recurso no encontrado en el almacenamiento.
	ErrNotFound = errors.New("not_found")

	// ErrConflict: conflicto de estado, duplicados o violación de restricciones únicas.
	ErrConflict = errors.New("conflict")

	// ErrDB: error del driver o infraestructura de datos (Mongo, SQL, etc.).
	ErrDB = errors.New("db_error")
)

// Wrap envuelve un error concreto (cause) con un sentinel de dominio y un mensaje opcional.
//
// Parámetros:
//   - cause: error original (puede ser nil).
//   - sentinel: uno de los Err* definidos arriba (ej. ErrDB, ErrNotFound).
//   - msg: contexto opcional de la operación fallida ("insert post", "parse objectid").
//
// Retorna:
//   - nil si cause es nil.
//   - error compuesto con sentinel y la causa.
//     Puede compararse con errors.Is(err, sentinel) en capas superiores.
//
// Ejemplos:
//   return services.Wrap(err, services.ErrDB, "insert post")
//   return services.Wrap(err, services.ErrInvalidID, "")
func Wrap(cause error, sentinel error, msg string) error {
	if cause == nil {
		return nil
	}
	if msg == "" {
		return fmt.Errorf("%w: %v", sentinel, cause)
	}
	return fmt.Errorf("%w: %s: %v", sentinel, msg, cause)
}
