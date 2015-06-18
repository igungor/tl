package tl

import (
	"bytes"
	"reflect"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	// single token tests
	var tests = []struct {
		s   string
		tok Token
	}{
		// operators and delimiters
		{s: ``, tok: Token{Token: ItemEOF}},
		{s: ` `, tok: Token{Token: ItemWhitespace, Literal: " "}},
		{s: `  `, tok: Token{Token: ItemWhitespace, Literal: "  "}},
		{s: "\t", tok: Token{Token: ItemWhitespace, Literal: "\t"}},
		{s: `:`, tok: Token{Token: ItemColon, Literal: ":"}},
		{s: `;`, tok: Token{Token: ItemSemicolon, Literal: ";"}},
		{s: `(`, tok: Token{Token: ItemOpenPar, Literal: "("}},
		{s: `)`, tok: Token{Token: ItemClosePar, Literal: ")"}},
		{s: `[`, tok: Token{Token: ItemOpenBracket, Literal: "["}},
		{s: `]`, tok: Token{Token: ItemCloseBracket, Literal: "]"}},
		{s: `{`, tok: Token{Token: ItemOpenBrace, Literal: "{"}},
		{s: `}`, tok: Token{Token: ItemCloseBrace, Literal: "}"}},
		{s: `<`, tok: Token{Token: ItemLeftAngle, Literal: "<"}},
		{s: `>`, tok: Token{Token: ItemRightAngle, Literal: ">"}},
		{s: `%`, tok: Token{Token: ItemPercent, Literal: "%"}},
		{s: `?`, tok: Token{Token: ItemQuestionMark, Literal: "?"}},
		{s: `!`, tok: Token{Token: ItemExclMark, Literal: "!"}},
		{s: `*`, tok: Token{Token: ItemAsterisk, Literal: "*"}},
		{s: `+`, tok: Token{Token: ItemPlus, Literal: "+"}},
		{s: `=`, tok: Token{Token: ItemEquals, Literal: "="}},
		{s: `_`, tok: Token{Token: ItemUnderscore, Literal: "_"}},
		{s: `.`, tok: Token{Token: ItemDot, Literal: "."}},
		{s: `,`, tok: Token{Token: ItemComma, Literal: ","}},
		{s: `@`, tok: Token{Token: ItemIllegal, Literal: "@"}},
		{s: `#`, tok: Token{Token: ItemHash, Literal: "#"}},
		{s: `---`, tok: Token{Token: ItemTripleMinus, Literal: "---"}},
		{s: `---functions---`, tok: Token{Token: ItemTripleMinus, Literal: "---"}},
		{s: `#decafbad`, tok: Token{Token: ItemHash, Literal: "#"}},

		// nat-const
		{s: `0`, tok: Token{Token: ItemNatConst, Literal: "0"}},
		{s: `12`, tok: Token{Token: ItemNatConst, Literal: "12"}},
		{s: `90123`, tok: Token{Token: ItemNatConst, Literal: "90123"}},

		// idents and ident-likes
		{s: `New`, tok: Token{Token: ItemNew, Literal: "New"}},
		{s: `Empty`, tok: Token{Token: ItemEmpty, Literal: "Empty"}},
		{s: `Final`, tok: Token{Token: ItemFinal, Literal: "Final"}},
		{s: `Newly`, tok: Token{Token: ItemUpperIdent, Literal: "Newly"}},
		{s: `Final_countdown`, tok: Token{Token: ItemUpperIdent, Literal: "Final_countdown"}},
		{s: `EmptyHands`, tok: Token{Token: ItemUpperIdent, Literal: "EmptyHands"}},
		{s: `functions---`, tok: Token{Token: ItemLowerIdent, Literal: "functions"}},
		{s: `getUser`, tok: Token{Token: ItemLowerIdent, Literal: "getUser"}},
		{s: `GetUser`, tok: Token{Token: ItemUpperIdent, Literal: "GetUser"}},
		{s: `int128`, tok: Token{Token: ItemLowerIdent, Literal: "int128"}},
		{s: `Int128`, tok: Token{Token: ItemUpperIdent, Literal: "Int128"}},
		{s: `user`, tok: Token{Token: ItemLowerIdent, Literal: "user"}},
		{s: `user#decafbad`, tok: Token{Token: ItemLowerIdent, Literal: "user#decafbad"}},
		{s: `users.user`, tok: Token{Token: ItemLowerIdent, Literal: "users.user"}},
		{s: `users.user#decafbad`, tok: Token{Token: ItemLowerIdent, Literal: "users.user#decafbad"}},
		{s: `User`, tok: Token{Token: ItemUpperIdent, Literal: "User"}},
		{s: `users.User`, tok: Token{Token: ItemUpperIdent, Literal: "users.User"}},

		// Illegal
		{s: `--a`, tok: Token{Token: ItemIllegal, Literal: "--a"}},
		{s: `-aa`, tok: Token{Token: ItemIllegal, Literal: "-aa"}},
	}

	for i, tt := range tests {
		scanner := NewScanner(bytes.NewBufferString(tt.s))
		token := scanner.Scan()

		if !reflect.DeepEqual(token, tt.tok) {
			t.Errorf("<%d> bad token for %q: got %#v, expected %#v", i, tt.s, token, tt.tok)
		}

		if scanner.Err() != nil {
			t.Fatal(scanner.Err())
		}
	}

	// multi token tests
	var multitokentests = []struct {
		s      string
		tokens []Token
	}{
		{`int ? = Int;`, []Token{
			{ItemLowerIdent, "int"},
			{ItemWhitespace, " "},
			{ItemQuestionMark, "?"},
			{ItemWhitespace, " "},
			{ItemEquals, "="},
			{ItemWhitespace, " "},
			{ItemUpperIdent, "Int"},
			{ItemSemicolon, ";"},
		},
		},
		{`users.user#decafbad;`, []Token{
			{ItemLowerIdent, "users.user#decafbad"},
			{ItemSemicolon, ";"},
		},
		},
	}

	for _, mt := range multitokentests {
		scanner := NewScanner(bytes.NewBufferString(mt.s))

		for _, tok := range mt.tokens {
			token := scanner.Scan()

			if !reflect.DeepEqual(token, tok) {
				t.Errorf("got %#v, expected %#v", token, tok)
			}

			if scanner.Err() != nil {
				t.Fatal(scanner.Err())
			}
		}
	}
}
