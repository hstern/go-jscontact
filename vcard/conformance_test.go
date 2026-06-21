// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package vcard

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	vc "github.com/emersion/go-vcard"
)

// TestConformanceRFC9555 drives every vCard fixture under testdata through the
// RFC 9555 conversion both directions: each parses, converts to a JSContact
// Card that passes Validate, converts back to a vCard, and the round-trip
// preserves the identity field (UID) and the core mapped properties.
func TestConformanceRFC9555(t *testing.T) {
	files, err := filepath.Glob("testdata/*.vcf")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no vCard conformance fixtures found")
	}
	for _, f := range files {
		f := f
		t.Run(filepath.Base(f), func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			src, err := vc.NewDecoder(bytes.NewReader(data)).Decode()
			if err != nil {
				t.Fatalf("decode vCard: %v", err)
			}

			card, err := FromVCard(src)
			if err != nil {
				t.Fatalf("FromVCard: %v", err)
			}
			if err := card.Validate(); err != nil {
				t.Fatalf("converted Card fails Validate: %v", err)
			}

			back, err := ToVCard(card)
			if err != nil {
				t.Fatalf("ToVCard: %v", err)
			}
			if got, want := back.Value(vc.FieldUID), src.Value(vc.FieldUID); got != want {
				t.Errorf("UID not preserved across round-trip: %q vs %q", got, want)
			}
			if src.Value(vc.FieldFormattedName) != "" && back.Value(vc.FieldFormattedName) == "" {
				t.Errorf("FN lost across round-trip")
			}
			if len(src[vc.FieldEmail]) != len(back[vc.FieldEmail]) {
				t.Errorf("EMAIL count changed: %d → %d", len(src[vc.FieldEmail]), len(back[vc.FieldEmail]))
			}
		})
	}
}
