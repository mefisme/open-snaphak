---
layout: guide
title: Code Signing Policy
description: How Snapmap+'s Windows binaries are code-signed. Free code signing provided by SignPath.io, certificate by SignPath Foundation.
---

# Code Signing Policy

Free code signing provided by SignPath.io, certificate by SignPath Foundation.

Snapmap+'s Windows binaries — the `snapmap-plus.exe` installer and the two overlay DLLs — are code-signed
with a certificate held by SignPath Foundation, applied automatically in the GitHub Actions release pipeline.
Every release is built from source in the open, and you can verify a build's provenance with
`gh attestation verify --repo doom-snapmap/snapmap-plus <file>`.
