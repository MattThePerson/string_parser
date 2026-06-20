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
            t.Errorf("failed to correctly format")
        }
        if got != test.Expected {
            t.Errorf("FORMAT: got: \"%s\"  want: \"%s\"\n", got, test.Expected)
        }
    }

}

// #region - PARSE TESTS -------------------------------------------------------

type ParseTestItem struct {
    Title    string
    Format   string
    Input    string
    Expected map[string]any
}

var (
    ParseTests = []ParseTestItem{
        { // PARSE: 1/15
            "basic text",
            "{actor} - {studio} - [{year}] {title}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
            map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": "2012", "title": "The Fountain of Testing"},
        },
        { // PARSE: 2/15
            "start pattern",
            "start shit {actor} - {studio} - [{year}] {title} [{sid}]",
            "start shit Actor McActor - LeopardsGate - [2012] The Fountain of Testing [1234]",
            map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": "2012", "title": "The Fountain of Testing", "sid": "1234"},
        },
        { // PARSE: 3/15
            ":d int",
            "{actor} - {studio} - [{year:d}] {title}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
            map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
        },
        { // PARSE: 4/15
            ":d int",
            "{actor} - {studio} - [{year:d}] {title}.{ext}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing.mp4",
            map[string]any{"actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing", "ext": "mp4"},
        },
        { // PARSE: 5/15
            "folders in path",
            "{parent_dir}/{actor} - {studio} - [{year:d}] {title}",
            "Movies/Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
            map[string]any{"parent_dir": "Movies", "actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
        },
        { // PARSE: 6/15
            "backslashes",
            "{collection}/{subcollection}/{actor} - {studio} - [{year:d}] {title}",
            "Movies\\Classics\\Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
            map[string]any{"collection": "Movies", "subcollection": "Classics", "actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
        },
        { // PARSE: 7/15 (same as above?)
            "IDK really",
            "{junk}/{parent_dir}/{actor} - {studio} - [{year:d}] {title}",
            "Media\\Movies\\Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
            map[string]any{"junk": "Media", "parent_dir": "Movies", "actor": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
        },
        { // PARSE: 8/15
            ":S type verb (no-spaces string field)",
            "{primary_actors} - {studio:S} - [{year:d}] {title}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
        },
        { // PARSE: 9/15
            "%Y-%m-%d date type verb",
            "{primary_actors} - {studio:S} - [{date_released:%Y-%m-%d}] {title}",
            "Actor McActor - LeopardsGate - [2012-06-19] The Fountain of Testing",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "date_released": "2012-06-19", "title": "The Fountain of Testing"},
        },
        { // PARSE: 10/15
            "%Y-%m partial date type verb",
            "{primary_actors} - {studio:S} - [{date_released_short:%Y-%m}] {title}",
            "Actor McActor - LeopardsGate - [2012-06] The Fountain of Testing",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "date_released_short": "2012-06", "title": "The Fountain of Testing"},
        },
        { // PARSE: 11/15
            "%Y.%m.%d dot-separated date type verb",
            "{primary_actors} - {studio:S} - [{date_released_alt:%Y.%m.%d}] {title}",
            "Actor McActor - LeopardsGate - [2012.06.19] The Fountain of Testing",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "date_released_alt": "2012.06.19", "title": "The Fountain of Testing"},
        },
        { // PARSE: 12/15
            ";opt optional section present",
            "{primary_actors} - {studio:S} - [{year:d}];opt [{date_released:%Y-%m-%d}] {title}",
            "Actor McActor - LeopardsGate - [2012] [2012-06-19] The Fountain of Testing",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "date_released": "2012-06-19", "title": "The Fountain of Testing"},
        },
        { // PARSE: 13/15
            ";opt optional section absent",
            "{primary_actors} - {studio:S} - [{year:d}];opt [{date_released:%Y-%m-%d}] {title}",
            "Actor McActor - LeopardsGate - [2012-06-19] The Fountain of Testing",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "date_released": "2012-06-19", "title": "The Fountain of Testing"},
        },
        { // PARSE: 14/15
            "multiple ;opt sections",
            "{primary_actors} - {studio:S} - [{year:d}];opt [{date_released:%Y-%m-%d}];opt [{date_released_alt:%Y.%m.%d}];opt {title}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing"},
        },
        { // PARSE: 15/15
            "{{{field}}} literal braces wrapping a field",
            "{primary_actors} - {studio:S} - [{year:d}] {title} {{{source_id:S}}}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing {ABC123}",
            map[string]any{"primary_actors": "Actor McActor", "studio": "LeopardsGate", "year": 2012, "title": "The Fountain of Testing", "source_id": "ABC123"},
        },
    }
)

func TestParse(t *testing.T) {
    for i, test := range ParseTests {
        parser := NewStringParser(test.Format)
        got, err := parser.Parse(test.Input)
        if err != nil {
            t.Errorf("Parse error")
        }
        if !mapsEqual(got, test.Expected) {
            got_fmt, _ := json.MarshalIndent(&got, "", "    ")
            exp_fmt, _ := json.MarshalIndent(&test.Expected, "", "    ")
            t.Errorf("PARSE TEST %d/%d:\n>title: %s\n>format: %s\n>input: %s\n>expected: %v\n>got: %v\n",
                i+1, len(ParseTests), test.Title, test.Format, test.Input, string(exp_fmt), string(got_fmt))
        }
    }
}

// #region - PARSE NEGATIVE TESTS ---------------------------------------------

type ParseNegativeTestItem struct {
    Title  string
    Format string
    Input  string
}

var (
    ParseNegativeTests = []ParseNegativeTestItem{
        { // PARSE NEG: 1/15 - separator mismatch (dash vs underscore)
            "separator mismatch",
            "{actor} - {studio} - [{year:d}] {title}",
            "Actor McActor _ LeopardsGate _ [2012] The Fountain of Testing",
        },
        { // PARSE NEG: 2/15 - bracket type mismatch (square vs round)
            "bracket type mismatch",
            "{actor} - {studio} - [{year:d}] {title}",
            "Actor McActor - LeopardsGate - (2012) The Fountain of Testing",
        },
        { // PARSE NEG: 3/15 - required trailing extension absent
            "missing required extension",
            "{actor} - {studio} - [{year:d}] {title}.{ext}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
        },
        { // PARSE NEG: 4/15 - path separator absent from input
            "missing path separator",
            "{dir}/{actor} - {studio} - [{year:d}] {title}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
        },
        { // PARSE NEG: 5/15 - dot date where dash-date expected
            "dot date for %Y-%m-%d",
            "{actor} - {studio} - [{date_released:%Y-%m-%d}] {title}",
            "Actor McActor - LeopardsGate - [2012.06.19] The Fountain of Testing",
        },
        { // PARSE NEG: 6/15 - dash date where dot-date expected
            "dash date for %Y.%m.%d",
            "{actor} - {studio} - [{date_released:%Y.%m.%d}] {title}",
            "Actor McActor - LeopardsGate - [2012-06-19] The Fountain of Testing",
        },
        { // PARSE NEG: 7/15 - year-only where full date expected
            "year-only for %Y-%m-%d",
            "{actor} - {studio} - [{date_released:%Y-%m-%d}] {title}",
            "Actor McActor - LeopardsGate - [2012] The Fountain of Testing",
        },
        { // PARSE NEG: 8/15 - full date where year-month expected
            "full date for %Y-%m",
            "{actor} - {studio} - [{date_released:%Y-%m}] {title}",
            "Actor McActor - LeopardsGate - [2012-06-19] The Fountain of Testing",
        },
        { // PARSE NEG: 9/15 - out-of-range date components
            "out-of-range date components",
            "{actor} - {studio} - [{date_released:%Y-%m-%d}] {title}",
            "Actor McActor - LeopardsGate - [2012-99-99] The Fountain of Testing",
        },
        { // PARSE NEG: 10/15 - non-date string in date field
            "non-date string in date field",
            "{actor} - {studio} - [{date_released:%Y-%m-%d}] {title}",
            "Actor McActor - LeopardsGate - [blahblah] The Fountain of Testing",
        },
        { // PARSE NEG: 11/15 - :d with alphabetic input
            ":d with alphabetic input",
            "{actor} - {studio} - [{year:d}] {title}",
            "Actor McActor - LeopardsGate - [abc] The Fountain of Testing",
        },
        { // PARSE NEG: 12/15 - :d with mixed alphanumeric
            ":d with mixed alphanumeric",
            "{actor} - {studio} - [{year:d}] {title}",
            "Actor McActor - LeopardsGate - [2012abc] The Fountain of Testing",
        },
        { // PARSE NEG: 13/15 - :S field contains spaces
            ":S field contains spaces",
            "{actor} - {studio:S} - [{year:d}] {title}",
            "Actor McActor - Leopards Gate - [2012] The Fountain of Testing",
        },
        { // PARSE NEG: 14/15 - trailing content after closing bracket
            "trailing content after closing bracket",
            "{actor} [{year:d}]",
            "Actor McActor [2012] SomeExtraGarbage",
        },
        { // PARSE NEG: 15/15 - extra quality tag after closing paren
            "extra quality tag after closing paren",
            "{title} ({year:d})",
            "The Movie (2023) [720p BluRay]",
        },
    }
)

func TestNegativeParse(t *testing.T) {
    for i, test := range ParseNegativeTests {
        parser := NewStringParser(test.Format)
        got, err := parser.Parse(test.Input)
        if err == nil || len(got) > 0 {
            got_fmt, _ := json.MarshalIndent(&got, "", "    ")
            t.Errorf("PARSE NEG TEST %d/%d: \"%s\" — expected failure\n>format: %s\n>input: %s\n>got: %v\n",
                i+1, len(ParseNegativeTests), test.Title, test.Format, test.Input, string(got_fmt))
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
    FULL_FORMATS = []string{
        "{rel_parent}/{primary_actors} - {studio:S} - [{year:d}];opt [{date_released_short:%Y-%m}];opt [{date_released:%Y-%m-%d}];opt [{date_released_alt:%Y.%m.%d}];opt [{line:S}];opt {title} [{secondary_actors}];opt {{{source_id:S}}};opt",
        "{rel_parent}/[{studio:S}] [{year:d}];opt [{date_released:%Y-%m-%d}];opt [{date_released_alt:%Y.%m.%d}];opt [{line:S}];opt {title} [{primary_actors}];opt {{{source_id:S}}};opt",
        "{rel_parent}/[{date_released:%Y-%m-%d}] {title} {{{source_id:S}}}",
        "{rel_parent}/[{date_released:%Y-%m-%d} {time_released}] {title} [{source_id_sec}];opt [{source_id}]",
        "{rel_parent}/{primary_actors} [{dvd_code:S}] [{date_released:%Y-%m-%d}];opt [{date_released_alt:%Y.%m.%d}];opt [{studio:S}];opt {title}",
        "{rel_parent}/[{dvd_code:S}]",
        "{rel_parent}/{primary_actors} - [{year:d}];opt [{date_released:%Y-%m-%d}];opt {title} ({year:d});opt",
        "{rel_parent}/{title} ({year:d}) [{primary_actors}];opt",
        "{rel_parent}/{title}",
    }

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
        { // complex test
            FULL_FORMATS,
            "some path/Scarlit Sandals, Chanel Convertible - TwistAndShout - [2023] [WhenThingsHappen] Sleeping in a Room",
            map[string]any{"rel_parent": "some path", "primary_actors": "Scarlit Sandals, Chanel Convertible", "studio": "TwistAndShout", "year": 2023, "line": "WhenThingsHappen", "title": "Sleeping in a Room"},
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
            t.Errorf("PARSE MULTI TEST %d/%d:\n>expected: %v\n>got: %v \n", i+1, len(ParseTestsMulti), string(exp_fmt), string(got_fmt))
        }
    }

}
