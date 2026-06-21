// Copyright 2026 The go-jscontact Authors
// SPDX-License-Identifier: Apache-2.0

package jscontact_test

import (
	"encoding/json"
	"fmt"

	jscontact "github.com/hstern/go-jscontact"
)

func ExampleParse() {
	card, err := jscontact.Parse([]byte(`{"@type":"Card","version":"1.0","uid":"abc","name":{"full":"John Doe"}}`))
	if err != nil {
		panic(err)
	}
	fmt.Println(card.UID, card.Name.Full)
	// Output: abc John Doe
}

func ExampleCard_Validate() {
	c := &jscontact.Card{Type: "Card", Version: "1.0"} // missing uid
	fmt.Println(c.Validate())
	// Output: jscontact: uid: is required (RFC 9553 §2.1.9)
}

// Marshaling a Card emits the "@type" discriminator first and defaults it to
// "Card" when unset.
func ExampleCard_MarshalJSON() {
	c := &jscontact.Card{Version: "1.0", UID: "abc"}
	b, _ := json.Marshal(c)
	fmt.Println(string(b))
	// Output: {"@type":"Card","version":"1.0","uid":"abc"}
}

func ExampleDecodeExtra() {
	c, _ := jscontact.Parse([]byte(`{"@type":"Card","version":"1.0","uid":"x","x-rank":{"score":42}}`))
	var rank struct {
		Score int `json:"score"`
	}
	ok, _ := jscontact.DecodeExtra(c.Extra, "x-rank", &rank)
	fmt.Println(ok, rank.Score)
	// Output: true 42
}
