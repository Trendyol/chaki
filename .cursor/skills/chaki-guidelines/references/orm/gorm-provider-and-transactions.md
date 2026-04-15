# Gorm Provider And Transactions

Use this page when the task needs context-aware database access or multi-step transactions.

## Access Pattern

- Register the ORM module with the appropriate driver, such as `orm.Module(postgresdriver.New())`.
- Inject `orm.GormProvider` instead of raw `*gorm.DB` when possible.
- Call `GormProvider.Get(ctx)` so the returned DB handle respects the active transaction and the incoming context.

## Transaction Rules

- Inject a `tx.Transactioner` when a service operation must commit or roll back multiple repository calls together.
- Wrap transactional work with `transactioner.Transaction(parentCtx, func(ctx context.Context) error { ... })`.
- Pass the transactional context through to every repository call inside the transaction.

## Testing Guidance

- Use the provided mock transactioner from `modules/orm/tx/mock` for unit tests; the Go package name is `txmock`.
- Keep transaction ownership in services. Repositories should consume the context they are given instead of opening their own transaction boundaries.
