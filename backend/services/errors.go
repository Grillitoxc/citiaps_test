// services/errors.go
package services

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid_input") // body inválido, tipos, constraints de dominio
	ErrInvalidID    = errors.New("invalid_id")    // parse de identificadores (ObjectID, UUID, etc.)
	ErrNotFound     = errors.New("not_found")     // recurso no existe
	ErrConflict     = errors.New("conflict")      // estado inconsistente / duplicados / violaciones únicas
	ErrDB           = errors.New("db_error")      // errores del driver/infra de datos
)

func Wrap(cause error, sentinel error, msg string) error {
	if cause == nil {
		return nil
	}
	if msg == "" {
		return fmt.Errorf("%w: %v", sentinel, cause)
	}
	return fmt.Errorf("%w: %s: %v", sentinel, msg, cause)
}
