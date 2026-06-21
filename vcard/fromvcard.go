// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package vcard

import (
	"strconv"
	"strings"

	vc "github.com/emersion/go-vcard"
	jscontact "github.com/hstern/go-jscontact"
)

// FromVCard converts a parsed vCard (RFC 6350) into a JSContact Card per the
// RFC 9555 mapping. Map-valued JSContact properties are keyed with stable,
// index-derived Ids since vCard has no equivalent of JSContact's Ids.
func FromVCard(card vc.Card) (*jscontact.Card, error) {
	c := &jscontact.Card{Type: "Card", Version: "1.0"}

	if uid := card.Value(vc.FieldUID); uid != "" {
		c.UID = uid
	}
	if k := string(card.Kind()); k != "" {
		c.Kind = k
	}
	if pid := card.Value(vc.FieldProductID); pid != "" {
		c.ProdID = pid
	}
	if rev := card.Value(vc.FieldRevision); rev != "" {
		c.Updated = jscontact.UTCDateTime(rev)
	}

	convertName(card, c)
	convertNicknames(card, c)
	convertOrganizations(card, c)
	convertTitles(card, c)
	convertEmails(card, c)
	convertPhones(card, c)
	convertOnlineServices(card, c)
	convertAddresses(card, c)
	convertResources(card, c)
	convertNotes(card, c)
	convertKeywords(card, c)
	convertAnniversaries(card, c)
	convertLanguages(card, c)
	convertMembers(card, c)

	return c, nil
}

func convertName(card vc.Card, c *jscontact.Card) {
	n := card.Name()
	fn := card.Value(vc.FieldFormattedName)
	if n == nil && fn == "" {
		return
	}
	name := &jscontact.Name{}
	if fn != "" {
		name.Full = fn
	}
	if n != nil {
		// Conventional reading order: prefix, given, additional, family, suffix.
		add := func(kind, val string) {
			if val != "" {
				name.Components = append(name.Components, jscontact.NameComponent{Kind: kind, Value: val})
			}
		}
		add("title", n.HonorificPrefix)
		add("given", n.GivenName)
		add("given2", n.AdditionalName)
		add("surname", n.FamilyName)
		add("credential", n.HonorificSuffix)
	}
	c.Name = name
}

func convertNicknames(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldNickname] {
		if f.Value == "" {
			continue
		}
		put(&c.Nicknames, idFor("nick", i), jscontact.Nickname{
			Name:     f.Value,
			Contexts: contextsFromParams(f.Params),
			Pref:     prefFromParams(f.Params),
		})
	}
}

func convertOrganizations(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldOrganization] {
		if f.Value == "" {
			continue
		}
		parts := strings.Split(f.Value, ";")
		org := jscontact.Organization{Name: parts[0]}
		for _, u := range parts[1:] {
			if u != "" {
				org.Units = append(org.Units, jscontact.OrgUnit{Name: u})
			}
		}
		put(&c.Organizations, idFor("org", i), org)
	}
}

func convertTitles(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldTitle] {
		if f.Value != "" {
			put(&c.Titles, idFor("title", i), jscontact.Title{Name: f.Value, Kind: "title"})
		}
	}
	for i, f := range card[vc.FieldRole] {
		if f.Value != "" {
			put(&c.Titles, idFor("role", i), jscontact.Title{Name: f.Value, Kind: "role"})
		}
	}
}

func convertEmails(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldEmail] {
		if f.Value == "" {
			continue
		}
		put(&c.Emails, idFor("e", i), jscontact.EmailAddress{
			Address:  f.Value,
			Contexts: contextsFromParams(f.Params),
			Pref:     prefFromParams(f.Params),
		})
	}
}

func convertPhones(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldTelephone] {
		if f.Value == "" {
			continue
		}
		ph := jscontact.Phone{Number: f.Value, Pref: prefFromParams(f.Params)}
		for _, t := range f.Params.Types() {
			lt := strings.ToLower(t)
			switch lt {
			case vc.TypeHome:
				ph.Contexts = setTrue(ph.Contexts, "private")
			case vc.TypeWork:
				ph.Contexts = setTrue(ph.Contexts, "work")
			default:
				if feat, ok := phoneFeatureFromType[lt]; ok {
					ph.Features = setFeature(ph.Features, feat)
				}
			}
		}
		put(&c.Phones, idFor("p", i), ph)
	}
}

func convertOnlineServices(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldIMPP] {
		if f.Value == "" {
			continue
		}
		put(&c.OnlineServices, idFor("svc", i), jscontact.OnlineService{
			URI:      f.Value,
			Contexts: contextsFromParams(f.Params),
			Pref:     prefFromParams(f.Params),
		})
	}
}

func convertAddresses(card vc.Card, c *jscontact.Card) {
	for i, a := range card.Addresses() {
		addr := jscontact.Address{Contexts: contextsFromParams(a.Params), Pref: prefFromParams(a.Params)}
		add := func(kind, val string) {
			if val != "" {
				addr.Components = append(addr.Components, jscontact.AddressComponent{Kind: kind, Value: val})
			}
		}
		add("postOfficeBox", a.PostOfficeBox)
		add("apartment", a.ExtendedAddress)
		add("name", a.StreetAddress)
		add("locality", a.Locality)
		add("region", a.Region)
		add("postcode", a.PostalCode)
		add("country", a.Country)
		// Note: vCard ADR carries a country name, not the ISO 3166-1 alpha-2
		// code JSContact's countryCode wants, so countryCode is left unset.
		put(&c.Addresses, idFor("adr", i), addr)
	}
}

func convertResources(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldPhoto] {
		put(&c.Media, idFor("photo", i), jscontact.Media{Kind: "photo", URI: f.Value, MediaType: f.Params.Get(vc.ParamMediaType)})
	}
	for i, f := range card[vc.FieldLogo] {
		put(&c.Media, idFor("logo", i), jscontact.Media{Kind: "logo", URI: f.Value, MediaType: f.Params.Get(vc.ParamMediaType)})
	}
	for i, f := range card[vc.FieldSound] {
		put(&c.Media, idFor("sound", i), jscontact.Media{Kind: "sound", URI: f.Value, MediaType: f.Params.Get(vc.ParamMediaType)})
	}
	for i, f := range card[vc.FieldURL] {
		put(&c.Links, idFor("link", i), jscontact.Link{URI: f.Value, Pref: prefFromParams(f.Params)})
	}
	for i, f := range card[vc.FieldKey] {
		put(&c.CryptoKeys, idFor("key", i), jscontact.CryptoKey{URI: f.Value, MediaType: f.Params.Get(vc.ParamMediaType)})
	}
}

func convertNotes(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldNote] {
		if f.Value != "" {
			put(&c.Notes, idFor("note", i), jscontact.Note{Note: f.Value})
		}
	}
}

func convertKeywords(card vc.Card, c *jscontact.Card) {
	for _, cat := range card.Categories() {
		if cat != "" {
			if c.Keywords == nil {
				c.Keywords = make(map[string]bool)
			}
			c.Keywords[cat] = true
		}
	}
}

func convertAnniversaries(card vc.Card, c *jscontact.Card) {
	if b := card.Value(vc.FieldBirthday); b != "" {
		put(&c.Anniversaries, "bday", jscontact.Anniversary{Kind: "birth", Date: anniversaryDate(b)})
	}
	if a := card.Value(vc.FieldAnniversary); a != "" {
		put(&c.Anniversaries, "anniversary", jscontact.Anniversary{Kind: "wedding", Date: anniversaryDate(a)})
	}
}

func convertLanguages(card vc.Card, c *jscontact.Card) {
	for i, f := range card[vc.FieldLanguage] {
		if f.Value != "" {
			put(&c.PreferredLanguages, idFor("lang", i), jscontact.LanguagePref{Language: f.Value, Pref: prefFromParams(f.Params)})
		}
	}
}

func convertMembers(card vc.Card, c *jscontact.Card) {
	for _, f := range card[vc.FieldMember] {
		if f.Value != "" {
			if c.Members == nil {
				c.Members = make(map[string]bool)
			}
			c.Members[f.Value] = true
		}
	}
}

// anniversaryDate parses a vCard date into a JSContact AnniversaryDate. A
// full RFC 3339 timestamp becomes a Timestamp; a calendar date (or partial
// date like "1970" or "--05") becomes a PartialDate.
func anniversaryDate(v string) jscontact.AnniversaryDate {
	if strings.Contains(v, "T") {
		return jscontact.AnniversaryDate{Timestamp: &jscontact.Timestamp{UTC: jscontact.UTCDateTime(v)}}
	}
	pd := &jscontact.PartialDate{}
	// Accept YYYY-MM-DD, YYYYMMDD, YYYY, and the vCard "--MM-DD" partial form.
	digits := strings.ReplaceAll(strings.TrimLeft(v, "-"), "-", "")
	switch {
	case strings.HasPrefix(v, "--") && len(digits) >= 2:
		pd.Month = atoiOr(digits[0:2])
		if len(digits) >= 4 {
			pd.Day = atoiOr(digits[2:4])
		}
	default:
		if len(digits) >= 4 {
			pd.Year = atoiOr(digits[0:4])
		}
		if len(digits) >= 6 {
			pd.Month = atoiOr(digits[4:6])
		}
		if len(digits) >= 8 {
			pd.Day = atoiOr(digits[6:8])
		}
	}
	return jscontact.AnniversaryDate{PartialDate: pd}
}

func atoiOr(s string) uint64 {
	n, _ := strconv.ParseUint(s, 10, 64)
	return n
}

func idFor(prefix string, i int) string { return prefix + strconv.Itoa(i+1) }

// put inserts v into the map pointed to by m, allocating it if needed.
func put[V any](m *map[string]V, key string, v V) {
	if *m == nil {
		*m = make(map[string]V)
	}
	(*m)[key] = v
}
