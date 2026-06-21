// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package vcard

import (
	"bytes"
	"strings"
	"testing"

	vc "github.com/emersion/go-vcard"
	jscontact "github.com/hstern/go-jscontact"
)

func decodeVCard(t *testing.T, s string) vc.Card {
	t.Helper()
	card, err := vc.NewDecoder(strings.NewReader(s)).Decode()
	if err != nil {
		t.Fatalf("decode vCard: %v", err)
	}
	return card
}

func encodeVCard(t *testing.T, card vc.Card) string {
	t.Helper()
	var buf bytes.Buffer
	if err := vc.NewEncoder(&buf).Encode(card); err != nil {
		t.Fatalf("encode vCard: %v", err)
	}
	return buf.String()
}

const sampleVCard = "BEGIN:VCARD\r\n" +
	"VERSION:4.0\r\n" +
	"UID:urn:uuid:abc\r\n" +
	"FN:John Doe\r\n" +
	"N:Doe;John;;Dr.;\r\n" +
	"ORG:Example Inc.;Engineering\r\n" +
	"TITLE:Staff Engineer\r\n" +
	"EMAIL;TYPE=work;PREF=1:jdoe@example.com\r\n" +
	"TEL;TYPE=\"work,cell\":+1-555-0100\r\n" +
	"ADR;TYPE=work:;;1 Main St;Zurich;;8000;CH\r\n" +
	"NOTE:hello\r\n" +
	"CATEGORIES:friend,colleague\r\n" +
	"BDAY:19700521\r\n" +
	"END:VCARD\r\n"

func TestFromVCard(t *testing.T) {
	c, err := FromVCard(decodeVCard(t, sampleVCard))
	if err != nil {
		t.Fatal(err)
	}
	if c.UID != "urn:uuid:abc" {
		t.Errorf("uid = %q", c.UID)
	}
	if c.Name == nil || c.Name.Full != "John Doe" {
		t.Fatalf("name.full = %+v", c.Name)
	}
	// N components in reading order: title(Dr.), given(John), surname(Doe).
	kinds := componentKinds(c.Name.Components)
	if strings.Join(kinds, ",") != "title,given,surname" {
		t.Errorf("name components = %v", kinds)
	}
	if len(c.Emails) != 1 {
		t.Fatalf("emails = %+v", c.Emails)
	}
	for _, e := range c.Emails {
		if e.Address != "jdoe@example.com" || !e.Contexts["work"] || e.Pref != 1 {
			t.Errorf("email mapped wrong: %+v", e)
		}
	}
	for _, ph := range c.Phones {
		if !ph.Contexts["work"] || !ph.Features["mobile"] {
			t.Errorf("phone TYPE not split into context+feature: %+v", ph)
		}
	}
	if !c.Keywords["friend"] || !c.Keywords["colleague"] {
		t.Errorf("categories→keywords lost: %+v", c.Keywords)
	}
	bday := c.Anniversaries["bday"]
	if bday.Kind != "birth" || bday.Date.PartialDate == nil || bday.Date.PartialDate.Year != 1970 || bday.Date.PartialDate.Month != 5 || bday.Date.PartialDate.Day != 21 {
		t.Errorf("bday mapped wrong: %+v", bday)
	}
	if err := c.Validate(); err != nil {
		t.Errorf("converted card should validate: %v", err)
	}
}

func TestToVCard(t *testing.T) {
	c := &jscontact.Card{
		Type: "Card", Version: "1.0", UID: "x", Kind: "individual",
		Name: &jscontact.Name{Full: "Jane Roe", Components: []jscontact.NameComponent{
			{Kind: "given", Value: "Jane"}, {Kind: "surname", Value: "Roe"},
		}},
		Emails: map[string]jscontact.EmailAddress{"e1": {Address: "jane@example.com", Contexts: jscontact.Contexts{"work": true}, Pref: 1}},
		Phones: map[string]jscontact.Phone{"p1": {Number: "+1-555-0199", Features: jscontact.Features{"mobile": true}}},
	}
	card, err := ToVCard(c)
	if err != nil {
		t.Fatal(err)
	}
	out := encodeVCard(t, card)
	for _, want := range []string{"FN:Jane Roe", "UID:x", "KIND:individual", "jane@example.com", "+1-555-0199"} {
		if !strings.Contains(out, want) {
			t.Errorf("vCard output missing %q\n%s", want, out)
		}
	}
}

func TestRoundTripCoreFields(t *testing.T) {
	orig := &jscontact.Card{
		Type: "Card", Version: "1.0", UID: "round-1", Kind: "individual",
		Name: &jscontact.Name{Components: []jscontact.NameComponent{
			{Kind: "given", Value: "Ada"}, {Kind: "surname", Value: "Lovelace"},
		}},
		Emails:        map[string]jscontact.EmailAddress{"e1": {Address: "ada@example.com", Contexts: jscontact.Contexts{"work": true}}},
		Organizations: map[string]jscontact.Organization{"o1": {Name: "Analytical", Units: []jscontact.OrgUnit{{Name: "Engines"}}}},
		Keywords:      map[string]bool{"math": true},
	}
	card, err := ToVCard(orig)
	if err != nil {
		t.Fatal(err)
	}
	back, err := FromVCard(card)
	if err != nil {
		t.Fatal(err)
	}
	if back.UID != "round-1" || back.Kind != "individual" {
		t.Errorf("identity fields lost: uid=%q kind=%q", back.UID, back.Kind)
	}
	if got := componentKinds(back.Name.Components); strings.Join(got, ",") != "given,surname" {
		t.Errorf("name round-trip changed order/kinds: %v", got)
	}
	var gotEmail jscontact.EmailAddress
	for _, e := range back.Emails {
		gotEmail = e
	}
	if gotEmail.Address != "ada@example.com" || !gotEmail.Contexts["work"] {
		t.Errorf("email round-trip lost data: %+v", gotEmail)
	}
	var gotOrg jscontact.Organization
	for _, o := range back.Organizations {
		gotOrg = o
	}
	if gotOrg.Name != "Analytical" || len(gotOrg.Units) != 1 || gotOrg.Units[0].Name != "Engines" {
		t.Errorf("org round-trip lost data: %+v", gotOrg)
	}
	if !back.Keywords["math"] {
		t.Errorf("keywords round-trip lost: %+v", back.Keywords)
	}
}

func componentKinds(comps []jscontact.NameComponent) []string {
	out := make([]string, len(comps))
	for i, c := range comps {
		out[i] = c.Kind
	}
	return out
}
