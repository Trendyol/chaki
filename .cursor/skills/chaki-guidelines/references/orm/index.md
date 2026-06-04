# ORM Module

Use this page when the task affects SQL database access, transactions, or repository abstractions.

## Read Next

- Context-aware DB access and transaction boundaries: [gorm-provider-and-transactions.md](gorm-provider-and-transactions.md)
- Generic repositories and repository composition: [repository-pattern.md](repository-pattern.md)

## Core Rules

- Prefer Chaki's ORM abstractions over raw database plumbing at call sites.
- Keep transactional boundaries at the service layer.
- Keep reusable repository behavior inside repository types rather than scattering GORM usage across services.
