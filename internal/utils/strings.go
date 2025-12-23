package utils

import "strings"

func ToSnakeCase(str string) string {
	var b strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			// Check if previous char was lowercase or number
			// Or if next char (if exists) is lowercase (handles acronyms)
			prev := str[i-1]
			if (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9') {
				b.WriteByte('_')
			}
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}
