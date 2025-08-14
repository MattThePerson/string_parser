// RUNNING
// go test
// go test -v
// go test -cover
package string_parser

import (
	"encoding/json"
	"testing"
)

// #region - FORMAT TESTS ------------------------------------------------------

type FormatTestItem struct {
	Format   string
	Input    map[string]any
	Expected string
}

var (
	FormatTests = []FormatTestItem{
		{
			"{actor} [{year:d}] {title} [{id}].{ext}",
			map[string]any{"actor": "Abbie"},
			"Here's your string, pervert",
		},
	}
)

func TestFormat(t *testing.T) {

	for _, test := range FormatTests {
		parser := NewStringParser(test.Format)
		got, err := parser.Format(test.Input)
		if err != nil {
			t.Errorf("ERR not nil")
		}
		if got != test.Expected {
			t.Errorf("FORMAT: got: \"%s\"  want: \"%s\"\n", got, test.Expected)
		}
	}

}

// #region - PARSE TESTS -------------------------------------------------------

type ParseTestItem struct {
	Format   string
	Input    string
	Expected map[string]any
}

var (
	ParseTests = []ParseTestItem{
		{ // basic text
			"{actor} - {studio} - [{year}] {title}",
			"Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
			map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": "2012", "title": "The Fountain of Testing"},
		},
		{ // start pattern
			"start shit {actor} - {studio} - [{year}] {title} [{sid}]",
			"start shit Actor McActor - LeopardsGate - [2012] The Fountain of Testing [1234]",
			map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": "2012", "title": "The Fountain of Testing", "sid": "1234"},
		},
		{ // :d int
			"{actor} - {studio} - [{year:d}] {title}",
			"Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
			map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
		},
		{ // :d int
			"{actor} - {studio} - [{year:d}] {title}.{ext}",
			"Actor McActor - LeopardsGate - [2012] The Fountain of Testing.mp4",
			map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing", "ext": "mp4"},
		},
		{ // folders in path
			"{parent_dir}/{actor} - {studio} - [{year:d}] {title}",
			"Movies/Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
			map[string]any{"parent_dir": "Movies", "actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
		},
		{ // backslashes
			"{collection}/{subcollection}/{actor} - {studio:%Y-%M-%D} - [{year:d}] {title}",
			"Movies\\Classics\\Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
			map[string]any{"collection": "Movies", "subcollection": "Classics", "actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
		},
		{ //
			"{junk}/{parent_dir}/{actor} - {studio} - [{year:d}] {title}",
			"Media\\Movies\\Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
			map[string]any{"junk": "Media", "parent_dir": "Movies", "actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
		},
	}
)

func TestParse(t *testing.T) {

	for i, test := range ParseTests {
		parser := NewStringParser(test.Format)
		got, err := parser.Parse(test.Input)
		if err != nil {
			t.Errorf("ERR not nil")
		}
		if !mapsEqual(got, test.Expected) {
			got_fmt, _ := json.MarshalIndent(&got, "", "    ")
			exp_fmt, _ := json.MarshalIndent(&test.Expected, "", "    ")
			t.Errorf("PARSE TEST %d/%d:\ngot: %v \nexpected: %v\n", i+1, len(ParseTests), string(got_fmt), string(exp_fmt))
		}
	}

}

func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if vb, ok := b[k]; !ok || vb != v {
			return false
		}
	}
	return true
}

// #region - PARSE TESTS -------------------------------------------------------

type ParseMultiTestItem struct {
	Formats  []string
	Input    string
	Expected map[string]any
}

var (
	ParseTestsMulti = []ParseMultiTestItem{
		{ // basic text
			[]string{"{actor} - {studio} - [{year}] {title}", "{other_stuff}"},
			"Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
			map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": "2012", "title": "The Fountain of Testing"},
		},
		{ // basic text
			[]string{"{actor} - {studio} - [{year}] {title}", "{actor} - {studio} - {title}"},
			"Actor McActor - LeopardsGate - The Fountain of Testing",
			map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "title": "The Fountain of Testing"},
		},
	}
)

func TestParseMulti(t *testing.T) {

	for i, test := range ParseTestsMulti {
		parser := NewStringParserFromList(test.Formats)
		got, err := parser.Parse(test.Input)
		if err != nil {
			t.Errorf("ERR not nil")
		}
		if !mapsEqual(got, test.Expected) {
			got_fmt, _ := json.MarshalIndent(&got, "", "    ")
			exp_fmt, _ := json.MarshalIndent(&test.Expected, "", "    ")
			t.Errorf("PARSE MULTI TEST %d/%d:\ngot: %v \nexpected: %v\n", i+1, len(ParseTestsMulti), string(got_fmt), string(exp_fmt))
		}
	}

}
