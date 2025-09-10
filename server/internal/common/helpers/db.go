package helpers

import "database/sql"

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

func Int64ToNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{Int64: n, Valid: true}
}

func OutIndexToNullInt64(n int64) sql.NullInt64 {
	if n < 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: n, Valid: true}
}
