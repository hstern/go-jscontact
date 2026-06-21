// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// Contact-channel property objects (RFC 9553 §2.3).

// EmailAddress is an email address (RFC 9553 §2.3.1).
type EmailAddress struct {
	Type     string                     `json:"@type,omitempty"` // "EmailAddress"
	Address  string                     `json:"address"`
	Contexts Contexts                   `json:"contexts,omitempty"`
	Pref     uint64                     `json:"pref,omitempty"`
	Label    string                     `json:"label,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (e EmailAddress) MarshalJSON() ([]byte, error) {
	if e.Type == "" {
		e.Type = "EmailAddress"
	}
	type alias EmailAddress
	return marshalExtra(alias(e), e.Extra)
}

func (e *EmailAddress) UnmarshalJSON(b []byte) error {
	type alias EmailAddress
	return unmarshalExtra(b, (*alias)(e), &e.Extra)
}

// OnlineService is an account on an online service (RFC 9553 §2.3.2).
type OnlineService struct {
	Type     string                     `json:"@type,omitempty"` // "OnlineService"
	Service  string                     `json:"service,omitempty"`
	URI      string                     `json:"uri,omitempty"`
	User     string                     `json:"user,omitempty"`
	Contexts Contexts                   `json:"contexts,omitempty"`
	Pref     uint64                     `json:"pref,omitempty"`
	Label    string                     `json:"label,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (o OnlineService) MarshalJSON() ([]byte, error) {
	if o.Type == "" {
		o.Type = "OnlineService"
	}
	type alias OnlineService
	return marshalExtra(alias(o), o.Extra)
}

func (o *OnlineService) UnmarshalJSON(b []byte) error {
	type alias OnlineService
	return unmarshalExtra(b, (*alias)(o), &o.Extra)
}

// Phone is a phone number (RFC 9553 §2.3.3).
type Phone struct {
	Type     string                     `json:"@type,omitempty"` // "Phone"
	Number   string                     `json:"number"`
	Features Features                   `json:"features,omitempty"`
	Contexts Contexts                   `json:"contexts,omitempty"`
	Pref     uint64                     `json:"pref,omitempty"`
	Label    string                     `json:"label,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (p Phone) MarshalJSON() ([]byte, error) {
	if p.Type == "" {
		p.Type = "Phone"
	}
	type alias Phone
	return marshalExtra(alias(p), p.Extra)
}

func (p *Phone) UnmarshalJSON(b []byte) error {
	type alias Phone
	return unmarshalExtra(b, (*alias)(p), &p.Extra)
}

// LanguagePref is a preferred language for contacting the entity (RFC 9553
// §2.3.4).
type LanguagePref struct {
	Type     string                     `json:"@type,omitempty"` // "LanguagePref"
	Language string                     `json:"language"`
	Contexts Contexts                   `json:"contexts,omitempty"`
	Pref     uint64                     `json:"pref,omitempty"`
	Extra    map[string]json.RawMessage `json:"-"`
}

func (l LanguagePref) MarshalJSON() ([]byte, error) {
	if l.Type == "" {
		l.Type = "LanguagePref"
	}
	type alias LanguagePref
	return marshalExtra(alias(l), l.Extra)
}

func (l *LanguagePref) UnmarshalJSON(b []byte) error {
	type alias LanguagePref
	return unmarshalExtra(b, (*alias)(l), &l.Extra)
}
