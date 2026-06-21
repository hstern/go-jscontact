// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateValid(t *testing.T) {
	c := &Card{Type: "Card", Version: "1.0", UID: "abc-123"}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
}

func TestValidateMissingMusts(t *testing.T) {
	c := &Card{} // no @type defaulting until marshal; empty version+uid
	err := c.Validate()
	if err == nil {
		t.Fatal("expected validation error")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("want *ValidationError, got %T", err)
	}
	paths := map[string]bool{}
	for _, fe := range ve.Errors {
		paths[fe.Path] = true
	}
	for _, want := range []string{"version", "uid"} {
		if !paths[want] {
			t.Errorf("missing expected error for %q (got %v)", want, ve.Errors)
		}
	}
}

func TestValidateBadType(t *testing.T) {
	c := &Card{Type: "Contact", Version: "1.0", UID: "x"}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "@type") {
		t.Fatalf("expected @type error, got %v", err)
	}
}

func TestValidateMembersRequireGroup(t *testing.T) {
	c := &Card{Type: "Card", Version: "1.0", UID: "x", Members: map[string]bool{"u1": true}}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "members") {
		t.Fatalf("expected members/group error, got %v", err)
	}
	c.Kind = "group"
	if err := c.Validate(); err != nil {
		t.Fatalf("group card with members should be valid, got %v", err)
	}
}

func TestValidatePrefRange(t *testing.T) {
	c := &Card{
		Type: "Card", Version: "1.0", UID: "x",
		Emails: map[string]EmailAddress{"e1": {Address: "a@b.c", Pref: 200}},
	}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "pref") {
		t.Fatalf("expected pref-range error, got %v", err)
	}
}

func TestValidateRequiredNestedFields(t *testing.T) {
	c := &Card{
		Type: "Card", Version: "1.0", UID: "x",
		Emails: map[string]EmailAddress{"e1": {}},                // missing address
		Phones: map[string]Phone{"p1": {}},                       // missing number
		Media:  map[string]Media{"m1": {URI: "https://x/a.png"}}, // missing kind
	}
	err := c.Validate()
	if err == nil {
		t.Fatal("expected errors for missing nested fields")
	}
	msg := err.Error()
	for _, want := range []string{`emails["e1"].address`, `phones["p1"].number`, `media["m1"].kind`} {
		if !strings.Contains(msg, want) {
			t.Errorf("missing %q in %q", want, msg)
		}
	}
}

func TestValidateBadIDKey(t *testing.T) {
	c := &Card{
		Type: "Card", Version: "1.0", UID: "x",
		Emails: map[string]EmailAddress{"bad key!": {Address: "a@b.c"}},
	}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "Id key") {
		t.Fatalf("expected Id-key error, got %v", err)
	}
}

func TestValidIDBoundaries(t *testing.T) {
	if validID("") || validID(strings.Repeat("a", 256)) {
		t.Error("empty and >255 must be invalid")
	}
	if !validID("a") || !validID(strings.Repeat("a", 255)) || !validID("A_z-9") {
		t.Error("valid Ids rejected")
	}
	if validID("a b") || validID("café") {
		t.Error("invalid characters accepted")
	}
}
