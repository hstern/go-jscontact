// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// Resource-style property objects (RFC 9553 §2.6). Each shares the common
// Resource members (uri, mediaType, contexts, pref, label) plus its own
// kind vocabulary.

// CryptoKey is a public key or certificate (RFC 9553 §2.6.1).
type CryptoKey struct {
	Type      string                     `json:"@type,omitempty"` // "CryptoKey"
	URI       string                     `json:"uri"`
	Kind      string                     `json:"kind,omitempty"`
	MediaType string                     `json:"mediaType,omitempty"`
	Contexts  Contexts                   `json:"contexts,omitempty"`
	Pref      uint64                     `json:"pref,omitempty"`
	Label     string                     `json:"label,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

func (c CryptoKey) MarshalJSON() ([]byte, error) {
	if c.Type == "" {
		c.Type = "CryptoKey"
	}
	type alias CryptoKey
	return marshalExtra(alias(c), c.Extra)
}

func (c *CryptoKey) UnmarshalJSON(b []byte) error {
	type alias CryptoKey
	return unmarshalExtra(b, (*alias)(c), &c.Extra)
}

// Directory is a link to the entity's entry in a directory (RFC 9553
// §2.6.2). ListAs gives the position within a directory listing.
type Directory struct {
	Type      string                     `json:"@type,omitempty"` // "Directory"
	Kind      string                     `json:"kind,omitempty"`  // "directory" | "entry"
	URI       string                     `json:"uri"`
	MediaType string                     `json:"mediaType,omitempty"`
	ListAs    uint64                     `json:"listAs,omitempty"`
	Contexts  Contexts                   `json:"contexts,omitempty"`
	Pref      uint64                     `json:"pref,omitempty"`
	Label     string                     `json:"label,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

func (d Directory) MarshalJSON() ([]byte, error) {
	if d.Type == "" {
		d.Type = "Directory"
	}
	type alias Directory
	return marshalExtra(alias(d), d.Extra)
}

func (d *Directory) UnmarshalJSON(b []byte) error {
	type alias Directory
	return unmarshalExtra(b, (*alias)(d), &d.Extra)
}

// Link is a generic link associated with the entity (RFC 9553 §2.6.3).
type Link struct {
	Type      string                     `json:"@type,omitempty"` // "Link"
	Kind      string                     `json:"kind,omitempty"`  // e.g. "contact"
	URI       string                     `json:"uri"`
	MediaType string                     `json:"mediaType,omitempty"`
	Contexts  Contexts                   `json:"contexts,omitempty"`
	Pref      uint64                     `json:"pref,omitempty"`
	Label     string                     `json:"label,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

func (l Link) MarshalJSON() ([]byte, error) {
	if l.Type == "" {
		l.Type = "Link"
	}
	type alias Link
	return marshalExtra(alias(l), l.Extra)
}

func (l *Link) UnmarshalJSON(b []byte) error {
	type alias Link
	return unmarshalExtra(b, (*alias)(l), &l.Extra)
}

// Media is a photo, sound, or logo for the entity (RFC 9553 §2.6.4).
type Media struct {
	Type      string                     `json:"@type,omitempty"` // "Media"
	Kind      string                     `json:"kind"`            // "photo" | "sound" | "logo"
	URI       string                     `json:"uri"`
	MediaType string                     `json:"mediaType,omitempty"`
	Contexts  Contexts                   `json:"contexts,omitempty"`
	Pref      uint64                     `json:"pref,omitempty"`
	Label     string                     `json:"label,omitempty"`
	Extra     map[string]json.RawMessage `json:"-"`
}

func (m Media) MarshalJSON() ([]byte, error) {
	if m.Type == "" {
		m.Type = "Media"
	}
	type alias Media
	return marshalExtra(alias(m), m.Extra)
}

func (m *Media) UnmarshalJSON(b []byte) error {
	type alias Media
	return unmarshalExtra(b, (*alias)(m), &m.Extra)
}
