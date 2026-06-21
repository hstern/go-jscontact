// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import (
	"fmt"
	"sort"
	"strings"
)

// Validation enforces the RFC 9553 §2 "MUST" requirements. Following the
// library's lenient-unmarshal / strict-marshal stance, decoding never fails
// on these; callers opt in by invoking Card.Validate.

// FieldError is a single validation failure located by a dotted/bracketed
// path from the Card root, e.g. `emails["e1"].address`.
type FieldError struct {
	Path    string
	Message string
}

func (e FieldError) String() string { return e.Path + ": " + e.Message }

// ValidationError aggregates every FieldError found in one Validate pass.
type ValidationError struct {
	Errors []FieldError
}

func (e *ValidationError) Error() string {
	switch len(e.Errors) {
	case 0:
		return "jscontact: invalid Card"
	case 1:
		return "jscontact: " + e.Errors[0].String()
	default:
		parts := make([]string, len(e.Errors))
		for i, fe := range e.Errors {
			parts[i] = fe.String()
		}
		return fmt.Sprintf("jscontact: %d validation errors: %s", len(e.Errors), strings.Join(parts, "; "))
	}
}

// validator accumulates field errors during a Validate pass.
type validator struct {
	errs []FieldError
}

func (v *validator) add(path, msg string) { v.errs = append(v.errs, FieldError{path, msg}) }

// id reports whether s is a valid Id (RFC 9553 §1.4.1): 1–255 octets drawn
// from the URL- and filename-safe base64 alphabet (A–Z, a–z, 0–9, '-', '_').
func validID(s string) bool {
	if len(s) == 0 || len(s) > 255 {
		return false
	}
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
		default:
			return false
		}
	}
	return true
}

// checkIDs validates the keys of an Id-keyed map property.
func (v *validator) checkIDs(prop string, keys []string) {
	for _, k := range keys {
		if !validID(k) {
			v.add(fmt.Sprintf("%s[%q]", prop, k), "invalid Id key (RFC 9553 §1.4.1)")
		}
	}
}

// checkPref validates that a pref value, when present, is within 1–100
// (RFC 9553 §1.5.3). A zero value means "absent" and is skipped.
func (v *validator) checkPref(path string, pref uint64) {
	if pref != 0 && pref > 100 {
		v.add(path, "pref must be in the range 1–100")
	}
}

// Validate checks the Card against the RFC 9553 §2 requirements and returns a
// *ValidationError describing every violation, or nil if the Card is valid.
func (c *Card) Validate() error {
	v := &validator{}

	if c.Type != "" && c.Type != "Card" {
		v.add("@type", `must be "Card"`)
	}
	if c.Version == "" {
		v.add("version", "is required (RFC 9553 §2.1.2)")
	}
	if c.UID == "" {
		v.add("uid", "is required (RFC 9553 §2.1.9)")
	}
	// members MUST NOT be set unless kind is "group" (RFC 9553 §2.1.6). Note
	// that members, relatedTo, keywords, and localizations are String[…] maps
	// keyed by UID / free string / language tag — not Id[…] maps — so their
	// keys are not subject to the §1.4.1 Id restrictions.
	if len(c.Members) > 0 && c.Kind != "group" {
		v.add("members", `is only valid when kind is "group"`)
	}

	// Contact channels: required fields + pref + Id keys.
	v.checkIDs("emails", keys(c.Emails))
	for k, e := range c.Emails {
		p := fmt.Sprintf("emails[%q]", k)
		if e.Address == "" {
			v.add(p+".address", "is required")
		}
		v.checkPref(p+".pref", e.Pref)
	}
	v.checkIDs("phones", keys(c.Phones))
	for k, ph := range c.Phones {
		p := fmt.Sprintf("phones[%q]", k)
		if ph.Number == "" {
			v.add(p+".number", "is required")
		}
		v.checkPref(p+".pref", ph.Pref)
	}
	v.checkIDs("onlineServices", keys(c.OnlineServices))
	for k, o := range c.OnlineServices {
		p := fmt.Sprintf("onlineServices[%q]", k)
		if o.URI == "" && o.User == "" {
			v.add(p, "must set at least one of uri or user")
		}
		v.checkPref(p+".pref", o.Pref)
	}
	v.checkIDs("preferredLanguages", keys(c.PreferredLanguages))
	for k, l := range c.PreferredLanguages {
		p := fmt.Sprintf("preferredLanguages[%q]", k)
		if l.Language == "" {
			v.add(p+".language", "is required")
		}
		v.checkPref(p+".pref", l.Pref)
	}

	// Addresses: components require value+kind; pref range.
	v.checkIDs("addresses", keys(c.Addresses))
	for k, a := range c.Addresses {
		p := fmt.Sprintf("addresses[%q]", k)
		for i, comp := range a.Components {
			cp := fmt.Sprintf("%s.components[%d]", p, i)
			if comp.Kind == "" {
				v.add(cp+".kind", "is required")
			}
			if comp.Value == "" {
				v.add(cp+".value", "is required")
			}
		}
		v.checkPref(p+".pref", a.Pref)
	}

	// Name components require value+kind.
	if c.Name != nil {
		for i, comp := range c.Name.Components {
			cp := fmt.Sprintf("name.components[%d]", i)
			if comp.Kind == "" {
				v.add(cp+".kind", "is required")
			}
			if comp.Value == "" {
				v.add(cp+".value", "is required")
			}
		}
	}

	// Resources: uri required (and Media.kind), pref range.
	v.checkIDs("media", keys(c.Media))
	for k, m := range c.Media {
		p := fmt.Sprintf("media[%q]", k)
		if m.URI == "" {
			v.add(p+".uri", "is required")
		}
		if m.Kind == "" {
			v.add(p+".kind", "is required")
		}
		v.checkPref(p+".pref", m.Pref)
	}
	v.checkIDs("links", keys(c.Links))
	v.checkIDs("cryptoKeys", keys(c.CryptoKeys))
	v.checkIDs("directories", keys(c.Directories))
	v.checkIDs("calendars", keys(c.Calendars))
	v.checkIDs("schedulingAddresses", keys(c.SchedulingAddresses))

	// Anniversaries: kind + a date value required.
	v.checkIDs("anniversaries", keys(c.Anniversaries))
	for k, a := range c.Anniversaries {
		p := fmt.Sprintf("anniversaries[%q]", k)
		if a.Kind == "" {
			v.add(p+".kind", "is required")
		}
		if a.Date.PartialDate == nil && a.Date.Timestamp == nil {
			v.add(p+".date", "is required")
		}
	}

	// Other Id-keyed properties (keywords is excluded: its keys are arbitrary
	// keyword strings, not Ids).
	v.checkIDs("nicknames", keys(c.Nicknames))
	v.checkIDs("organizations", keys(c.Organizations))
	v.checkIDs("titles", keys(c.Titles))
	v.checkIDs("notes", keys(c.Notes))
	for k, n := range c.Notes {
		if n.Note == "" {
			v.add(fmt.Sprintf("notes[%q].note", k), "is required")
		}
	}
	v.checkIDs("personalInfo", keys(c.PersonalInfo))
	for k, pi := range c.PersonalInfo {
		p := fmt.Sprintf("personalInfo[%q]", k)
		if pi.Kind == "" {
			v.add(p+".kind", "is required")
		}
		if pi.Value == "" {
			v.add(p+".value", "is required")
		}
	}

	if len(v.errs) == 0 {
		return nil
	}
	sort.Slice(v.errs, func(i, j int) bool { return v.errs[i].Path < v.errs[j].Path })
	return &ValidationError{Errors: v.errs}
}

// keys returns the keys of a map of any value type, for Id validation.
func keys[V any](m map[string]V) []string {
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
