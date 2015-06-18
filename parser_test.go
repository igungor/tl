package tl

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParser_parseBoxedTypeIdent(t *testing.T) {
	var tests = []struct {
		s string
		b BoxedTypeIdent
	}{
		{`Int`, BoxedTypeIdent{"Int"}},
		{`Long`, BoxedTypeIdent{"Long"}},
		{`Int128`, BoxedTypeIdent{"Int128"}},
		{`String`, BoxedTypeIdent{"String"}},
		{`Double`, BoxedTypeIdent{"Double"}},
	}

	for i, tt := range tests {
		parser := NewParser(bytes.NewBufferString(tt.s))

		// initial parse
		parser.next()

		b := parser.parseBoxedTypeIdent()

		if parser.Err() != nil {
			t.Errorf("got error: %v", parser.Err())
		}

		if !reflect.DeepEqual(b, tt.b) {
			t.Errorf("<%d> bad type for %q: got %#v, expected %#v", i, tt.b, b, tt.b)
		}
	}
}

func TestParser_parseTypeIdent(t *testing.T) {
	var tests = []struct {
		s string
		b TypeIdent
	}{
		{s: `int`, b: TypeIdent{"int"}},
		{s: `user`, b: TypeIdent{"user"}},
		{s: `users.user`, b: TypeIdent{"users.user"}},
		{s: `User`, b: TypeIdent{"User"}},
		{s: `users.User`, b: TypeIdent{"users.User"}},
		{s: `#`, b: TypeIdent{"#"}},
	}

	for i, tt := range tests {
		parser := NewParser(bytes.NewBufferString(tt.s))

		// initial parse
		parser.next()

		b := parser.parseTypeIdent()

		if parser.Err() != nil {
			t.Errorf("got error: %v", parser.Err())
		}

		if !reflect.DeepEqual(b, tt.b) {
			t.Errorf("<%d> bad type for %q: got %#v, expected %#v", i, tt.b, b, tt.b)
		}
	}
}

func TestParser_parseCombinatorId(t *testing.T) {
	var tests = []struct {
		s  string
		id CombinatorId
	}{
		{s: `int`, id: CombinatorId{NewIdent("int")}},
		{s: `user`, id: CombinatorId{NewIdent("user")}},
		{s: `users.user`, id: CombinatorId{NewIdent("users.user")}},
		{s: `users.user#decafbad`, id: CombinatorId{NewIdent("users.user#decafbad")}},
		{s: `user#decafbad`, id: CombinatorId{NewIdent("user#decafbad")}},
		{s: `_`, id: CombinatorId{NewIdent("_")}},
	}

	for i, tt := range tests {
		parser := NewParser(bytes.NewBufferString(tt.s))

		// initial parse
		parser.next()

		id := parser.parseCombinatorId()

		if parser.Err() != nil {
			t.Errorf("got error: %v", parser.Err())
		}

		if !reflect.DeepEqual(id, tt.id) {
			t.Errorf("<%d> bad type for %q: got %#v, expected %#v", i, tt.id, id, tt.id)
		}
	}
}

func TestParser_parseFullCombinatorId(t *testing.T) {
	var tests = []struct {
		s  string
		id FullCombinatorId
	}{
		{s: `int`, id: FullCombinatorId{NewIdent("int")}},
		{s: `user`, id: FullCombinatorId{NewIdent("user")}},
		{s: `users.user`, id: FullCombinatorId{NewIdent("users.user")}},
		{s: `users.user#decafbad`, id: FullCombinatorId{NewIdent("users.user#decafbad")}},
		{s: `user#decafbad`, id: FullCombinatorId{NewIdent("user#decafbad")}},
		{s: `_`, id: FullCombinatorId{NewIdent("_")}},
	}

	for i, tt := range tests {
		parser := NewParser(bytes.NewBufferString(tt.s))

		// initial parse
		parser.next()

		id := parser.parseFullCombinatorId()

		if parser.Err() != nil {
			t.Errorf("got error: %v", parser.Err())
		}

		if !reflect.DeepEqual(id, tt.id) {
			t.Errorf("<%d> bad type for %q: got %#v, expected %#v", i, tt.id, id, tt.id)
		}
	}
}

func TestParser_parseBuiltinCombinatorDecl(t *testing.T) {
	var tests = []struct {
		s string
		b BuiltinCombDecl
	}{
		{`int ?= Int;`, BuiltinCombDecl{FullCombinatorId{NewIdent("int")}, BoxedTypeIdent{"Int"}}},
		{`long ?= Long;`, BuiltinCombDecl{FullCombinatorId{NewIdent("long")}, BoxedTypeIdent{"Long"}}},
		{`double ?= Double;`, BuiltinCombDecl{FullCombinatorId{NewIdent("double")}, BoxedTypeIdent{"Double"}}},
		{`string ?= String;`, BuiltinCombDecl{FullCombinatorId{NewIdent("string")}, BoxedTypeIdent{"String"}}},
	}

	for i, tt := range tests {
		parser := NewParser(bytes.NewBufferString(tt.s))

		// initial parse
		parser.next()

		b := parser.parseBuiltinCombinatorDecl()

		if parser.Err() != nil {
			t.Errorf("got error: %v", parser.Err())
		}

		if !reflect.DeepEqual(b, tt.b) {
			t.Errorf("<%d> bad type for %q: got %#v, expected %#v", i, tt.b, b, tt.b)
		}
	}
}
