package string_parser

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Parse attempts to parse data from a string using StringParser's list of formats
func (sp *StringParser) Parse(input string) (map[string]any, error) {

	for _, format := range sp.formats {
		data, err := parseUsingFormat(input, format)
		if err == nil {
			return data, nil
		}
	}

	return map[string]any{}, fmt.Errorf("no matching format found for input: \"%s\"", input)
}

// parseUsingFormat parses data from a string given a single format
func parseUsingFormat(input string, format string) (map[string]any, error) {

	data := map[string]any{}

	format_runes := []rune(format)
	input_runes := []rune(input)

	extractor_items := getExtractorItems(format_runes)

	start_i := 0
	for _, eitem := range extractor_items {
		if eitem.Name == "" {
			start_i = len(eitem.EndPattern)
			continue
		}
		// get end_i
		end_i := start_i + getSubtringStartIndex(input_runes[start_i:], eitem.EndPattern)
		if end_i < start_i {
			return data, fmt.Errorf("end pattern not in string: %s", string(eitem.EndPattern))
		}
		// extract value
		subrunes := input_runes[start_i:end_i]
		value, err := runesToType(subrunes, eitem.TypeVerb)
		if err != nil {
			return data, err
		}
		data[eitem.Name] = value
		start_i = end_i + len(eitem.EndPattern)
	}

	return data, nil
}

// ExtractorItem
type ExtractorItem struct {
	Name       string
	TypeVerb   string
	EndPattern []rune
}

// getExtractorItems
func getExtractorItems(format_str []rune) []ExtractorItem {

	items := []ExtractorItem{}

	start_i := getNextAttributeStartIndex(format_str, 0)
	begin_pattern := unescapeBraces(format_str[:start_i])
	items = append(items, ExtractorItem{"", "", begin_pattern})

	for start_i < len(format_str) {
		// get end index
		end_i := start_i + getSubtringStartIndex(format_str[start_i:], []rune{'}'})
		if end_i < start_i {
			log.Fatalf("error in string format: no closing brace found }")
		}
		// extract brace contents
		substring := format_str[start_i+1 : end_i]
		type_verb := "any"
		delim_idx := getSubtringStartIndex(substring, []rune{':'})
		if delim_idx != -1 {
			type_verb = string(substring[delim_idx+1:])
			substring = substring[:delim_idx]
		}
		// get end pattern
		start_i = getNextAttributeStartIndex(format_str, end_i)
		end_pattern := unescapeBraces(format_str[end_i+1 : start_i])
		items = append(items, ExtractorItem{
			Name:       string(substring),
			TypeVerb:   type_verb,
			EndPattern: end_pattern,
		})
	}

	return items
}

// getNextAttributeStartIndex
func getNextAttributeStartIndex(str []rune, start_i int) int {
	for i := start_i; i < len(str); i++ {
		if str[i] == '{' && (len(str)-1 == i || str[i+1] != '{') {
			return i
		}
	}
	return len(str)
}

// getSubtringStartIndex returns the start index of a subtring within another string (string means []rune for utf compat)
func getSubtringStartIndex(str []rune, substr []rune) int {
	if len(substr) == 0 {
		return len(str)
	}
	for i := 0; i < len(str)-len(substr)+1; i++ {
		if runesAreEqual(str[i:i+len(substr)], substr) {
			return i
		}
	}
	return -1
}

type parsedDateFormat struct {
	pattern string
	verbs   []string // verb letter per capture group, e.g. ["Y","m","d"]
}

// parseDateFormat converts a strftime-style format to a regex with capture groups and tracks verb order
func parseDateFormat(format string) parsedDateFormat {
	pattern := "^"
	verbs := []string{}
	for i := 0; i < len(format); {
		if format[i] == '%' && i+1 < len(format) {
			verb := string(format[i+1])
			switch format[i+1] {
			case 'Y':
				pattern += `(\d{4})`
			case 'm', 'd':
				pattern += `(\d{2})`
			default:
				pattern += "(" + regexp.QuoteMeta(verb) + ")"
			}
			verbs = append(verbs, verb)
			i += 2
		} else {
			pattern += regexp.QuoteMeta(string(format[i]))
			i++
		}
	}
	return parsedDateFormat{pattern + "$", verbs}
}

// runesToType
func runesToType(input []rune, type_verb string) (any, error) {

	input_str := string(input)

	if strings.Contains(type_verb, "%") {
		df := parseDateFormat(type_verb)
		rx, _ := regexp.Compile(df.pattern)
		matches := rx.FindStringSubmatch(input_str)
		if matches == nil {
			return "", fmt.Errorf("value %q does not match date format %q", input_str, type_verb)
		}
		for i, verb := range df.verbs {
			val, _ := strconv.Atoi(matches[i+1])
			switch verb {
			case "m":
				if val < 1 || val > 12 {
					return "", fmt.Errorf("month %d out of range in %q", val, input_str)
				}
			case "d":
				if val < 1 || val > 31 {
					return "", fmt.Errorf("day %d out of range in %q", val, input_str)
				}
			}
		}
		return input_str, nil
	}

	switch type_verb {

	case "any":
		return input_str, nil

	case "d":
		value, err := strconv.Atoi(input_str)
		if err != nil {
			return "", fmt.Errorf("value %q is not a valid integer", input_str)
		}
		return value, nil

	case "S":
		if strings.ContainsAny(input_str, " \t") {
			return "", fmt.Errorf("value %q contains whitespace, :S requires a no-space string", input_str)
		}
		return input_str, nil

	}

	return input_str, nil

}
