package tl

import (
	"fmt"
	"io"
	"log"
)

// Parser holds the parser's internal state while consuming a given
// token.
type Parser struct {
	s   *Scanner
	tok Token // one token look-ahead
	err error // sticky error

	Trace  bool // parsing mode
	indent int  // indentation used for tracing output
}

// NewParser returns a Parser from the given io.Reader.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Err() error {
	return p.err
}

func (p *Parser) setErr(err error) {
	if p.err == nil {
		p.err = err
	}
}

// Parse is the entry-point to the parser.
//
// TL-program ::= constr-declarations { --- functions --- fun-declarations | --- types --- constr-declarations }
//
func (p *Parser) Parse() (program Program) {
	defer un(trace(p, "ParseProgram"))

	p.next()

	// constr-declarations
	for {
		if p.tok.Token == ItemEOF {
			return
		}

		p.parseDecl()
		p.expectSemi()

		if p.tok.Token == ItemTripleMinus {
			if p.tok.Literal == "functions" {
				p.next()
				p.expect(ItemTripleMinus)
				break
			} else {
				p.setErr(fmt.Errorf("expected functions separator"))
			}
		}
	}

	// fun-declarations
	for {
		if p.tok.Token == ItemEOF {
			return
		}

		p.parseDecl()
		p.expectSemi()

		if p.tok.Token == ItemTripleMinus {
			if p.tok.Literal == "types" {
				break
			} else {
				p.setErr(fmt.Errorf("expected types separator"))
			}
		}
	}

	// types-declarations
	for {
		if p.tok.Token == ItemEOF {
			return
		}

		p.parseDecl()
		p.expectSemi()

		if p.tok.Token == ItemTripleMinus {
			p.setErr(fmt.Errorf("unexpected separator"))
		}
	}
}

// parseConstrDecls consumes constructor-declarations till EOF or function separator.
//
// constr-declarations ::= { declaration ; }
// constr-declaration ::= combinator-decl | partial-app-decl | final-decldeclaration ;
//
func (p *Parser) parseConstrDecl() {
	defer un(trace(p, "parseConstrDecl"))
}

// parseFuncDecls consumes function-declarations till EOF or type separator.
//
// fun-declarations ::= { declaration ; }
// function-declaration ::= combinator-decl | partial-app-decl | final-decldeclaration ;
//
func (p *Parser) parseFuncDecl() {
	defer un(trace(p, "parseFuncDecl"))
}

// parseDeclaration consumes a generic declaration.
//
// declaration ::= combinator-decl | partial-app-decl | final-decl
//
func (p *Parser) parseDecl() {
	defer un(trace(p, "parseDecl"))

	tok := p.tok

	switch tok.Token {
	case ItemLowerIdent:
		p.parseBuiltinCombinatorDecl()
	case ItemUpperIdent:
		p.parsePartialAppDecl()
	case ItemNew, ItemFinal, ItemEmpty:
		p.parseFinalDecl()
	default:
		p.setErr(fmt.Errorf("unexpected token"))
		return
	}
}

// parseBuiltinCombinatorDecl consumes a builtin combinator declaration.
//
// builtin-combinator-decl ::= full-combinator-id ? = boxed-type-ident ;
//
func (p *Parser) parseBuiltinCombinatorDecl() BuiltinCombDecl {
	defer un(trace(p, "parseBuiltinCombinatorDecl"))

	id := p.parseFullCombinatorId()
	p.expect(ItemQuestionMark)
	p.expect(ItemEquals)
	result := p.parseBoxedTypeIdent()
	p.expectSemi()

	return BuiltinCombDecl{id, result}
}

// parseCombinatorDecl consumes a combinator declaration.
//
// combinator-decl ::= full-combinator-id { opt-args } { args } `=` result-type `;`
//
// user#decafbad {id:int} name:string = User;
//
func (p *Parser) parseCombinatorDecl() {
	defer un(trace(p, "parseCombinatorDecl"))

	p.next()
}

// parsePartialAppDecl consumes a partial-app-decl.
//
// partial-app-decl ::= partial-type-app-decl | partial-comb-app-decl
// partial-type-app-decl ::= boxed-type-ident subexpr { subexpr } ; | boxed-type-ident < expr { , expr } > ;
// partial-comb-app-decl ::= combinator-id subexpr { subexpr } ;
//
func (p *Parser) parsePartialAppDecl() {
	defer un(trace(p, "parsePartialAppDecl"))

	p.next()
}

// parseFinalDeclaration consumes a final declaration.
//
// final-decl ::= New boxed-type-ident ; | Final boxed-type-ident ; | Empty boxed-type-ident ;
//
func (p *Parser) parseFinalDecl() {
	defer un(trace(p, "parseFinalDecl"))

	p.next()
}

// parseResultType consumes a result-type.
//
// result-type ::= boxed-type-ident < subexpr { , subexpr } >
//
func (p *Parser) parseResultType() {
	defer un(trace(p, "parseResultType"))
}

// parseOptionalArgs parses a combinator's optional arguments. All optional
// fields must be explicitly named. (using '_' is not allowed)
//
// opt-args ::= '{' var-ident { var-ident } : [excl-mark] type-expr '}'
//
func (p *Parser) parseOptionalArgs() {
	defer un(trace(p, "parseOptionalArgs"))
}

// parseArgs consumes required arguments of a combinator declaration.
//
// Required arguments can be grouped by parens, but not required. Also, '_' can
// be used as names of one or more fields. An anonymous field may be declared
// using a type-entry, functionally equivalent to `_:int`:
//
//   `getUser (id:int) = User;`
//   `getUser id:int = User;`
//   `getUser (_:int) = User;`
//   `getUser _:int = User;`
//   `getUser int = User;`
//
// args ::= var-ident-opt ':' [ conditional-def ] [ '!' ] type-term
// args ::= [ var-ident-opt ':' ] [ multiplicity *] '[' { args } ']'
// args ::= '(' var-ident-opt { var-ident-opt } : [!] type-term ')'
// args ::= [ '!' ] type-term
//
func (p *Parser) parseArgs() {
	defer un(trace(p, "parseArgs"))
}

// parseExpr consumes multiple subexpr's.
//
// expr ::= { subexpr }
//
func (p *Parser) parseExpr() {
	defer un(trace(p, "parseExpr"))
}

// parseSubExpr consumes a subexpr.
//
// subexpr ::= term | nat-const '+' subexpr | subexpr '+' nat-const
//
func (p *Parser) parseSubExpr() {
	defer un(trace(p, "parseSubExpr"))
}

// parseTerm consumes a term.
//
// term ::= '(' expr ')' | type-ident | var-ident | nat-const | % term | type-ident '<' expr { ',' expr } '>'
//
func (p *Parser) parseTerm() {
	defer un(trace(p, "parseTerm"))
}

// parseTypeIdent consumes a type-ident.
//
// type-ident ::= boxed-type-ident | lc-ident-ns | #
//
func (p *Parser) parseTypeIdent() TypeIdent {
	defer un(trace(p, "parseTypeIdent"))

	tok := p.tok
	if tok.Token == ItemLowerIdent && tok.HasName() {
		p.error("got lc-ident-full, expected lc-ident-ns")
	}

	p.expect(ItemUpperIdent, ItemLowerIdent, ItemHash)

	return TypeIdent{tok.Literal}
}

// parseBoxedType consumes a boxed-type-ident.
//
// boxed-type-ident ::= uc-ident-ns
//
func (p *Parser) parseBoxedTypeIdent() BoxedTypeIdent {
	defer un(trace(p, "parseBoxedTypeIdent"))

	name := p.tok.Literal
	p.expect(ItemUpperIdent)

	return BoxedTypeIdent{name}
}

// parseCombinatorId consumes a combinator-id
//
// combinator-id ::= lc-ident-ns | _
// e.g.: user | group.user | _
//
func (p *Parser) parseCombinatorId() CombinatorId {
	// FIXME: handle lc-ident-full case.
	defer un(trace(p, "parseCombinatorId"))

	tok := p.tok
	p.expect(ItemLowerIdent, ItemUnderscore)

	return CombinatorId{Ident{tok}}
}

// parseFullCombinatorId consumes a full-combinator-id
//
// full-combinator-id ::= lc-ident-full | _
// e.g.: user#decafbad | user | group.user | group.user#decafbad
//
func (p *Parser) parseFullCombinatorId() FullCombinatorId {
	defer un(trace(p, "parseFullCombinatorId"))

	tok := p.tok
	p.expect(ItemLowerIdent, ItemUnderscore)

	return FullCombinatorId{Ident{tok}}
}

// next advances to the next non-whitespace token.
func (p *Parser) next() {
	for {
		p.tok = p.s.Scan()

		// skip whitespace
		if p.tok.Token != ItemWhitespace {
			break
		}
	}
}

// expect checks if the current token is in the given items list, then advances
// to the next non-whitespace token.
func (p *Parser) expect(items ...Item) {
	var errc int

	for _, item := range items {
		if p.tok.Token == item {
			errc = 0
			break
		}
		errc += 1
	}

	if errc != 0 {
		p.errorExpected(fmt.Sprintf("got: %s, want: '%v'", p.tok, items[:errc]))
	}
	p.next()
}

func (p *Parser) expectSemi() {
	if p.tok.Token != ItemSemicolon {
		p.setErr(fmt.Errorf("expected semicolon"))
	}

	p.next()
}

func (p *Parser) printTrace(a ...interface{}) {
	if !p.Trace {
		return
	}

	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . . "
	const n = len(dots)
	i := 2 * p.indent
	for i > n {
		fmt.Print(dots)
		i -= n
	}
	// i <= n
	fmt.Print(dots[0:i])
	fmt.Println(a...)
}

func trace(p *Parser, msg string) *Parser {
	p.printTrace(msg, "(")
	p.indent++
	return p
}

func un(p *Parser) {
	p.indent--
	p.printTrace(")")
}

// func (p *Parser) error(msg string)         { p.errors = append(p.errors, errors.New(msg)) }
func (p *Parser) error(msg string)         { log.Fatal(msg) }
func (p *Parser) errorExpected(msg string) { p.error("expected " + msg) }
