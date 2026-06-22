package string_parser

import (
	"sort"
	"strings"
)

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
// Variants are sorted most-specific first: primary key is count of optional fields present (more = first),
// secondary key is sum of type-verb weights (date format > :d > :S > any), so more-restrictive patterns
// are tried before less-restrictive ones of equal length.
func expandOptionalSections(format string) []string {
	tokens := strings.Split(format, " ")

	type tok struct {
		content  string
		optional bool
	}
	parsed := make([]tok, len(tokens))
	optIndices := []int{}
	optWeights := []int{}
	for i, t := range tokens {
		if strings.HasSuffix(t, ";opt") {
			content := t[:len(t)-4]
			parsed[i] = tok{content, true}
			optIndices = append(optIndices, i)
			optWeights = append(optWeights, optTokenWeight(content))
		} else {
			parsed[i] = tok{t, false}
		}
	}

	if len(optIndices) == 0 {
		return []string{format}
	}

	n := len(optIndices)

	type variant struct {
		score int
		text  string
	}
	variants := make([]variant, 1<<n)

	for mask := 0; mask < (1 << n); mask++ {
		present := map[int]bool{}
		popcount, weightSum := 0, 0
		for b, idx := range optIndices {
			if mask&(1<<b) != 0 {
				present[idx] = true
				popcount++
				weightSum += optWeights[b]
			}
		}
		parts := []string{}
		for i, t := range parsed {
			if !t.optional || present[i] {
				parts = append(parts, t.content)
			}
		}
		variants[mask] = variant{popcount*1000 + weightSum, strings.Join(parts, " ")}
	}

	sort.Slice(variants, func(i, j int) bool {
		return variants[i].score > variants[j].score
	})

	result := make([]string, len(variants))
	for i, v := range variants {
		result[i] = v.text
	}
	return result
}

// optTokenWeight returns the type-verb restrictiveness of an optional token.
// Higher weight = more restrictive = tried first.
func optTokenWeight(token string) int {
	start := strings.Index(token, ":")
	if start == -1 {
		return 0
	}
	rest := token[start+1:]
	end := strings.Index(rest, "}")
	if end == -1 {
		return 0
	}
	verb := rest[:end]
	if strings.Contains(verb, "%") {
		return 3
	}
	switch verb {
	case "d":
		return 2
	case "S":
		return 1
	}
	return 0
}
