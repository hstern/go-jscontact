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

## Design

- **Core is dependency-free.** The `jscontact` package depends only on the Go
  standard library.
- **vCard conversion is opt-in.** Only the `jscontact/vcard` sub-package pulls
  in [`github.com/emersion/go-vcard`](https://github.com/emersion/go-vcard);
  importing the core object model stays dependency-free.
- **Wire fidelity first.** Lenient unmarshal, strict marshal; byte-stable
  round-trip; `@type` discriminators emit first; ordered components and
  `String[Boolean]` sets preserved verbatim.

## Install

```sh
go get github.com/hstern/go-jscontact
```

Requires Go 1.26 or newer.

## Quickstart

### Parse, inspect, validate

```go
package main

import (
	"fmt"

	"github.com/hstern/go-jscontact"
)

func main() {
	data := []byte(`{
		"@type": "Card",
		"version": "1.0",
		"uid": "22B2C7DF-9120-4969-8460-05956FE6B065",
		"name": {"components": [
			{"kind": "given", "value": "John"},
			{"kind": "surname", "value": "Doe"}
		]},
		"emails": {"e1": {"address": "jdoe@example.com", "contexts": {"work": true}}}
	}`)

	card, err := jscontact.Parse(data)
	if err != nil {
		panic(err)
	}
	if err := card.Validate(); err != nil {
		panic(err) // *jscontact.ValidationError, with field paths
	}

	fmt.Println(card.Name.Components[0].Value) // John
	fmt.Println(card.Emails["e1"].Address)     // jdoe@example.com
}
```

### Build and marshal

`encoding/json` works directly on a `Card`; the codec emits `@type` first and
preserves any extension members.

```go
card := &jscontact.Card{
	Version: "1.0", // @type defaults to "Card" on marshal
	UID:     "abc",
	Name:    &jscontact.Name{Full: "Jane Roe"},
}
b, _ := json.Marshal(card)
// {"@type":"Card","version":"1.0","uid":"abc","name":{"@type":"Name","full":"Jane Roe"}}
```

### Typed extensions

Unknown members round-trip losslessly through each object's `Extra` map. Decode
them into your own types instead of reaching for `any`:

```go
var rank struct {
	Score int `json:"score"`
}
if ok, _ := jscontact.DecodeExtra(card.Extra, "x-rank", &rank); ok {
	fmt.Println(rank.Score)
}
```

### Convert to and from vCard (RFC 9555)

```go
import (
	"github.com/emersion/go-vcard"
	jsvcard "github.com/hstern/go-jscontact/vcard"
)

src, _ := vcard.NewDecoder(r).Decode() // a vcard.Card from go-vcard
card, _ := jsvcard.FromVCard(src)      // → *jscontact.Card

out, _ := jsvcard.ToVCard(card)        // *jscontact.Card → vcard.Card
_ = vcard.NewEncoder(w).Encode(out)
```

Conversion covers the common property set and is intentionally asymmetric where
RFC 9555 is; see the package documentation for coverage and the lossy
directions.

## Versioning & stability

This is the `v0.1.x` series: the API may still change before `v1.0.0`. The
library tracks RFC 9553 via `jscontact.SpecVersion`. SemVer is independent of
the spec version.

## License

[Apache-2.0](LICENSE). Copyright 2026 The go-jscontact Authors.
