%{
package tmsh

import (
	"fmt"
)

const (
	ltmNodeNode = iota
	ltmPoolNode
	ltmVirtualNode
	keyNode
	structNode
	scalarNode
)

type nodeType int

type node struct {
	kind    nodeType
	value    string
	children []node
}

type Token struct {
	token   int
	literal string
}

%}

%union{
	token   Token
	ltm     node
	object  node
	pair    node
	members []node
	value   node
}

%type<ltm>     ltm
%type<object>  object
%type<pair>    pair
%type<members> members
%type<value>   value

%token<token> ILLEGAL
%token<token> EOF
%token<token> WS
%token<token> NEWLINE
%token<token> L_BRACE
%token<token> R_BRACE
%token<token> IDENT
%token<token> LTM

%%

ltm
	: LTM IDENT IDENT object
	{
		var kind nodeType
		switch $2.literal {
		case "node":
			kind = ltmNodeNode
		case "pool":
			kind = ltmPoolNode
		case "virtual":
			kind = ltmVirtualNode
		}
		yylex.(*Lexer).result = node{
			kind: kind,
			value: $3.literal,
			children: []node{$4},
		}
	}

object
	: L_BRACE R_BRACE
	{
		$$ = node{kind: structNode, value: "", children: []node{}}
	}
	| L_BRACE NEWLINE R_BRACE
	{
		$$ = node{kind: structNode, value: "", children: []node{}}
	}
	| L_BRACE NEWLINE members R_BRACE
	{
		$$ = node{kind: structNode, value: "", children: $3}
	}

members
  : pair
	{
		$$ = []node{$1}
	}
  | members pair
	{
		$$ = append($1, $2)
	}

pair
  : IDENT value NEWLINE
	{
		$$ = node{kind: keyNode, value: $1.literal, children: []node{$2}}
	}
	| IDENT object NEWLINE
	{
		$$ = node{kind: keyNode, value: $1.literal, children: []node{$2}}
	}

value
	: IDENT value
	{
		s := fmt.Sprintf("%s %s", $1.literal, $2.value)
		$$ = node{kind: scalarNode, value: s, children: []node{}}
	}
	| IDENT
	{
		$$ = node{kind: scalarNode, value: $1.literal, children: []node{}}
	}

%%

type Lexer struct {
	s      *Scanner
	result node
}

func (l *Lexer) Lex(lval *yySymType) int {
	var tok int
	var lit string

	for {
		tok, lit = l.s.Scan()
		if tok != WS {
			break
		}
	}

	lval.token = Token{token: tok, literal: lit}

	if tok == EOF {
		return 0
	}

	return tok
}

func (l *Lexer) Error(e string) {
	panic(e)
}
