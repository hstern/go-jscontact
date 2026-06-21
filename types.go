// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// This file holds the RFC 9553 §1.4 common value types and the small leaf
// component objects that the higher-level property types compose.

// UTCDateTime is an RFC 3339 date-time in UTC, ending in a literal "Z"
// (RFC 9553 §1.4.5), e.g. "2026-06-21T01:18:22Z".
type UTCDateTime string

// PatchObject is a set of patches keyed by JSON Pointer (RFC 6901) relative
// to a Card, used by the Card.Localizations property (RFC 9553 §1.4.3,
// §2.7.1). A key maps to the replacement value at that pointer; the raw
// JSON value is preserved so byte-stable round-tripping is possible.
type PatchObject map[string]json.RawMessage

// Boolean-valued string sets (RFC 9553 "String[Boolean]"). Every value is
// true; the type models a set whose membership is the set of keys. These
// are plain maps so encoding/json renders them with deterministically
// ordered keys.
type (
	// Contexts is the set of contexts in which a property applies, e.g.
	// {"work": true} (RFC 9553 §1.5.1).
	Contexts map[string]bool

	// Features is the set of a phone's capabilities, e.g. {"voice": true}
	// (RFC 9553 §2.3.3).
	Features map[string]bool
)

// NameComponent is one ordered piece of a Name (RFC 9553 §2.2.1).
type NameComponent struct {
	Type     string `json:"@type,omitempty"` // "NameComponent"
	Kind     string `json:"kind"`            // given, surname, title, given2, surname2, credential, generation, separator
	Value    string `json:"value"`
	Phonetic string `json:"phonetic,omitempty"`
}

// AddressComponent is one ordered piece of an Address (RFC 9553 §2.5.1).
type AddressComponent struct {
	Type     string `json:"@type,omitempty"` // "AddressComponent"
	Kind     string `json:"kind"`            // number, name, locality, region, postcode, country, separator, ...
	Value    string `json:"value"`
	Phonetic string `json:"phonetic,omitempty"`
}

// OrgUnit is a unit within an Organization, in descending hierarchical
// order (RFC 9553 §2.2.3).
type OrgUnit struct {
	Type   string `json:"@type,omitempty"` // "OrgUnit"
	Name   string `json:"name"`
	SortAs string `json:"sortAs,omitempty"`
}

// Pronouns expresses how to refer to the entity (RFC 9553 §2.2.4).
type Pronouns struct {
	Type     string   `json:"@type,omitempty"` // "Pronouns"
	Pronouns string   `json:"pronouns"`
	Contexts Contexts `json:"contexts,omitempty"`
	Pref     uint64   `json:"pref,omitempty"`
}

// Author identifies who wrote a Note (RFC 9553 §2.8.3).
type Author struct {
	Type string `json:"@type,omitempty"` // "Author"
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`
}

// Timestamp is a point in time (RFC 9553 §1.4.6); one of the value types an
// Anniversary date may take.
type Timestamp struct {
	Type string      `json:"@type,omitempty"` // "Timestamp"
	UTC  UTCDateTime `json:"utc"`
}

// PartialDate is a calendar date that may omit components (RFC 9553 §2.8.1);
// the other value type an Anniversary date may take.
type PartialDate struct {
	Type          string `json:"@type,omitempty"` // "PartialDate"
	Year          uint64 `json:"year,omitempty"`
	Month         uint64 `json:"month,omitempty"`
	Day           uint64 `json:"day,omitempty"`
	CalendarScale string `json:"calendarScale,omitempty"`
}

// Relation expresses how another entity relates to the Card (RFC 9553
// §2.1.8). The relation set values are always true.
type Relation struct {
	Type     string                     `json:"@type,omitempty"` // "Relation"
	Relation map[string]bool            `json:"relation,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (r Relation) MarshalJSON() ([]byte, error) {
	if r.Type == "" {
		r.Type = "Relation"
	}
	type alias Relation
	return marshalExtra(alias(r), r.Extra)
}

func (r *Relation) UnmarshalJSON(b []byte) error {
	type alias Relation
	return unmarshalExtra(b, (*alias)(r), &r.Extra)
}
