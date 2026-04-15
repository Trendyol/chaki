# Repository Pattern

Use this page when the task creates or extends repository types.

## Preferred Pattern

- Use the generic repository abstraction from `github.com/Trendyol/chaki/modules/orm/repository` for standard CRUD behavior through `repository.New[Id, T](gp)`.
- Keep custom data access in a repository type that composes or embeds the generic repository.

## Repository Rules

- Repositories should accept `context.Context` and keep data-access logic close to the persistence boundary.
- Use the generic repository for common operations before dropping to raw GORM methods.
- Use lower-level helpers such as `Context(ctx)` or `ParseQuery(ctx, q)` only when the generic abstraction is not enough.

## Design Guidance

- Keep business logic in services and persistence logic in repositories.
- If repeated query names, column labels, or identifiers appear across repository methods, extract them into constants.
