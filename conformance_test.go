// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// TestConformanceRFC9553 drives every fixture under testdata/rfc9553 through
// the library and asserts the RFC 9553 conformance properties: each fixture
// parses, validates against the §2 MUSTs, and re-marshals idempotently
// (Parse→Marshal→Parse→Marshal is byte-stable), which exercises @type-first
// ordering, set/ordered-component fidelity, and extension passthrough.
func TestConformanceRFC9553(t *testing.T) {
	files, err := filepath.Glob("testdata/rfc9553/*.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no conformance fixtures found")
	}
	for _, f := range files {
		f := f
		t.Run(filepath.Base(f), func(t *testing.T) {
			raw, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			card, err := Parse(raw)
			if err != nil {
				t.Fatalf("Parse: %v", err)
			}
			if err := card.Validate(); err != nil {
				t.Fatalf("Validate: %v", err)
			}

			m1, err := json.Marshal(card)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			again, err := Parse(m1)
			if err != nil {
				t.Fatalf("re-Parse: %v", err)
			}
			m2, err := json.Marshal(again)
			if err != nil {
				t.Fatalf("re-Marshal: %v", err)
			}
			if string(m1) != string(m2) {
				t.Fatalf("re-marshal not idempotent:\n1: %s\n2: %s", m1, m2)
			}

			// The canonical re-marshal must always lead with the @type.
			if len(m1) < 16 || string(m1[:16]) != `{"@type":"Card",` {
				t.Fatalf("@type not emitted first: %s", m1[:min(40, len(m1))])
			}
		})
	}
}
