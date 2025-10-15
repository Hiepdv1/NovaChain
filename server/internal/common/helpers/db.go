package helpers

import (
	"database/sql"
	"fmt"
)

func StringPtrToNullString(s *string) sql.NullString {
	if s == nil || *s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func FloatToNullString(n float64) sql.NullString {
	return sql.NullString{
		String: fmt.Sprintf("%.8f", n),
		Valid:  true,
	}
}
func Int64ToNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}

func Int32ToNullInt32(n int32) sql.NullInt32 {
	return sql.NullInt32{Int32: n, Valid: true}
}

func OutIndexToNullInt64(n int64) sql.NullInt64 {
	if n < 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: n, Valid: true}
}
