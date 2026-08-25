package sqlite

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/nipalab/nipa/internal/domain"
)

func TestHandleError_Nil(t *testing.T) {
	require.NoError(t, handleError(nil))
}

func TestHandleError_RecordNotFound(t *testing.T) {
	err := handleError(sql.ErrNoRows)
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 404, domErr.Code)
	require.Equal(t, "record not found", domErr.Message)
}

func TestHandleError_DatabaseError(t *testing.T) {
	err := handleError(errors.New("connection refused"))
	require.Error(t, err)

	var domErr *domain.Error
	require.ErrorAs(t, err, &domErr)
	require.Equal(t, 500, domErr.Code)
	require.Equal(t, "database error", domErr.Message)
	require.Equal(t, "connection refused", domErr.InternalMessage)
}

func TestNullTimePtr_Valid(t *testing.T) {
	ts := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	got := nullTimePtr(sql.NullTime{Time: ts, Valid: true})
	require.NotNil(t, got)
	require.Equal(t, ts, *got)
}

func TestNullTimePtr_Invalid(t *testing.T) {
	got := nullTimePtr(sql.NullTime{})
	require.Nil(t, got)
}

func TestNullInt64Ptr_Valid(t *testing.T) {
	got := nullInt64Ptr(sql.NullInt64{Int64: 42, Valid: true})
	require.NotNil(t, got)
	require.Equal(t, int64(42), *got)
}

func TestNullInt64Ptr_Invalid(t *testing.T) {
	got := nullInt64Ptr(sql.NullInt64{})
	require.Nil(t, got)
}

func TestTimePtrToNullTime_NonNil(t *testing.T) {
	ts := time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC)
	got := timePtrToNullTime(&ts)
	require.True(t, got.Valid)
	require.Equal(t, ts, got.Time)
}

func TestTimePtrToNullTime_Nil(t *testing.T) {
	got := timePtrToNullTime(nil)
	require.False(t, got.Valid)
}
