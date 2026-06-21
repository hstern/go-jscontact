// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// Card is the top-level JSContact object: a single contact card
// (RFC 9553 §2). Its fields are declared in the order @type, version, then
// the remaining properties so that the "@type" discriminator marshals first,
// which interop consumers rely on.
//
// Map-valued properties use the JSContact "Id[…]" shape, where each key is a
// caller-assigned Id (RFC 9553 §1.4.1). Unknown members decode into Extra
// and re-marshal verbatim, so a Card round-trips without data loss.
type Card struct {
	Type      string              `json:"@type"`   // MUST be "Card"
	Version   string              `json:"version"` // MUST be set; current spec value "1.0"
	Created   UTCDateTime         `json:"created,omitempty"`
	Kind      string              `json:"kind,omitempty"`
	Language  string              `json:"language,omitempty"`
	Members   map[string]bool     `json:"members,omitempty"`
	ProdID    string              `json:"prodId,omitempty"`
	RelatedTo map[string]Relation `json:"relatedTo,omitempty"`
	UID       string              `json:"uid"` // MUST be set
	Updated   UTCDateTime         `json:"updated,omitempty"`

	Name          *Name                   `json:"name,omitempty"`
	Nicknames     map[string]Nickname     `json:"nicknames,omitempty"`
	Organizations map[string]Organization `json:"organizations,omitempty"`
	SpeakToAs     *SpeakToAs              `json:"speakToAs,omitempty"`
	Titles        map[string]Title        `json:"titles,omitempty"`

	Emails             map[string]EmailAddress  `json:"emails,omitempty"`
	OnlineServices     map[string]OnlineService `json:"onlineServices,omitempty"`
	Phones             map[string]Phone         `json:"phones,omitempty"`
	PreferredLanguages map[string]LanguagePref  `json:"preferredLanguages,omitempty"`

	Calendars           map[string]Calendar          `json:"calendars,omitempty"`
	SchedulingAddresses map[string]SchedulingAddress `json:"schedulingAddresses,omitempty"`
	Addresses           map[string]Address           `json:"addresses,omitempty"`

	CryptoKeys  map[string]CryptoKey `json:"cryptoKeys,omitempty"`
	Directories map[string]Directory `json:"directories,omitempty"`
	Links       map[string]Link      `json:"links,omitempty"`
	Media       map[string]Media     `json:"media,omitempty"`

	Localizations map[string]PatchObject `json:"localizations,omitempty"`

	Anniversaries map[string]Anniversary  `json:"anniversaries,omitempty"`
	Keywords      map[string]bool         `json:"keywords,omitempty"`
	Notes         map[string]Note         `json:"notes,omitempty"`
	PersonalInfo  map[string]PersonalInfo `json:"personalInfo,omitempty"`

	// Extra holds members not defined by RFC 9553 (vendor/IANA extensions).
	// They round-trip verbatim. Use DecodeExtra/EncodeExtra for typed access.
	Extra map[string]json.RawMessage `json:"-"`
}

func (c Card) MarshalJSON() ([]byte, error) {
	if c.Type == "" {
		c.Type = "Card"
	}
	type alias Card
	return marshalExtra(alias(c), c.Extra)
}

func (c *Card) UnmarshalJSON(b []byte) error {
	type alias Card
	return unmarshalExtra(b, (*alias)(c), &c.Extra)
}

// Parse decodes a JSContact Card from its JSON encoding. It is liberal
// (RFC 9553 follows Postel's law): unknown members are preserved in Extra
// rather than rejected. Call Card.Validate for strict checking.
func Parse(b []byte) (*Card, error) {
	var c Card
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
