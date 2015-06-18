# TL (Type Language)

## Notation (Extended Backus-Naur Form - EBNF)

EBNF is used to make a formal description of a formal language which can be a
computer programming language. 

    Production  = production_name "=" [ Expression ] "." .
    Expression  = Alternative { "|" Alternative } .
    Alternative = Term { Term } .
    Term        = production_name | token [ "…" token ] | Group | Option | Repetition .
    Group       = "(" Expression ")" .
    Option      = "[" Expression "]" .
    Repetition  = "{" Expression "}" .

Productions are expressions constructed from terms and the following operators,
in increasing precedence:

    |   alternation
    ()  grouping
    []  option (0 or 1 times)
    {}  repetition (0 to n times)

## Definitions

* `schema` is a collection of all the data type descriptions. This is used to define some
  agreed-to system of types.

* `combinator` is a function that takes arguments of certain types and returns a value of some
  other type.

* `arity` is a non-negative integer, the number of combinator arguments.

* `combinator identifier` is an identifier beginning with a lowercase Roman letter that uniquely
  identifies a combinator.

* `combinator number/name` is a 32-bit number that uniquely identifies a combinator. CRC32 of the
  string containing the combinator description without the final semicolon, with one space between
  contiguous lexemes.

    getUser#9fa45g7c (Vector int) = Vector User;
            ^^^^^^^^ -> combinator name.

    it is calculated with CRC32 sum of `getUser Vector int = Vector User`.

* full-combinator-id: user#01234567 (combinator-id with a name (including the hash))

* `combinator description` is a string of format:

    combinator_name type_arg_1 ... type_arg_N = type_res

  where N stands for arity of the combinator. type_res is the combinator value type.

* `constructor` is a combinator that cannot be computed. This is used to represent composite data
  types. E.g.:

    int_tree IntTree int IntTree = IntTree

  Alongside combinator `empty_tree = Inttree`, may be used to define a composite data type called
  `IntTree` that takes on values in the form of binary trees with integers as nodes.

* `function (functional combinator)` is a combinator which may be computed on condition that the
  requisite number of arguments of requisite types are provided. The result of the computation is
  an expression consisting of constructors and base type values only.

* `normal form` is an expression consisting only of constructors and base type values; that which
  is normally the result of computing a function.

* `type identifier` is an identifier that normally starts with a capital letter in Roman script
  and uniquely identifies the type.

* `type number/name` is a 32-bit number that uniquely identifies a type; it normally is the sum of
  CRC32 values of the descriptions of the type constructors.

* `description of Type T` is a collection of the descriptions of all constructors that take on `Type
  T` values. Description of Type `IntTree`:

    int_tree IntTree int IntTree = IntTree;
    empty_tree = IntTree;

* `polymorphic type` is a type whose description contains parameters (type variables) in lieu of
  actual types.

    cons {alpha:Type} alpha (List alpha) = List alpha;
    nil {alpha:Type} = List alpha;

  Above is the description of Type `List alpha` where `List` is a polymorphic type of `arity 1` and
  `alpha` is a type variable which appears as the constructors **optional parameter (in curly
  braces)**.

* `value of Type T`. For example. Let `Combinator int_tree` have the index number 17, where as
  Combinator `empty_tree` has the index number 239. Then the value of Type `IntTree` is, for
  example `17 17 239 1 239 2 239` which is more conveniently written as:

    int_tree int_tree empty_tree 1 empty_tree 2 empty_tree

  From the standpoint of a high-level language, this is:

    int_tree (int_tree (empty_tree) 1 (empty_tree)) 2 (empty_tree): IntTree

* `boxed type` is a typedef which starts with the constructor number. Since every constructor has a
  uniquely determined value type, the first number (id) in any boxed type uniquely defines its type.
  A boxed type identifier is always capitalized.

* `bare type` is a type whose values do not contain a constructor number, which is implied instead.
  For example, `3 4` is a value of the int_couple bare type, defined using

    int_couple int int = IntCouple

  The corresponding boxed type is `IntCouple`; if 404 is the constructor index number for
  `int_couple`, then 
  
    404 3 4
    
  is the value for the `IntCouple` boxed type which corresponds to the value of the bare type
  `int_couple`.

**Conceptually, only boxed types should be used everywhere.**

However, for speed and compactness, bare types have to be used (for instance, an array of 10,000
bare int values is 40,000 bytes long, whereas boxed Int values take up twice as much space;
therefore, when transmitting a large array of integer identifiers, say, it is more efficient to use
the Vector int type rather than Vector Int).  In addition, all base types (int, long, double,
string) are bare.

## General Syntax

A TL program consists of a stream of tokens (separated by whitespace).

`TL-program ::= constr-declarations { --- functions --- fun-declarations | --- types ---
constr-declarations }`

The constructor- and function declarations are nearly identical in their syntax (they are both
combinators)

```
constr-declarations ::= { declaration ; }
fun-declarations ::= { declaration ; }
```

TL consists of two sections separated by "---functions---"

  - declarations of built-in type and aggregate types
  - declared functions

If additional type declarations are required after functions have been declared, "---types---" is
used.

Each declaration ends with a semicolon.

Example:

```
user#d23c81a3 id:int first_name:string last_name:string = User;

---functions---
getUser#b0f732d5 int = User;
getUsers#2d84d5f5 (Vector int) = Vector User;  # this is called a `combinator`
```

There are 3 declarations in this tl file.

1. `user`, which has the following properties:
    - id:int
    - first_name:string
    - last_name:string
2. getUser function, accepts an integer (int), returns a User
3. getUsers functions, accepts a list of integers (Vector int), and returns a list of User's
   (Vector User)

`user` constructor has been explicitly assigned a hex number (d23c81a3). It is the CRC32
calculation of the signature `int ? = Int`

`getUsers` function declaration also assigned a hex number. It is the CRC32 of the string
`getUsers Vector int = Vector User`. (**all parenthesis are removed!**)

```go
h := hash.NewIEEE()
h.Write([]byte("user id:int first_name:string last_name:string = User")
fmt.Printf("%x\n", h.Sum32())
```

**combinator declaration**

```
combinator-decl ::= lc-ident-ns [ # hex-digit *8 ]  { opt-args } { args } = result-type ;
                    ^^^^^^^^^^^ ^^^^^^^^^^^^^^^^^^  ^^^^^^^^^^^^ ^^^^^^^^ ^ ^^^^^^^^^^^ ^
                    pong        #34777ec5(optional) {t:Type}     id:int   = Pong        ;

e.g. pong#34777ec5 {t:Type} id:int = Pong;
```

There are also “pseudo-declarations” that are allowed only to declare **built-in** types.

```
builtin-combinator-decl ::= full-combinator-id ? = boxed-type-ident ;

e.g.: int ? = Int;
```

Let us say that we need to represent users as triplets containing one integer
(user ID) and two strings (first and last names). The requisite data structure
is the triplet int, string, string which may be declared as follows:

    user int string string = User;

On the other hand, a group may be described by a similar triplet consisting of
a group ID, its name, and description:

    group int string string = Group;

For the difference between User and Group to be clear, it is convenient to
assign names to some or all of the fields:

    user id:int first_name:string last_name:string = User;
    group id:int title:string description:string = Group;

If the User type needs to be extended at a later time by having records with
some additional field added to it, it could be accomplished as follows:

    userv2 id:int unread_messages:int first_name:string last_name:string in_groups:vector int = User;

### RPC query example

`getUsers([2, 3, 4])`. This query will be serialized into a sequence of 32-bit integers as follows:

    0x2d84d5f5 0x1cb5c415 0x3 0x2 0x3 0x4
    ^^^^^^^^^^ ^^^^^^^^^^ ^^^ ^^^ ^^^ ^^^
     │         │          │   │   │   └─ 4 (arg3)
     │         │          │   │   └───── 3 (arg2)
     │         │          │   └───────── 2 (arg1)
     │         │          └───────────── len(args)
     │         └──────────────────────── Vector int
     └────────────────────────────────── getUsers

  - TL serialization yields sequences of 32-bit integers. Little Endian. Above query corresponds to
    following byte stream:

`f5d5842d 15c4b51c 03000000 02000000 03000000 04000000`

Response might look like:

`0x1cb5c415 0x3 0xd23c81a3 0x2 0x74655005 0x00007265 0x72615006 0x72656b 0xc67599d1 0x3 0xd23c81a3
0x4 0x686f4a04 0x6e 0x656f4403`

which corresponds to:

`[{"id":2,"first_name":"Peter", "last_name":"Parker"},{},{"id":4,"first_name":"John","last_name":"Doe"}]`

## Formal Grammar

### Character classes:

    comment:    //
    lc-letter:  a...z
    uc-letter:  A...Z
    digit:      0...9
    hex-digit:  digit | abcdef
    underscore: _
    letter:     lc-letter | uc-letter
    ident-char: letter | digit | underscore

### Identifiers and keywords:

    lc-ident:        lc-letter { ident-char }
    uc-ident:        uc-letter { ident-char }
    namespace-ident: lc-ident
    lc-ident-ns:     [ namespace-ident . ] lc-ident
    uc-ident-ns:     [ namespace-ident . ] uc-ident
    lc-ident-full:   lc-ident-ns [ # hex-digit *8 ]

### Tokens:

    underscore:      _
    colon:           :
    semicolon:       ;
    open-par:        (
    close-par:       )
    open-bracket:    [
    close-bracket:   ]
    open-brace:      {
    close-brace:     }
    triple-minus:    ---
    equals:          =
    hash:            #
    question-mark:   ?
    percent:         %
    plus:            +
    langle:          <
    rangle:          >
    comma:           ,
    dot:             .
    asterisk:        *
    excl-mark:       !
    Final-kw:        Final
    New-kw:          New
    Empty-kw:        Empty
    nat-const:       digit { digit }
    lc-ident-full
    lc-ident
    uc-ident-ns

`Final`, `New` and `Empty` are special tokens (reserved word).

### Combinators

    combinator-decl ::= full-combinator-id { opt-args } { args } `=` result-type `;`
    full-combinator-id ::= lc-ident-full | `_`
    combinator-id ::= lc-ident-ns | `_`
    opt-args ::= `{` var-ident { var-ident } : [excl-mark] type-expr `}`
    args ::= var-ident-opt `:` [ conditional-arg-def ] [ `!` ] type-term
    args ::= [ var-ident-opt `:` ] [ multiplicity `*`] `[` { args } `]`
    args ::= `(` var-ident-opt { var-ident-opt } `:` [`!`] type-term `)`
    args ::= [ `!` ] type-term
    multiplicity ::= nat-term
    var-ident-opt ::= var-ident | `_`
    conditional-arg-def ::= var-ident [ `.` nat-const ] `?`
    result-type ::= boxed-type-ident { subexpr }
    result-type ::= boxed-type-ident `<` subexpr { `,` subexpr } `>`

* A combinator identifier is either an identifier starting with a lowercase Latin letter (lc-ident)
  or a namespace identifier (also lc-ident) followed by a period and another lc-ident. Therefore
  `cons` and `lists.get` are valid combinator identifiers.

* A combinator has a name, also knowns as a number, -a 32-bit number that unambiguously determines
  it. It is either calculated automatically or it is explicitly assigned in the declarations. H
  hash mark (#) and exactly 8 hex digits are added to the identifier of the combinator being
  defined.

* A combinators declaration begins with its identifier, to which its name (separated by a hash
  mark) may have been added.

* After the combinator identifier comes the main part of the declaration, which consists of
  declarations of `fields` (or variables), including an indication of their `types`.

* First come declarations of optional fields (of which there may be several or none at all). Then
  there are the declarations of the required fields (there may not be any of these either).

* Any identifier that begins with an uppercase or lowercase letter and which does not contain
  references to a namespace can be a field (variable) identifier. Using uc-ident for identifiers of
  variable types and lc-ident for other variables is a good practice.

* A combinator declaration contains the equals sign (=) and the result type (it may be composite or
  appearing for the first time). The result type may be polymorphic and/or dependent; any fields of
  the defined constructor's fields of type `Type` or `#` may be returned.

* A combinator declaration is terminated with a semicolon (;).

* A constructor's `fields`, `variables` and `arguments` mean the same thing.

## Tokenization

  - Each scan returns a single token.
