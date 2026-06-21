// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import "encoding/json"

// Helpers for working with open-extension members (the Extra maps that every
// JSContact object carries; RFC 9553 §1.5 reserved/extension members). They
// keep the public surface free of any-typed values: a consumer decodes an
// extension into its own typed struct.

// DecodeExtra unmarshals the extension member named key from extra into dst.
// It reports whether the member was present. A nil or absent member is
// treated as "not present" and leaves dst untouched.
func DecodeExtra(extra map[string]json.RawMessage, key string, dst any) (bool, error) {
	raw, ok := extra[key]
	if !ok || len(raw) == 0 {
		return false, nil
	}
	if err := json.Unmarshal(raw, dst); err != nil {
		return true, err
	}
	return true, nil
}

// EncodeExtra marshals v and stores it under key in extra, allocating the map
// if needed, and returns the (possibly newly allocated) map.
func EncodeExtra(extra map[string]json.RawMessage, key string, v any) (map[string]json.RawMessage, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return extra, err
	}
	if extra == nil {
		extra = make(map[string]json.RawMessage, 1)
	}
	extra[key] = raw
	return extra, nil
}
