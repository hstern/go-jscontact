// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package vcard

import (
	"sort"
	"strconv"
	"strings"

	vc "github.com/emersion/go-vcard"
	jscontact "github.com/hstern/go-jscontact"
)

// ToVCard converts a JSContact Card into a vCard (RFC 6350) per the RFC 9555
// mapping. The result targets vCard 4.0. Conversion is asymmetric where the
// formats diverge: JSContact features without a vCard equivalent are dropped
// rather than smuggled into nonstandard properties.
func ToVCard(c *jscontact.Card) (vc.Card, error) {
	card := make(vc.Card)
	card.SetValue(vc.FieldVersion, "4.0")

	if c.UID != "" {
		card.SetValue(vc.FieldUID, c.UID)
	}
	if c.Kind != "" {
		card.SetKind(vc.Kind(c.Kind))
	}
	if c.ProdID != "" {
		card.SetValue(vc.FieldProductID, c.ProdID)
	}
	if c.Updated != "" {
		card.SetValue(vc.FieldRevision, string(c.Updated))
	}

	nameToVCard(c, card)
	nicknamesToVCard(c, card)
	orgsToVCard(c, card)
	titlesToVCard(c, card)
	emailsToVCard(c, card)
	phonesToVCard(c, card)
	onlineToVCard(c, card)
	addressesToVCard(c, card)
	resourcesToVCard(c, card)
	notesToVCard(c, card)
	keywordsToVCard(c, card)
	anniversariesToVCard(c, card)
	languagesToVCard(c, card)
	membersToVCard(c, card)

	return card, nil
}

func nameToVCard(c *jscontact.Card, card vc.Card) {
	if c.Name == nil {
		return
	}
	if c.Name.Full != "" {
		card.SetValue(vc.FieldFormattedName, c.Name.Full)
	}
	if len(c.Name.Components) == 0 {
		return
	}
	n := &vc.Name{}
	for _, comp := range c.Name.Components {
		switch comp.Kind {
		case "surname":
			n.FamilyName = join(n.FamilyName, comp.Value)
		case "given":
			n.GivenName = join(n.GivenName, comp.Value)
		case "given2":
			n.AdditionalName = join(n.AdditionalName, comp.Value)
		case "title":
			n.HonorificPrefix = join(n.HonorificPrefix, comp.Value)
		case "credential", "generation":
			n.HonorificSuffix = join(n.HonorificSuffix, comp.Value)
		}
	}
	card.AddName(n)
}

func nicknamesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Nicknames) {
		nn := c.Nicknames[k]
		card.Add(vc.FieldNickname, &vc.Field{Value: nn.Name, Params: paramsFor(nn.Contexts, nn.Pref)})
	}
}

func orgsToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Organizations) {
		org := c.Organizations[k]
		parts := []string{org.Name}
		for _, u := range org.Units {
			parts = append(parts, u.Name)
		}
		card.Add(vc.FieldOrganization, &vc.Field{Value: strings.Join(parts, ";")})
	}
}

func titlesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Titles) {
		t := c.Titles[k]
		field := vc.FieldTitle
		if t.Kind == "role" {
			field = vc.FieldRole
		}
		card.Add(field, &vc.Field{Value: t.Name})
	}
}

func emailsToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Emails) {
		e := c.Emails[k]
		card.Add(vc.FieldEmail, &vc.Field{Value: e.Address, Params: paramsFor(e.Contexts, e.Pref)})
	}
}

func phonesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Phones) {
		ph := c.Phones[k]
		p := paramsFor(ph.Contexts, ph.Pref)
		for _, feat := range sortedSet(ph.Features) {
			if t, ok := typeFromPhoneFeature[feat]; ok {
				p.Add(vc.ParamType, t)
			}
		}
		card.Add(vc.FieldTelephone, &vc.Field{Value: ph.Number, Params: p})
	}
}

func onlineToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.OnlineServices) {
		o := c.OnlineServices[k]
		uri := o.URI
		if uri == "" {
			uri = o.User
		}
		if uri != "" {
			card.Add(vc.FieldIMPP, &vc.Field{Value: uri, Params: paramsFor(o.Contexts, o.Pref)})
		}
	}
}

func addressesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Addresses) {
		a := c.Addresses[k]
		adr := &vc.Address{}
		for _, comp := range a.Components {
			switch comp.Kind {
			case "postOfficeBox":
				adr.PostOfficeBox = join(adr.PostOfficeBox, comp.Value)
			case "apartment", "floor", "building", "room":
				adr.ExtendedAddress = join(adr.ExtendedAddress, comp.Value)
			case "name", "number":
				adr.StreetAddress = join(adr.StreetAddress, comp.Value)
			case "locality":
				adr.Locality = join(adr.Locality, comp.Value)
			case "region":
				adr.Region = join(adr.Region, comp.Value)
			case "postcode":
				adr.PostalCode = join(adr.PostalCode, comp.Value)
			case "country":
				adr.Country = join(adr.Country, comp.Value)
			}
		}
		adr.Field = &vc.Field{Params: paramsFor(a.Contexts, a.Pref)}
		card.AddAddress(adr)
	}
}

func resourcesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Media) {
		m := c.Media[k]
		field := vc.FieldPhoto
		switch m.Kind {
		case "logo":
			field = vc.FieldLogo
		case "sound":
			field = vc.FieldSound
		}
		card.Add(field, mediaField(m.URI, m.MediaType))
	}
	for _, k := range sortedKeys(c.Links) {
		card.Add(vc.FieldURL, &vc.Field{Value: c.Links[k].URI, Params: paramsFor(c.Links[k].Contexts, c.Links[k].Pref)})
	}
	for _, k := range sortedKeys(c.CryptoKeys) {
		card.Add(vc.FieldKey, mediaField(c.CryptoKeys[k].URI, c.CryptoKeys[k].MediaType))
	}
}

func notesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Notes) {
		card.Add(vc.FieldNote, &vc.Field{Value: c.Notes[k].Note})
	}
}

func keywordsToVCard(c *jscontact.Card, card vc.Card) {
	if cats := sortedSet(c.Keywords); len(cats) > 0 {
		card.SetCategories(cats)
	}
}

func anniversariesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.Anniversaries) {
		a := c.Anniversaries[k]
		v := dateString(a.Date)
		if v == "" {
			continue
		}
		switch a.Kind {
		case "birth":
			card.SetValue(vc.FieldBirthday, v)
		case "wedding":
			card.SetValue(vc.FieldAnniversary, v)
		}
	}
}

func languagesToVCard(c *jscontact.Card, card vc.Card) {
	for _, k := range sortedKeys(c.PreferredLanguages) {
		l := c.PreferredLanguages[k]
		p := make(vc.Params)
		if l.Pref != 0 {
			p.Set(vc.ParamPreferred, strconv.FormatUint(l.Pref, 10))
		}
		card.Add(vc.FieldLanguage, &vc.Field{Value: l.Language, Params: p})
	}
}

func membersToVCard(c *jscontact.Card, card vc.Card) {
	for _, m := range sortedSet(c.Members) {
		card.Add(vc.FieldMember, &vc.Field{Value: m})
	}
}

func mediaField(uri, mediaType string) *vc.Field {
	f := &vc.Field{Value: uri, Params: make(vc.Params)}
	if mediaType != "" {
		f.Params.Set(vc.ParamMediaType, mediaType)
	}
	return f
}

// dateString renders an AnniversaryDate as a vCard date value.
func dateString(d jscontact.AnniversaryDate) string {
	if d.Timestamp != nil {
		return string(d.Timestamp.UTC)
	}
	if d.PartialDate == nil {
		return ""
	}
	p := d.PartialDate
	switch {
	case p.Year != 0 && p.Month != 0 && p.Day != 0:
		return pad4(p.Year) + pad2(p.Month) + pad2(p.Day)
	case p.Year != 0:
		return pad4(p.Year)
	case p.Month != 0 && p.Day != 0:
		return "--" + pad2(p.Month) + pad2(p.Day)
	default:
		return ""
	}
}

func pad4(n uint64) string { return zero(strconv.FormatUint(n, 10), 4) }
func pad2(n uint64) string { return zero(strconv.FormatUint(n, 10), 2) }

func zero(s string, width int) string {
	for len(s) < width {
		s = "0" + s
	}
	return s
}

func join(existing, add string) string {
	if existing == "" {
		return add
	}
	return existing + " " + add
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedSet(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k, v := range m {
		if v {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}
