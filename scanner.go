package tl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

const (
	eof        rune = -1
	notReadYet rune = -2
)

// Scanner implements a TL (Type Language) lexer.
type Scanner struct {
	r   *bufio.Reader
	ch  rune  // one rune look-ahead
	err error // sticky error
}

// NewScanner returns a Scanner which tokenizes a TL program
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:  bufio.NewReader(r),
		ch: notReadYet,
	}
}

// Scan returns the next token from the underlying reader.
func (s *Scanner) Scan() Token {
	ch := s.peek()

	switch {
	case isWhitespace(ch):
		return s.scanWhitespace()
	case isLetter(ch):
		return s.scanIdent()
	case isDigit(ch):
		return s.scanNumber()
	default:
		ch := s.Next() // always advance.
		switch ch {
		case eof:
			return Token{Token: ItemEOF}
		case '/':
			if s.ch == '/' {
				s.r.ReadBytes('\n')
			}
		case '-':
			return s.scanTripleMinus()
		case '#':
			return Token{ItemHash, string(ch)}
		case '.':
			return Token{ItemDot, string(ch)}
		case ',':
			return Token{ItemComma, string(ch)}
		case ':':
			return Token{ItemColon, string(ch)}
		case ';':
			return Token{ItemSemicolon, string(ch)}
		case '_':
			return Token{ItemUnderscore, string(ch)}
		case '=':
			return Token{ItemEquals, string(ch)}
		case '%':
			return Token{ItemPercent, string(ch)}
		case '?':
			return Token{ItemQuestionMark, string(ch)}
		case '!':
			return Token{ItemExclMark, string(ch)}
		case '*':
			return Token{ItemAsterisk, string(ch)}
		case '+':
			return Token{ItemPlus, string(ch)}
		case '(':
			return Token{ItemOpenPar, string(ch)}
		case ')':
			return Token{ItemClosePar, string(ch)}
		case '{':
			return Token{ItemOpenBrace, string(ch)}
		case '}':
			return Token{ItemCloseBrace, string(ch)}
		case '[':
			return Token{ItemOpenBracket, string(ch)}
		case ']':
			return Token{ItemCloseBracket, string(ch)}
		case '<':
			return Token{ItemLeftAngle, string(ch)}
		case '>':
			return Token{ItemRightAngle, string(ch)}
		default:
			return Token{ItemIllegal, string(ch)}
		}
	}

	panic("unreachable")
}

// Next reads and returns the next rune from the underlying reader.
func (s *Scanner) Next() rune {
	next := s.peek()
	s.ch = s.next()
	return next
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// scanTripleMinus consumes triple-minus separator.
func (s *Scanner) scanTripleMinus() Token {
	minus1, minus2 := s.Next(), s.Next()

	if minus1 == '-' && minus2 == '-' {
		return Token{ItemTripleMinus, string("---")}
	}

	return Token{ItemIllegal, string("-") + string(minus1) + string(minus2)}
}

// scanWhitespace consumes the current rune and all contiguous whitespaces.
func (s *Scanner) scanWhitespace() Token {
	var buf bytes.Buffer

	for isWhitespace(s.ch) {
		buf.WriteRune(s.ch)
		s.Next()
	}

	return Token{ItemWhitespace, buf.String()}
}

// scanIdent consumes an identifier-like token.
//
// ident ::= letter { ident-char }
//
func (s *Scanner) scanIdent() Token {
	var buf bytes.Buffer
	var item Item

	// definitely letter. lower or upper?

	// handle uppercase case
	if isUpperLetter(s.ch) {
		for isIdentChar(s.ch) {
			buf.WriteRune(s.ch)
			s.Next()
		}

		switch buf.String() {
		case "New":
			item = ItemNew
		case "Empty":
			item = ItemEmpty
		case "Final":
			item = ItemFinal
		default:
			item = ItemUpperIdent
		}

		return Token{item, buf.String()}
	}

	// handle lowercase case.

	// consume all ident-chars.
	for isIdentChar(s.ch) {
		buf.WriteRune(s.ch)
		s.Next()
	}

	// handle namespace
	if s.ch == '.' {
		buf.WriteRune(s.ch)

		s.Next()

		// uc-ident-ns
		if isUpperLetter(s.ch) {
			// consume all ident-chars
			for isIdentChar(s.ch) {
				buf.WriteRune(s.ch)
				s.Next()
			}

			return Token{ItemUpperIdent, buf.String()}
		}

		// lc-ident-ns
		if isLowerLetter(s.ch) {
			// consume all ident-chars
			for isIdentChar(s.ch) {
				buf.WriteRune(s.ch)
				s.Next()
			}
		} else {
			s.setErr(fmt.Errorf("expected letter, got %q", s.ch))
		}
	}

	// lc-ident-full
	if s.ch == '#' {
		buf.WriteRune(s.ch)

		// expect 8 hex-digits
		for i := 0; i < 8; i++ {
			s.Next()
			if !isHexDigit(s.ch) {
				s.setErr(fmt.Errorf("expected hexdigit, got %q", s.ch))
			}
			buf.WriteRune(s.ch)
		}
		s.Next()
	}

	return Token{ItemLowerIdent, buf.String()}
}

// scanNumber consumes a number.
func (s *Scanner) scanNumber() Token {
	var buf bytes.Buffer

	for isDigit(s.ch) {
		buf.WriteRune(s.ch)
		s.Next()
	}

	return Token{ItemNatConst, buf.String()}
}

func (s *Scanner) next() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		ch = eof
	}
	return ch
}

// peek return the next rune in the reader without advancing the scanner.
func (s *Scanner) peek() rune {
	if s.ch == notReadYet {
		s.ch = s.next()
	}

	return s.ch
}

// setErr records the first error encountered.
func (s *Scanner) setErr(err error) {
	if s.err == nil || s.err == io.EOF {
		s.err = err
	}
}

// isLowerLetter reports whether the rune is a lowercase letter.
//
// lc-letter ::= a | b | … | z
//
func isLowerLetter(ch rune) bool { return 'a' <= ch && ch <= 'z' }

// isUpperLetter reports whether the rune is a uppercase letter.
//
// uc-letter ::= A | B | … | Z
//
func isUpperLetter(ch rune) bool { return 'A' <= ch && ch <= 'Z' }

// isLetter reports whether the rune is a letter.
//
// letter ::= lc-letter | uc-letter
//
func isLetter(ch rune) bool { return isLowerLetter(ch) || isUpperLetter(ch) }

// isDigit reports whether the rune is a digit.
//
// digit ::= 0 | 1 | … | 9
//
func isDigit(ch rune) bool { return '0' <= ch && ch <= '9' }

// isHexDigit reports whether the rune is a hex digit.
//
// hex-digit ::= digit | a | b | c | d | e | f
//
func isHexDigit(ch rune) bool { return isDigit(ch) || ('a' <= ch && ch <= 'f') }

// isIdentChar reports whether the rune is a valid identifier character.
//
// ident_char ::= letter | digit | underscore
//
func isIdentChar(ch rune) bool { return isLetter(ch) || isDigit(ch) || ch == '_' }

// isWhitespace reports whether the rune is a valid whitespace separator.
//
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' }
