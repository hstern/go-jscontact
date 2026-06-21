// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
)

// This file holds the shared codec machinery that every JSContact object
// type reuses. Two invariants drive the design (RFC 9553 §1.4):
//
//   - The "@type" discriminator is emitted first. Because encoding/json
//     marshals struct fields in declaration order, every object type
//     declares Type as its first field; marshalExtra preserves that order.
//   - Unknown members round-trip losslessly. Each object carries an
//     Extra map (tagged json:"-"); unmarshalExtra routes every member the
//     typed struct does not recognize into Extra, and marshalExtra splices
//     them back in, so a decode→encode cycle never drops data.
//
// The per-type MarshalJSON/UnmarshalJSON methods are one-line wrappers
// around these helpers, using a local alias type to shed the method set
// (and thus avoid infinite recursion through json.Marshal).

// marshalExtra marshals v with standard struct-field ordering, then splices
// the open-extension members from extra into the resulting JSON object.
// v must be an alias type with no MarshalJSON method of its own.
func marshalExtra(v any, extra map[string]json.RawMessage) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if len(extra) == 0 {
		return b, nil
	}
	keys := make([]string, 0, len(extra))
	for k := range extra {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.Grow(len(b) + 32*len(keys))
	buf.Write(b[:len(b)-1]) // everything but the closing '}'
	if len(b) > 2 {         // object already has members (not "{}")
		buf.WriteByte(',')
	}
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		if err := json.Compact(&buf, extra[k]); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// unmarshalExtra decodes b into the typed struct v, then collects every
// member of b that v does not declare into *extra. v must be a pointer to
// an alias type with no UnmarshalJSON method of its own.
func unmarshalExtra(b []byte, v any, extra *map[string]json.RawMessage) error {
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	var all map[string]json.RawMessage
	if err := json.Unmarshal(b, &all); err != nil {
		return err
	}
	known := knownJSONKeys(reflect.TypeOf(v))
	for k := range all {
		if known[k] {
			delete(all, k)
		}
	}
	if len(all) > 0 {
		*extra = all
	} else {
		*extra = nil
	}
	return nil
}

// knownJSONKeys returns the set of JSON member names a struct type declares,
// excluding fields tagged json:"-" (notably the Extra field itself).
func knownJSONKeys(t reflect.Type) map[string]bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	m := make(map[string]bool, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		name, _, _ := strings.Cut(tag, ",")
		if name == "" {
			name = f.Name
		}
		m[name] = true
	}
	return m
}
