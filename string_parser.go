package string_parser

type StringParser struct {
	Formats []string
}

func NewStringParser(formats []string) *StringParser {
	formats = ExpandFormats(formats)
	return &StringParser{
		Formats: formats,
	}
}

func (sp *StringParser) Parse(input string) (map[string]any, error) {
	data := map[string]any{
		"nothing": "here",
	}

	// ...

	return data, nil
}

func (sp *StringParser) Format(data map[string]any) (string, error) {

	// ...

	return "Here's your string, pervert", nil
}

func ExpandFormats(formats []string) []string {

	// ...

	return formats
}
