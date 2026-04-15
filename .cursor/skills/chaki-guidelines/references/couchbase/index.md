# Couchbase Module

Use this page when the task affects Couchbase cluster wiring, bucket access, or tracing.

## Read Next

- Cluster config, options, and tracing hooks: [cluster-config-and-tracing.md](cluster-config-and-tracing.md)

## Core Rules

- Let the module own cluster construction.
- Keep connection settings and timeouts config-driven.
- Use request context and parent spans when performing bucket and collection operations.
