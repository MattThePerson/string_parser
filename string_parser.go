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
	formats := []string{format}
	if strings.Contains(format, "/") {
		formats = append(formats, strings.ReplaceAll(format, "/", "\\"))
	} else if strings.Contains(format, "\\") {
		formats = append(formats, strings.ReplaceAll(format, "\\", "/"))
	}
	return formats
}

func ExpandFormatsList(formats []string) []string {
	expanded := []string{}
	for _, f := range formats {
		expanded = append(expanded, ExpandFormats(f)...)
	}
	return expanded
}
