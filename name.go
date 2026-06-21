// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// Name and the other naming-related property objects (RFC 9553 §2.2).

// Name is the entity's name (RFC 9553 §2.2.1). The Components slice is
// ordered: array order is semantic and must be preserved verbatim.
type Name struct {
	Type             string                     `json:"@type,omitempty"` // "Name"
	Components       []NameComponent            `json:"components,omitempty"`
	IsOrdered        bool                       `json:"isOrdered,omitempty"`
	DefaultSeparator string                     `json:"defaultSeparator,omitempty"`
	Full             string                     `json:"full,omitempty"`
	SortAs           map[string]string          `json:"sortAs,omitempty"`
	PhoneticScript   string                     `json:"phoneticScript,omitempty"`
	PhoneticSystem   string                     `json:"phoneticSystem,omitempty"`
	Extra            map[string]json.RawMessage `json:"-"`
}

func (n Name) MarshalJSON() ([]byte, error) {
	if n.Type == "" {
		n.Type = "Name"
	}
	type alias Name
	return marshalExtra(alias(n), n.Extra)
}

func (n *Name) UnmarshalJSON(b []byte) error {
	type alias Name
	return unmarshalExtra(b, (*alias)(n), &n.Extra)
}

// Nickname is an informal name (RFC 9553 §2.2.2).
type Nickname struct {
	Type     string                     `json:"@type,omitempty"` // "Nickname"
	Name     string                     `json:"name"`
	Contexts Contexts                   `json:"contexts,omitempty"`
	Pref     uint64                     `json:"pref,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (n Nickname) MarshalJSON() ([]byte, error) {
	if n.Type == "" {
		n.Type = "Nickname"
	}
	type alias Nickname
	return marshalExtra(alias(n), n.Extra)
}

func (n *Nickname) UnmarshalJSON(b []byte) error {
	type alias Nickname
	return unmarshalExtra(b, (*alias)(n), &n.Extra)
}

// Organization is a company or group the entity is associated with
// (RFC 9553 §2.2.3). Units are in descending hierarchical order.
type Organization struct {
	Type     string                     `json:"@type,omitempty"` // "Organization"
	Name     string                     `json:"name,omitempty"`
	Units    []OrgUnit                  `json:"units,omitempty"`
	SortAs   string                     `json:"sortAs,omitempty"`
	Contexts Contexts                   `json:"contexts,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (o Organization) MarshalJSON() ([]byte, error) {
	if o.Type == "" {
		o.Type = "Organization"
	}
	type alias Organization
	return marshalExtra(alias(o), o.Extra)
}

func (o *Organization) UnmarshalJSON(b []byte) error {
	type alias Organization
	return unmarshalExtra(b, (*alias)(o), &o.Extra)
}

// Title is a job title or role (RFC 9553 §2.2.5).
type Title struct {
	Type           string                     `json:"@type,omitempty"` // "Title"
	Name           string                     `json:"name"`
	Kind           string                     `json:"kind,omitempty"` // "title" (default) | "role"
	OrganizationID string                     `json:"organizationId,omitempty"`
	Extra          map[string]json.RawMessage `json:"-"`
}

func (t Title) MarshalJSON() ([]byte, error) {
	if t.Type == "" {
		t.Type = "Title"
	}
	type alias Title
	return marshalExtra(alias(t), t.Extra)
}

func (t *Title) UnmarshalJSON(b []byte) error {
	type alias Title
	return unmarshalExtra(b, (*alias)(t), &t.Extra)
}

// SpeakToAs describes how to address or refer to the entity (RFC 9553
// §2.2.4).
type SpeakToAs struct {
	Type              string                     `json:"@type,omitempty"` // "SpeakToAs"
	GrammaticalGender string                     `json:"grammaticalGender,omitempty"`
	Pronouns          map[string]Pronouns        `json:"pronouns,omitempty"`
	Extra             map[string]json.RawMessage `json:"-"`
}

func (s SpeakToAs) MarshalJSON() ([]byte, error) {
	if s.Type == "" {
		s.Type = "SpeakToAs"
	}
	type alias SpeakToAs
	return marshalExtra(alias(s), s.Extra)
}

func (s *SpeakToAs) UnmarshalJSON(b []byte) error {
	type alias SpeakToAs
	return unmarshalExtra(b, (*alias)(s), &s.Extra)
}
