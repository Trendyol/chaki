# Request Parsing And Validation

Use this page when changing request structs, validation tags, or response status behavior.

## Request Tag Rules

- Use `json:"field"` for request body fields.
- Use `query:"field"` for query parameters.
- Use `param:"field"` for path parameters.
- Use `reqHeader:"Header-Name"` for request headers.
- Use `cookie:"cookie_name"` for cookies.
- Use `validate:"..."` for validator rules.

## Validation Rules

- Typed handlers automatically parse request data into the request struct and run validation when validation tags are present.
- Validation is based on `github.com/go-playground/validator/v10`.
- Use custom validation rules through Chaki validation hooks when built-in tags are not enough.

## Response Rules

- Keep request structs aligned with the actual transport contract so Swagger generation remains accurate.
- Use `route.NoParam` when a route has no inputs at all.
- If a route accepts only query or path input, model those fields explicitly with `query` or `param` tags instead of forcing a JSON body.
- If response status depends on response content, implement `Status() int` on the response type instead of repeating manual status-setting logic.
