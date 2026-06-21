// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import (
	"bytes"
	"encoding/json"
	"testing"
)

// roundTrip parses input, re-marshals it, and asserts the bytes are stable.
// JSContact objects are byte-stable modulo canonical member ordering, so the
// input fixtures are written in the canonical order the codec emits: "@type"
// first, then struct-declaration order, then extension members sorted by key.
func roundTrip(t *testing.T, name, input string) *Card {
	t.Helper()
	c, err := Parse([]byte(input))
	if err != nil {
		t.Fatalf("%s: Parse: %v", name, err)
	}
	got, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("%s: Marshal: %v", name, err)
	}
	if !bytes.Equal(got, []byte(input)) {
		t.Fatalf("%s: round-trip not byte-stable\n in: %s\nout: %s", name, input, got)
	}
	return c
}

func TestRoundTripMinimal(t *testing.T) {
	roundTrip(t, "minimal", `{"@type":"Card","version":"1.0","uid":"abc"}`)
}

func TestRoundTripTypeFirst(t *testing.T) {
	// "@type" must emit first even when the input lists it last.
	in := `{"version":"1.0","uid":"x","@type":"Card"}`
	c, err := Parse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"@type":"Card","version":"1.0","uid":"x"}`
	if string(got) != want {
		t.Fatalf("@type not first:\ngot:  %s\nwant: %s", got, want)
	}
}

func TestRoundTripFullCard(t *testing.T) {
	// A card exercising every property group, in canonical emission order.
	in := `{"@type":"Card","version":"1.0","kind":"individual","uid":"22B2C7DF",` +
		`"name":{"@type":"Name","components":[{"@type":"NameComponent","kind":"given","value":"John"},{"@type":"NameComponent","kind":"surname","value":"Doe"}]},` +
		`"organizations":{"o1":{"@type":"Organization","name":"Example","units":[{"@type":"OrgUnit","name":"Eng"}]}},` +
		`"titles":{"t1":{"@type":"Title","name":"Engineer"}},` +
		`"emails":{"e1":{"@type":"EmailAddress","address":"jdoe@example.com","contexts":{"work":true}}},` +
		`"phones":{"p1":{"@type":"Phone","number":"tel:+1-555-0100","features":{"voice":true}}},` +
		`"addresses":{"a1":{"@type":"Address","components":[{"@type":"AddressComponent","kind":"locality","value":"Zürich"}],"countryCode":"CH"}},` +
		`"anniversaries":{"k1":{"@type":"Anniversary","kind":"birth","date":{"@type":"PartialDate","year":1970,"month":5}}},` +
		`"keywords":{"friend":true},` +
		`"notes":{"n1":{"@type":"Note","note":"met at conf"}}}`
	c := roundTrip(t, "full", in)

	if c.Kind != "individual" {
		t.Errorf("kind = %q", c.Kind)
	}
	if c.Name == nil || len(c.Name.Components) != 2 || c.Name.Components[0].Value != "John" {
		t.Errorf("name components not decoded: %+v", c.Name)
	}
	if !c.Emails["e1"].Contexts["work"] {
		t.Errorf("email work context lost")
	}
	an := c.Anniversaries["k1"]
	if an.Date.PartialDate == nil || an.Date.PartialDate.Year != 1970 || an.Date.PartialDate.Month != 5 {
		t.Errorf("anniversary partial date wrong: %+v", an.Date)
	}
}

func TestComponentOrderPreserved(t *testing.T) {
	// Ordered components must never be reordered.
	in := `{"@type":"Card","version":"1.0","uid":"x","name":{"@type":"Name","components":[` +
		`{"@type":"NameComponent","kind":"surname","value":"Doe"},` +
		`{"@type":"NameComponent","kind":"given","value":"John"}]}}`
	c := roundTrip(t, "order", in)
	if c.Name.Components[0].Kind != "surname" || c.Name.Components[1].Kind != "given" {
		t.Fatalf("component order changed: %+v", c.Name.Components)
	}
}

func TestExtraPassthrough(t *testing.T) {
	// Unknown members survive a decode→encode cycle, both on the Card and on
	// a nested object, and sort by key on output.
	// Extension members emit sorted by key, so the fixture lists them sorted.
	in := `{"@type":"Card","version":"1.0","uid":"x",` +
		`"emails":{"e1":{"@type":"EmailAddress","address":"a@b.c","x-verified":true}},` +
		`"x-aaa":[1,2],"x-source":"crm","x-zzz":1}`
	c := roundTrip(t, "extra", in)

	if string(c.Extra["x-source"]) != `"crm"` {
		t.Errorf("card extra x-source lost: %q", c.Extra["x-source"])
	}
	if _, ok := c.Emails["e1"].Extra["x-verified"]; !ok {
		t.Errorf("nested extra x-verified lost")
	}
}

func TestExtraTypedHelpers(t *testing.T) {
	c, err := Parse([]byte(`{"@type":"Card","version":"1.0","uid":"x","x-rank":{"score":42}}`))
	if err != nil {
		t.Fatal(err)
	}
	var rank struct {
		Score int `json:"score"`
	}
	ok, err := DecodeExtra(c.Extra, "x-rank", &rank)
	if err != nil || !ok {
		t.Fatalf("DecodeExtra: ok=%v err=%v", ok, err)
	}
	if rank.Score != 42 {
		t.Errorf("score = %d", rank.Score)
	}
	if _, err := EncodeExtra(c.Extra, "x-new", map[string]int{"n": 1}); err != nil {
		t.Fatalf("EncodeExtra: %v", err)
	}
	if string(c.Extra["x-new"]) != `{"n":1}` {
		t.Errorf("EncodeExtra stored %q", c.Extra["x-new"])
	}
}

func TestAnniversaryTimestamp(t *testing.T) {
	// A Timestamp-valued anniversary date (vs the default PartialDate).
	in := `{"@type":"Card","version":"1.0","uid":"x",` +
		`"anniversaries":{"k1":{"@type":"Anniversary","kind":"death","date":{"@type":"Timestamp","utc":"1995-02-01T00:00:00Z"}}}}`
	c := roundTrip(t, "timestamp", in)
	d := c.Anniversaries["k1"].Date
	if d.Timestamp == nil || d.Timestamp.UTC != "1995-02-01T00:00:00Z" {
		t.Fatalf("timestamp date wrong: %+v", d)
	}
	if d.PartialDate != nil {
		t.Errorf("partial date should be nil for a Timestamp")
	}
}

func TestSetMembersStayTrue(t *testing.T) {
	c := roundTrip(t, "sets", `{"@type":"Card","version":"1.0","uid":"x","keywords":{"a":true,"b":true}}`)
	if !c.Keywords["a"] || !c.Keywords["b"] {
		t.Fatalf("keyword set wrong: %+v", c.Keywords)
	}
}

func TestParseInvalidJSON(t *testing.T) {
	if _, err := Parse([]byte(`{not json`)); err == nil {
		t.Fatal("expected error on malformed JSON")
	}
}
