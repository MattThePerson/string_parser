# String Parser

A Go library for extracting structured data from strings using format patterns.

## Format syntax

Fields are written as `{name}` or `{name:verb}` inside a format string. Literal text outside braces is matched exactly. Use `{{` and `}}` to include literal brace characters.

```
"{primary_actors} - {studio:S} - [{year:d}] {title}"
```

### Optional sections

Append `;opt` to a token to make that section optional. The parser will try all combinations, preferring the most-inclusive match.

```
"{primary_actors} - {studio:S} - [{year:d}];opt [{date_released:%Y-%m-%d}];opt {title}"
```

---

## Type verbs

| Verb | Name | Description |
|---|---|---|
| *(none)* | any | Accepts any string value. This is the default. |
| `:d` | integer | Parses the value as a base-10 integer. Fails if the value is not a valid integer. |
| `:S` | single-word | Accepts only strings with no whitespace. Useful for IDs, codes, or slugs. |
| `:K` | code | Like `:S` (no whitespace), but also rejects plain integers. Use for alphanumeric codes like `ABC-123` where a bare number would be ambiguous with `:d` fields. |
| `:P` | path | Matches up to the **last** occurrence of the field's separator instead of the first. Use for fields that represent filesystem paths, which may themselves contain the separator character (e.g. `/`). |
| `:%Y-%m-%d` | date | Matches a date string against a strftime-style pattern. Supported directives: `%Y` (4-digit year), `%m` (2-digit month, 01–12), `%d` (2-digit day, 01–31). Fails if the value does not match the pattern or contains out-of-range values. |

### `:P` example

With format `{rel_parent:P}/{title}` and input `A:/full/path/to/thing/filename here`:

- **Without `:P`** — stops at first `/`: `rel_parent = "A:"`, `title = "full/path/to/thing/filename here"`
- **With `:P`** — stops at last `/`: `rel_parent = "A:/full/path/to/thing"`, `title = "filename here"`

### Date format examples

| Verb | Matches |
|---|---|
| `:%Y-%m-%d` | `2024-06-19` |
| `:%Y-%m` | `2024-06` |
| `:%Y.%m.%d` | `2024.06.19` |
