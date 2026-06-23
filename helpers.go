package string_parser

// getLastSubstringStartIndex returns the start index of the LAST occurrence of substr within str, or -1 if not found
func getLastSubstringStartIndex(str []rune, substr []rune) int {
	if len(substr) == 0 {
		return len(str)
	}
	for i := len(str) - len(substr); i >= 0; i-- {
		if runesAreEqual(str[i:i+len(substr)], substr) {
			return i
		}
	}
	return -1
}

func runesAreEqual(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// unescapeBraces replaces {{ with { and }} with } in a rune slice (used for format literal patterns)
func unescapeBraces(s []rune) []rune {
	result := []rune{}
	for i := 0; i < len(s); {
		if s[i] == '{' && i+1 < len(s) && s[i+1] == '{' {
			result = append(result, '{')
			i += 2
		} else if s[i] == '}' && i+1 < len(s) && s[i+1] == '}' {
			result = append(result, '}')
			i += 2
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return result
}
