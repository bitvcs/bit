package sqlite

import (
	"database/sql"
	"errors"
	"time"

	"github.com/nipalab/nipa/internal/domain"
)

func handleError(err error) error {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.NewErrorRecordNotFound()
		}
		return domain.NewErrorDatabase(err.Error())
	}
	return nil
}

func nullTimePtr(t sql.NullTime) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
