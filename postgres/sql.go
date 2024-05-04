package postgres

import "fmt"

func FormatLimitOffset(limit, offset int) string {
	if offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	}
	return fmt.Sprintf(`LIMIT %d`, limit)
}

func FormatOrderBy(column string, direction string) string {
	if column != "" {
		return fmt.Sprintf(`ORDER BY %s %s`, column, direction)
	}
	return ""
}
