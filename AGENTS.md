# Contributing & agent conventions

Conventions for working in `go-jscontact`. This is a standards-implementing
library (RFC 9553 / RFC 9555), not an application — decisions favor wire and
schema fidelity over ergonomic shortcuts.

## Build & test

```sh
go build ./...
go test ./...
go vet ./...
gofmt -l .          # must print nothing
```

CI runs three required checks on every PR: `static` (`go vet` + `gofmt`),
`test` (`go test ./...`), and `lint` (golangci-lint).

## Code conventions

- **Go 1.26+.** Standard library only in the core `jscontact` package. The
  sole external dependency, `github.com/emersion/go-vcard`, is confined to the
  `jscontact/vcard` sub-package (RFC 9555 conversion).
- **Per-file header.** Every `.go` file begins with exactly:

  ```go
  // Copyright 2026 The go-jscontact Authors
  // SPDX-License-Identifier: Apache-2.0
  ```

- **Lenient unmarshal, strict marshal.** Decode whatever the wire provides;
  route unknown members to the open extension field; validate at the marshal
  boundary and via an opt-in `Validate`.
- **Byte-stable round-trip.** Open-extension fields are `json.RawMessage`, not
  `map[string]any`. `@type` discriminators marshal first.
- **Wire fidelity.** `String[Boolean]` members are sets; `components` arrays are
  ordered — preserve both verbatim.

## Commits & PRs

- Imperative, concise commit subjects (`Add Card metadata types`).
- One logical unit per commit; each PR is independently reviewable and keeps CI
  green.
- Reference the relevant spec section in the body where it aids review.

## License

By contributing you agree your contributions are licensed under Apache-2.0.
