---
title: PostgreSQL 16
description: Official PostgreSQL 16 docs record for SQL semantics, administration behavior, and feature expectations in AthenaWork environments.
doc_type: canonical_source
status: draft
tool: postgresql
vendor: PostgreSQL Global Development Group
version_policy: pinned_minor
version: "16"
docs_base: https://www.postgresql.org/docs/16/
allowed_domains:
  - www.postgresql.org
entrypoints:
  - https://www.postgresql.org/docs/16/
consult_when:
  - checking SQL behavior or PostgreSQL feature semantics
  - reviewing server configuration or administration details
  - local environments use PostgreSQL 16 images or packages
avoid_when:
  - the task is only about AthenaWork workflow rules
notes:
  - Match this record to the image tag used by local compose files.
---

# PostgreSQL 16

## Summary

Use the PostgreSQL 16 manual when AthenaWork tasks depend on database semantics or Postgres administration details.

## References

- https://www.postgresql.org/docs/16/
