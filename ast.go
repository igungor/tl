package tl

import (
	"fmt"
	"hash/crc32"
	"unicode"
	"unicode/utf8"
)

// Node	represents a node in abstract syntax tree.
type Node interface {
	node()
}

// Declaration represents a constructor or function declaration. All
// declaration nodes implement the Declaration interface.
type Declaration interface {
	Node
	declNode()
}

// Combinator represents a constructor, function or a built-in combinator whose
// combinator-name can be computed.
type Combinator interface {
	Node
	Name() string
}

// Program represents a TL program.
//
// https://core.telegram.org/mtproto/TL-formal
type Program struct {
	Constructors []ConstrDecl

	// Optional
	Functions []FuncDecl
	Types     []ConstrDecl
}

type Expr struct{}
type Ident struct {
	Name Token
}

func (i Ident) Text() string {
	return i.Name.Literal
}

type CombinatorId struct {
	Id Ident
}

type FullCombinatorId struct {
	Id Ident
}

type Arg struct{}
type OptionalArg struct{}

type TypeIdent struct {
	Name string
}
type BoxedTypeIdent struct {
	Name string
}
type BareType struct{}
type ResultType struct{}

// A declaration is represented by one of the following declaration nodes.
//
type (
	CombDecl struct {
		Id      FullCombinatorId
		OptArgs []OptionalArg
		Args    []Arg
		Result  ResultType
	}

	BuiltinCombDecl struct {
		Id     FullCombinatorId
		Result BoxedTypeIdent
	}

	ConstrDecl struct{}
	FuncDecl   struct{}
)

// node implementations for declaration nodes.
//
func (ConstrDecl) node()      {}
func (FuncDecl) node()        {}
func (CombDecl) node()        {}
func (BuiltinCombDecl) node() {}

// declNode() ensures that only declaration nodes can be assigned to a declaration node.
//
func (ConstrDecl) declNode()      {}
func (FuncDecl) declNode()        {}
func (CombDecl) declNode()        {}
func (BuiltinCombDecl) declNode() {}

// combinatorNode implementations for combinator nodes.
//
func (ConstrDecl) Name() {}
func (FuncDecl) Name()   {}
func (CombDecl) Name()   {}

//
// Constructors
//

func NewIdent(name string) Ident {
	ch, _ := utf8.DecodeRuneInString(name)

	if name == "_" {
		return Ident{Name: Token{Token: ItemUnderscore, Literal: name}}
	} else if unicode.IsLower(ch) {
		return Ident{Name: Token{Token: ItemLowerIdent, Literal: name}}
	} else {
		return Ident{Name: Token{Token: ItemUpperIdent, Literal: name}}
	}
}

// computeCRC32 calculates the combinator-name for the given combinator-description.
// e.g. int128 ? = Int128 => a8509bda
func computeCRC32(s string) (string, error) {
	h := crc32.NewIEEE()
	_, err := h.Write([]byte(s))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum32()), nil
}
