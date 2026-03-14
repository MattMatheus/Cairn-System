---
title: Canonical Sources Index
description: Registry of official vendor-doc records used by AthenaWork agents for external tool semantics.
doc_type: index
status: draft
---

# Canonical Sources Index

## Summary

This folder holds official external source records. Each file identifies the vendor docs AthenaWork agents should trust for a tool, plus the version policy that controls retrieval.

## Intended Audience

- AthenaWork maintainers
- agent authors

## Preconditions

- the tool materially affects development work
- the team can identify the official doc set

## Main Flow

Create one file per tool or per pinned major/minor line when version behavior matters.

Current records:

- [[go]]
- [[podman-5]]
- [[postgresql-16]]
- [[react-19]]
- [[terraform-1-9]]

## Failure Modes

- mixing multiple tools in one file
- omitting version policy
- linking to tutorials or community blogs instead of official docs

## References

- [[CANONICAL_SOURCE_REGISTRY]]
