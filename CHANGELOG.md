# Changelog

All notable changes to this project are documented here. The format is based
on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-06-21

Initial release.

### Added

- Typed `Card` object model for RFC 9553 §2, with the §1.4 value types
  (`Id`, `UTCDateTime`, `PatchObject`, components) — hand-written, each
  godoc'd to its spec section.
- Custom JSON codec: `@type` discriminators emit first, and unknown members
  round-trip losslessly through each object's open `Extra` map. Ordered
  `components` arrays and `String[Boolean]` sets are preserved verbatim.
- `Parse` entry point and `DecodeExtra`/`EncodeExtra` typed-extension helpers.
- `Card.Validate` enforcing the §2 MUSTs (required `@type`/`version`/`uid`,
  `members`⇒`kind==group`, Id-key validity, `pref` range, required nested
  fields), reporting a `*ValidationError` with field paths.
- `jscontact/vcard` sub-package: RFC 9555 `FromVCard`/`ToVCard` conversion,
  depending on `github.com/emersion/go-vcard` (the core package stays
  dependency-free).
- RFC 9553 / RFC 9555 conformance fixtures with a dedicated CI job.
- `SpecVersion` constant (`"RFC 9553"`).

[0.1.0]: https://github.com/hstern/go-jscontact/releases/tag/v0.1.0
