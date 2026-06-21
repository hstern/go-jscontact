// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// Additional-information property objects (RFC 9553 §2.8).

// Anniversary is a significant date such as a birthday (RFC 9553 §2.8.1).
type Anniversary struct {
	Type  string                     `json:"@type,omitempty"` // "Anniversary"
	Kind  string                     `json:"kind"`            // "birth" | "death" | "wedding"
	Date  AnniversaryDate            `json:"date"`
	Place *Address                   `json:"place,omitempty"`
	Extra map[string]json.RawMessage `json:"-"`
}

func (a Anniversary) MarshalJSON() ([]byte, error) {
	if a.Type == "" {
		a.Type = "Anniversary"
	}
	type alias Anniversary
	return marshalExtra(alias(a), a.Extra)
}

func (a *Anniversary) UnmarshalJSON(b []byte) error {
	type alias Anniversary
	return unmarshalExtra(b, (*alias)(a), &a.Extra)
}

// AnniversaryDate holds an Anniversary's date, which RFC 9553 §2.8.1 allows
// to be either a PartialDate (the default type) or a Timestamp. Exactly one
// of the two fields is set after a successful decode.
type AnniversaryDate struct {
	PartialDate *PartialDate
	Timestamp   *Timestamp
}

func (d AnniversaryDate) MarshalJSON() ([]byte, error) {
	switch {
	case d.Timestamp != nil:
		return json.Marshal(d.Timestamp)
	case d.PartialDate != nil:
		return json.Marshal(d.PartialDate)
	default:
		return []byte("null"), nil
	}
}

func (d *AnniversaryDate) UnmarshalJSON(b []byte) error {
	var probe struct {
		Type string  `json:"@type"`
		UTC  *string `json:"utc"`
	}
	if err := json.Unmarshal(b, &probe); err != nil {
		return err
	}
	// defaultType is PartialDate; a Timestamp is identified by @type or, when
	// @type is omitted, by the presence of the "utc" member.
	if probe.Type == "Timestamp" || (probe.Type == "" && probe.UTC != nil) {
		var t Timestamp
		if err := json.Unmarshal(b, &t); err != nil {
			return err
		}
		d.PartialDate, d.Timestamp = nil, &t
		return nil
	}
	var p PartialDate
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	d.Timestamp, d.PartialDate = nil, &p
	return nil
}

// PersonalInfo is a piece of personal information such as a hobby or area of
// expertise (RFC 9553 §2.8.4).
type PersonalInfo struct {
	Type   string                     `json:"@type,omitempty"` // "PersonalInfo"
	Kind   string                     `json:"kind"`            // "expertise" | "hobby" | "interest"
	Value  string                     `json:"value"`
	Level  string                     `json:"level,omitempty"` // "high" | "medium" | "low"
	ListAs uint64                     `json:"listAs,omitempty"`
	Label  string                     `json:"label,omitempty"`
	Extra  map[string]json.RawMessage `json:"-"`
}

func (p PersonalInfo) MarshalJSON() ([]byte, error) {
	if p.Type == "" {
		p.Type = "PersonalInfo"
	}
	type alias PersonalInfo
	return marshalExtra(alias(p), p.Extra)
}

func (p *PersonalInfo) UnmarshalJSON(b []byte) error {
	type alias PersonalInfo
	return unmarshalExtra(b, (*alias)(p), &p.Extra)
}

// Note is a free-text note about the entity (RFC 9553 §2.8.3).
type Note struct {
	Type    string                     `json:"@type,omitempty"` // "Note"
	Note    string                     `json:"note"`
	Created UTCDateTime                `json:"created,omitempty"`
	Author  *Author                    `json:"author,omitempty"`
	Extra   map[string]json.RawMessage `json:"-"`
}

func (n Note) MarshalJSON() ([]byte, error) {
	if n.Type == "" {
		n.Type = "Note"
	}
	type alias Note
	return marshalExtra(alias(n), n.Extra)
}

func (n *Note) UnmarshalJSON(b []byte) error {
	type alias Note
	return unmarshalExtra(b, (*alias)(n), &n.Extra)
}
