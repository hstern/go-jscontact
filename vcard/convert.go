// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

// Package vcard converts between the JSContact Card object model
// (github.com/hstern/go-jscontact) and vCard (RFC 6350), following the
// mapping defined by RFC 9555 (JSContact: Converting from and to vCard).
//
// This is the one part of go-jscontact that takes an external dependency:
// it delegates all vCard parsing and serialization to
// github.com/emersion/go-vcard rather than re-implementing RFC 6350. The
// core jscontact package stays dependency-free; only importers of this
// sub-package pull in go-vcard.
//
// Coverage. FromVCard and ToVCard map the common, widely-used properties:
// UID, KIND, FN/N, NICKNAME, ORG, TITLE/ROLE, EMAIL, TEL, ADR, IMPP, URL,
// PHOTO/LOGO/SOUND, KEY, NOTE, CATEGORIES, BDAY/ANNIVERSARY, LANG, PRODID,
// REV, MEMBER, and RELATED. Properties outside this set are not dropped on
// the vCard→JSContact path: unmapped vCard fields are not yet carried, which
// is the documented v0.1.0 limitation (see the package tests). The mapping
// is intentionally asymmetric where RFC 9555 is — round-tripping is lossy
// for properties the two formats model differently.
package vcard

import (
	"strconv"
	"strings"

	vc "github.com/emersion/go-vcard"
	jscontact "github.com/hstern/go-jscontact"
)

// contextsFromParams derives JSContact contexts from vCard TYPE parameters
// (RFC 9555 §2.2): home → "private", work → "work".
func contextsFromParams(p vc.Params) jscontact.Contexts {
	var ctx jscontact.Contexts
	for _, t := range p.Types() {
		switch strings.ToLower(t) {
		case vc.TypeHome:
			ctx = setTrue(ctx, "private")
		case vc.TypeWork:
			ctx = setTrue(ctx, "work")
		}
	}
	return ctx
}

// typesFromContexts is the inverse of contextsFromParams.
func typesFromContexts(ctx jscontact.Contexts) []string {
	var types []string
	if ctx["private"] {
		types = append(types, vc.TypeHome)
	}
	if ctx["work"] {
		types = append(types, vc.TypeWork)
	}
	return types
}

// phoneFeatureFromType maps a vCard TEL TYPE to a JSContact phone feature,
// returning ("", false) for TYPE values that are contexts, not features.
var phoneFeatureFromType = map[string]string{
	vc.TypeVoice:     "voice",
	vc.TypeCell:      "mobile",
	vc.TypeFax:       "fax",
	vc.TypeVideo:     "video",
	vc.TypePager:     "pager",
	vc.TypeText:      "text",
	vc.TypeTextPhone: "textphone",
}

// typeFromPhoneFeature is the inverse of phoneFeatureFromType.
var typeFromPhoneFeature = map[string]string{
	"voice":     vc.TypeVoice,
	"mobile":    vc.TypeCell,
	"fax":       vc.TypeFax,
	"video":     vc.TypeVideo,
	"pager":     vc.TypePager,
	"text":      vc.TypeText,
	"textphone": vc.TypeTextPhone,
}

// prefFromParams reads a vCard PREF parameter (1–100) as a JSContact pref;
// 0 means "absent".
func prefFromParams(p vc.Params) uint64 {
	if s := p.Get(vc.ParamPreferred); s != "" {
		if n, err := strconv.ParseUint(s, 10, 64); err == nil {
			return n
		}
	}
	return 0
}

// paramsFor builds vCard parameters from JSContact contexts and pref.
func paramsFor(ctx jscontact.Contexts, pref uint64) vc.Params {
	p := make(vc.Params)
	for _, t := range typesFromContexts(ctx) {
		p.Add(vc.ParamType, t)
	}
	if pref != 0 {
		p.Set(vc.ParamPreferred, strconv.FormatUint(pref, 10))
	}
	return p
}

func setTrue(m jscontact.Contexts, k string) jscontact.Contexts {
	if m == nil {
		m = make(jscontact.Contexts)
	}
	m[k] = true
	return m
}

func setFeature(m jscontact.Features, k string) jscontact.Features {
	if m == nil {
		m = make(jscontact.Features)
	}
	m[k] = true
	return m
}
