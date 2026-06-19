# Notes

## Tagging and versioning

Use: \
`git tag v0.1.0` \
or \
`git tag -a v1.0.0 -m "Description"`

Then: \
`git push origin v0.1.0`

To clone/go get specific version: \
`github.com/MattThePerson/string_parser@v0.1.0`

## Formats I hope to be able to support (maybe not all)

```yml
filename_formats = [
  '{primary_actors} - {studio:S} - [{year:d}];opt [{date_released_short:%Y-%m}];opt [{date_released:%Y-%m-%d}];opt [{date_released_alt:%Y.%m.%d}];opt [{line:S}];opt {title} [{secondary_actors}];opt {{{source_id:S}}};opt',
  '[{studio:S}] [{year:d}];opt [{date_released:%Y-%m-%d}];opt [{date_released_alt:%Y.%m.%d}];opt [{line:S}];opt {title} [{primary_actors}];opt {{{source_id:S}}};opt',
  '[{date_released:%Y-%m-%d}] {title} {{{source_id:S}}}',
  '[{date_released:%Y-%m-%d} {time_released}] {title} [{source_id_sec}];opt [{source_id}]',
  '{primary_actors} [{dvd_code:S}] [{date_released:%Y-%m-%d}];opt [{date_released_alt:%Y.%m.%d}];opt [{studio:S}];opt {title}',
  '[{dvd_code:S}]',
  '{primary_actors} - [{year:d}];opt [{date_released:%Y-%m-%d}];opt {title} ({year:d});opt',
  '{title} ({year:d}) [{primary_actors}];opt',
  '{title}',
]
```
