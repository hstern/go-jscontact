// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// Addressing and calendaring property objects (RFC 9553 §2.4–§2.5).

// Address is a physical/mailing address (RFC 9553 §2.5.1). Components is an
// ordered slice: array order is semantic and preserved verbatim.
type Address struct {
	Type             string                     `json:"@type,omitempty"` // "Address"
	Components       []AddressComponent         `json:"components,omitempty"`
	IsOrdered        bool                       `json:"isOrdered,omitempty"`
	CountryCode      string                     `json:"countryCode,omitempty"`
	Coordinates      string                     `json:"coordinates,omitempty"`
	TimeZone         string                     `json:"timeZone,omitempty"`
	Contexts         Contexts                   `json:"contexts,omitempty"`
	Full             string                     `json:"full,omitempty"`
	DefaultSeparator string                     `json:"defaultSeparator,omitempty"`
	Pref             uint64                     `json:"pref,omitempty"`
	PhoneticScript   string                     `json:"phoneticScript,omitempty"`
	PhoneticSystem   string                     `json:"phoneticSystem,omitempty"`
	Extra            map[string]json.RawMessage `json:"-"`
}

func (a Address) MarshalJSON() ([]byte, error) {
	if a.Type == "" {
		a.Type = "Address"
	}
	type alias Address
	return marshalExtra(alias(a), a.Extra)
}

func (a *Address) UnmarshalJSON(b []byte) error {
	type alias Address
	return unmarshalExtra(b, (*alias)(a), &a.Extra)
}

// Calendar is a link to the entity's calendar or free/busy data (RFC 9553
// §2.4.1).
type Calendar struct {
	Type      string                     `json:"@type,omitempty"` // "Calendar"
	Kind      string                     `json:"kind,omitempty"`  // "calendar" | "freeBusy"
	URI       string                     `json:"uri"`
	MediaType string                     `json:"mediaType,omitempty"`
	Contexts  Contexts                   `json:"contexts,omitempty"`
	Pref      uint64                     `json:"pref,omitempty"`
	Label     string                     `json:"label,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

func (c Calendar) MarshalJSON() ([]byte, error) {
	if c.Type == "" {
		c.Type = "Calendar"
	}
	type alias Calendar
	return marshalExtra(alias(c), c.Extra)
}

func (c *Calendar) UnmarshalJSON(b []byte) error {
	type alias Calendar
	return unmarshalExtra(b, (*alias)(c), &c.Extra)
}

// SchedulingAddress is an address for scheduling, e.g. an iMIP email
// (RFC 9553 §2.4.2).
type SchedulingAddress struct {
	Type     string                     `json:"@type,omitempty"` // "SchedulingAddress"
	URI      string                     `json:"uri"`
	Contexts Contexts                   `json:"contexts,omitempty"`
	Pref     uint64                     `json:"pref,omitempty"`
	Label    string                     `json:"label,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (s SchedulingAddress) MarshalJSON() ([]byte, error) {
	if s.Type == "" {
		s.Type = "SchedulingAddress"
	}
	type alias SchedulingAddress
	return marshalExtra(alias(s), s.Extra)
}

func (s *SchedulingAddress) UnmarshalJSON(b []byte) error {
	type alias SchedulingAddress
	return unmarshalExtra(b, (*alias)(s), &s.Extra)
}
