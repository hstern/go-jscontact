// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "testing"

func TestSpecVersion(t *testing.T) {
	if SpecVersion != "RFC 9553" {
		t.Fatalf("SpecVersion = %q, want %q", SpecVersion, "RFC 9553")
	}
}
