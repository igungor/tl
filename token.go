//go:generate stringer -type=Item -output=token_string.go

package tl

import "fmt"
import "strings"

// Token represents a lexical token.
type Token struct {
	// Type of token
	Token Item

	// Literal value of token
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("<Token: %+q><Literal: %+q>", t.Token, t.Literal)
}

// HasName reports whether the token literal is lc-ident-full.
func (t Token) HasName() bool {
	// only lc-ident-full has a name
	if t.Token != ItemLowerIdent {
		return false
	}

	fields := strings.Split(t.Literal, "#")
	if len(fields) == 1 || (len(fields) == 2 && fields[1] == "") {
		return false
	}
	return true
}

// HasNamespace reports whether the token literal is lc-ident-ns or uc-ident-ns.
func (t Token) HasNamespace() bool {
	// only lc-ident-full has a name
	if t.Token != ItemLowerIdent || t.Token != ItemUpperIdent {
		return false
	}

	fields := strings.Split(t.Literal, ".")
	if len(fields) == 1 || (len(fields) == 2 && fields[1] == "") {
		return false
	}
	return true
}

// Item represents a lexical token type.
type Item int

// List of tokens - https://core.telegram.org/mtproto/TL-formal#tokens
const (
	ItemIllegal Item = iota
	ItemWhitespace
	ItemEOF

	ItemUnderscore   // _
	ItemColon        // :
	ItemSemicolon    // ;
	ItemOpenPar      // (
	ItemClosePar     // )
	ItemOpenBracket  // [
	ItemCloseBracket // ]
	ItemOpenBrace    // {
	ItemCloseBrace   // }
	ItemLeftAngle    // <
	ItemRightAngle   // >
	ItemTripleMinus  // ---
	ItemEquals       // =
	ItemHash         // #
	ItemExclMark     // !
	ItemQuestionMark // ?
	ItemPercent      // %
	ItemPlus         // +
	ItemComma        // ,
	ItemDot          // .
	ItemAsterisk     // *

	ItemNatConst // 4, 42, 421
	// LowerIdent represents lc-ident, lc-ident-ns and lc-ident-full
	ItemLowerIdent // user | user#decafbad | users.user | users.user#decafbad
	// UpperIdent represents uc-ident, uc-ident-ns
	ItemUpperIdent // User | group.User

	ItemFinal // Final
	ItemNew   // New
	ItemEmpty // Empty
)
