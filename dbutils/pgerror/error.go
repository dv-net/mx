package pgerror

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUniqueViolation = errors.New("unique constraint violation")
	ErrNotFound        = errors.New("record not found")
)

func ParseError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return &UniqueConstraintError{
				Err:        ErrUniqueViolation,
				Constraint: pgErr.ConstraintName,
				Table:      pgErr.TableName,
				Detail:     pgErr.Detail,
			}
		case "02000":
			return &NotFoundError{
				Err:    ErrNotFound,
				Table:  pgErr.TableName,
				Detail: pgErr.Detail,
			}
		}
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return &NotFoundError{
			Err:    ErrNotFound,
			Detail: "no rows found in query",
		}
	}

	return err
}

type UniqueConstraintError struct {
	Err        error
	Constraint string
	Table      string
	Detail     string
}

func (e *UniqueConstraintError) Error() string {
	return e.Err.Error() + ": " + e.Detail
}

func (e *UniqueConstraintError) Unwrap() error {
	return e.Err
}

type NotFoundError struct {
	Err    error
	Table  string
	Detail string
}

func (e *NotFoundError) Error() string {
	return e.Err.Error() + ": " + e.Detail
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}
