package postgres

import "fmt"

func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}

func FormatOrderBy(column string) string {
	if column != "" {
		return fmt.Sprintf(`ORDER BY %s`, column)
	}
	return ""
}
