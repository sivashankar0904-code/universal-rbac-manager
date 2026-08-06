package repo

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"urm/internal/apperr"
)

func translateErr(err error) error {
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
