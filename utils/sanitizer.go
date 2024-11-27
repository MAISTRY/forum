package utils

import (
	"regexp"
)

func InputSanitizer(inp string) string {
	replacements := map[string]string{
		"<": "&lt;",
		">": "&gt;",
		"&": "&amp;",
		`"`: "&quot;",
		"'": "&apos;",
		"=": "&#x3D;",
		"%": "&#x25;",
		":": "&#x3A;",
		";": "&#x3B;",
		"#": "&#x23;",
		"@": "&#x40;",
		"{": "&#x7B;",
		"}": "&#x7D;",
		"[": "&#x5B;",
		"]": "&#x5D;",
		"(": "&#x28;",
		")": "&#x29;",
		"?": "&#x3F;",
		"+": "&#x2B;",
		"-": "&#x2D;",
	}

	pattern := regexp.MustCompile(`[<>&"'=%:;#@{}\[\]()?+-]`)

	sanitizedInp := pattern.ReplaceAllStringFunc(inp, func(match string) string {
		return replacements[match]
	})

	return sanitizedInp
}
