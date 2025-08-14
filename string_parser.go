package string_parser

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

	// ...

	return []string{format}
}

func ExpandFormatsList(formats []string) []string {

	// ...

	return formats
}
