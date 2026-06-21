// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

// Package jscontact implements RFC 9553 — JSContact: A JSON Representation
// of Contact Data, the JSON-native successor to vCard.
//
// The typed Card object model, its custom byte-stable JSON codec, and
// Validate land in subsequent development phases; the vCard conversion
// (RFC 9555) lives in the jscontact/vcard sub-package. This file is the
// initial scaffold — see the project README for current status.
package jscontact

// SpecVersion is the RFC 9553 version this build implements.
const SpecVersion = "RFC 9553"
