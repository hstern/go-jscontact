# go-jscontact

[![CI](https://github.com/hstern/go-jscontact/actions/workflows/ci.yml/badge.svg)](https://github.com/hstern/go-jscontact/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/hstern/go-jscontact.svg)](https://pkg.go.dev/github.com/hstern/go-jscontact)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](LICENSE)

A Go library implementing **[RFC 9553 — JSContact: A JSON Representation of
Contact Data](https://www.rfc-editor.org/rfc/rfc9553.html)** — the typed,
JSON-native successor to vCard.

`go-jscontact` provides the typed `Card` object model and a custom JSON codec
that is byte-stable on round-trip and routes unknown members to an open
extension field. A separate sub-package, `jscontact/vcard`, implements the
**[RFC 9555](https://www.rfc-editor.org/rfc/rfc9555.html)** vCard ⇄ JSContact
conversion.

> **Status: pre-publication.** The first release will be tagged `v0.1.0`. The
> object model, codec, validation, and conversion land across the development
> phases tracked in the repository. This scaffold currently exposes only
> `SpecVersion`.

## Design

- **Core is dependency-free.** The `jscontact` package depends only on the Go
  standard library.
- **vCard conversion is opt-in.** Only the `jscontact/vcard` sub-package pulls
  in [`github.com/emersion/go-vcard`](https://github.com/emersion/go-vcard);
  importing the core object model stays dependency-free.
- **Wire fidelity first.** Lenient unmarshal, strict marshal; byte-stable
  round-trip; `@type` discriminators emit first.

## Install

```sh
go get github.com/hstern/go-jscontact
```

Requires Go 1.26 or newer.

## License

[Apache-2.0](LICENSE). Copyright 2026 The go-jscontact Authors.
