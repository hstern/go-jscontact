// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package vcard_test

import (
	"fmt"
	"strings"

	govcard "github.com/emersion/go-vcard"
	jsvcard "github.com/hstern/go-jscontact/vcard"
)

func ExampleFromVCard() {
	const vcf = "BEGIN:VCARD\r\n" +
		"VERSION:4.0\r\n" +
		"UID:urn:uuid:abc\r\n" +
		"FN:Jane Roe\r\n" +
		"EMAIL;TYPE=work:jane@example.com\r\n" +
		"END:VCARD\r\n"

	src, err := govcard.NewDecoder(strings.NewReader(vcf)).Decode()
	if err != nil {
		panic(err)
	}
	card, err := jsvcard.FromVCard(src)
	if err != nil {
		panic(err)
	}
	fmt.Println(card.UID, card.Name.Full)
	// Output: urn:uuid:abc Jane Roe
}
