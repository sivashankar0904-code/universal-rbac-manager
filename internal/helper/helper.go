package helper

import (
	"context"
	"errors"
	"urm/internal/apperr"

	"github.com/fatih/structs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// StructToMapStringWithTag converts a struct to map[string]interface, with the keys set as the
// tag value set
func (s *Service) StructToMapStringWithTag(source interface{}) (map[string]interface{}, error) {
	converter := structs.New(source)
	converter.TagName = "json"
	return converter.Map(), nil
}

// validateRequiredFields checks that every value in fields is non-empty,
// keyed by field name for the error response (e.g. {"key": o.Key}).
func (s *Service) ValidateRequiredFields(object map[string]interface{}, validation_fields []string) error {
	missing := map[string]string{}
	for _, name := range validation_fields {
		value, ok := object[name]
		if !ok || value == "" || value == nil {
			missing[name] = "required"
		}
	}
	if len(missing) > 0 {
		return apperr.Validation("invalid organisation", missing)
	}
	return nil
}

func (s *Service) TranslateErr(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.Wrap(err, apperr.KindNotFound, "resource not found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return apperr.Wrap(err, apperr.KindConflict, "resource already exists")
		case "23503":
			return apperr.Wrap(err, apperr.KindValidation, "referenced resource does not exist")
		}
	}

	return apperr.Wrap(err, apperr.KindInternal, "database error")
}
