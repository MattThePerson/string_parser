package string_parser

import "strings"

type StringParser struct {
	formats []string
}

func NewStringParser(format string) *StringParser {
	formats := ExpandFormats(format)
	return &StringParser{
		formats: formats,
	}
}

func NewStringParserFromList(formats []string) *StringParser {
	formats = ExpandFormatsList(formats)
	return &StringParser{
		formats: formats,
	}
}

func ExpandFormats(format string) []string {
	// Expand ;opt optional sections first, then expand path separators for each variant
	result := []string{}
	for _, f := range expandOptionalSections(format) {
		result = append(result, f)
		if strings.Contains(f, "/") {
			result = append(result, strings.ReplaceAll(f, "/", "\\"))
		} else if strings.Contains(f, "\\") {
			result = append(result, strings.ReplaceAll(f, "\\", "/"))
		}
	}
	return result
}

func ExpandFormatsList(formats []string) []string {
	expanded := []string{}
	for _, f := range formats {
		expanded = append(expanded, ExpandFormats(f)...)
	}
	return expanded
}

// expandOptionalSections generates all 2^n format variants by including/excluding tokens marked ;opt.
// Variants are ordered most-inclusive first so Parse() matches greedily.
func expandOptionalSections(format string) []string {
	tokens := strings.Split(format, " ")

	type tok struct {
		content  string
		optional bool
	}
	parsed := make([]tok, len(tokens))
	optIndices := []int{}
	for i, t := range tokens {
		if strings.HasSuffix(t, ";opt") {
			parsed[i] = tok{t[:len(t)-4], true}
			optIndices = append(optIndices, i)
		} else {
			parsed[i] = tok{t, false}
		}
	}

	if len(optIndices) == 0 {
		return []string{format}
	}

	n := len(optIndices)
	result := []string{}
	for mask := (1 << n) - 1; mask >= 0; mask-- {
		present := map[int]bool{}
		for b, idx := range optIndices {
			if mask&(1<<b) != 0 {
				present[idx] = true
			}
		}
		parts := []string{}
		for i, t := range parsed {
			if !t.optional || present[i] {
				parts = append(parts, t.content)
			}
		}
		result = append(result, strings.Join(parts, " "))
	}
	return result
}
