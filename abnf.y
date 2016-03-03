// Copyright 2011 Bobby Powers. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// based off of Appendix A from http://dinosaur.compilertools.net/yacc/

%{

package main

import (
    "fmt"
    "unicode"
    "strconv"
    "strings"
)

var _parseResult interface{}

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union{
    val interface{}
}

// any non-terminal which returns a value needs a type, which is
// really a field name in the above union struct
%type <val> stat expr quant vwsp

// same for terminals
%token <val> WSP NUMBER CHARSET STRING IDENT

%left UNITE
%left '/'
%left JOIN
%left QUANT
%left '*'

%start stat
%%

stat    :   expr
        { _parseResult = $1 }
    ;

expr    :    '(' vwsp expr vwsp ')' %prec UNITE
        { $$  =  $3 }
    |   '[' vwsp expr vwsp ']'
        { $$ = Quant(NewQt(0, 1), $3) }
    |    quant expr                %prec  QUANT
        { $$  =  Quant($1, $2) }
    |    expr vwsp expr            %prec  JOIN
        { $$  =  Join($1, $3) }
    |    expr vwsp '/' vwsp expr
        { $$  =  Alter($1, $5) }
    |   IDENT
        { $$  =  Ident($1) }
    |   STRING
        { $$  =  Str($1) }
    |   CHARSET
        { $$  =  CharSet($1) }
    ;

quant   :   '*'
        { $$ = NewQt(-1, -1) }
    |    NUMBER '*'
        { $$ = NewQt($1, -1) }
    |    '*' NUMBER
        { $$ = NewQt(-1, $2) }
    |    NUMBER '*' NUMBER
        { $$ = NewQt($1, $3) }
    ;

vwsp    : /* empty */
        { $$ = "" }
    |   vwsp WSP
        { $$ = $1 }

%%      /*  start  of  programs  */

type yyLex struct {
    s string
    pos int
    errs string
}

const (
    OPT_STRING = iota
    OPT_IDENT
    OPT_CHARSET
    OPT_QUANT
    OPT_JOIN
    OPT_ALTER
)

type QuantParam struct {
    min int
    max int
}

type Exp struct {
    opt int
    v []interface{}
}

func (l *yyLex) Lex(lval *yySymType) int {
    if l.pos == len(l.s) {
        return 0
    }

    var is_break bool = false
    c := rune(l.s[l.pos])
    l.pos += 1

    if c == ' ' || c == '\t' || c == '\n' {
        for c == ' ' || c == '\t' || c == '\n' {
            if l.pos == len(l.s) {
                return 0
            }
            c = rune(l.s[l.pos])
            l.pos += 1
        }

        l.pos -= 1
        return WSP
    }

    if c == '"' {
        if l.pos == len(l.s) {
            l.errs = "Quote missing"
            return 0
        } else {
            c = rune(l.s[l.pos])
            l.pos += 1
        }
        data := []byte{}
        for c != '"' {
            data = append(data, []byte(string(c))...)
            if l.pos == len(l.s) {
                l.errs = "Quote missing"
                return 0
            }
            c = rune(l.s[l.pos])
            l.pos += 1
        }
        lval.val = string(data)
        return STRING
    }

    // In ABNF, an identifier can start with a digit. i.e. "1UTF-CON"
    if unicode.IsDigit(c) || unicode.IsLetter(c) {
        isNumber := true
        data := []byte{}
        for unicode.IsDigit(c) || unicode.IsLetter(c) || c == '-' {
            if !unicode.IsDigit(c) {
                isNumber = false
            }
            data = append(data, []byte(string(c))...)
            if l.pos == len(l.s) {
                is_break = true
                break
            }
            c = rune(l.s[l.pos])
            l.pos += 1
        }
        if !is_break {
            l.pos -= 1
        }
        if isNumber {
            v_num, err := strconv.Atoi(string(data))
            if err != nil {
                l.errs = err.Error()
                return 0
            }
            lval.val = v_num
            return NUMBER
        } else {
            lval.val = string(data)
            return IDENT
        }
    }

    if c == '%' {
        if l.expect('x') != 0 {
            return 0
        }
        c = rune(l.s[l.pos])
        l.pos += 1
        data := []byte{}
        for (c >= '0' && c <= '9') ||
            (c >= 'A' && c <= 'F') ||
            c == '-' {

            data = append(data, []byte(string(c))...)
            if l.pos == len(l.s) {
                is_break = true
                break
            }
            c = rune(l.s[l.pos])
            l.pos += 1
        }
        if !is_break {
            l.pos -= 1
        }
        lval.val = string(data)
        return CHARSET
    }

    return int(c)
}

func (l *yyLex) expect(v rune) int {
    if l.pos == len(l.s) {
        l.errs = fmt.Sprintf("End of string. Expecting %s", v)
        return -1
    }
    c := rune(l.s[l.pos])
    if c == v {
        l.pos += 1
        return 0
    } else {
        l.errs = fmt.Sprintf("Error string. Expecting %s Pos=%d", v, l.pos)
        return -1
    }
}

func (l *yyLex) Error(s string) {
    fmt.Printf("Error=\"%s\" Pos=%d\n", s, l.pos)
}

func NewQt(min interface{}, max interface{}) QuantParam {

    v_min, ok := min.(int)
    if !ok {
        v_min = -1
    }

    v_max, ok := max.(int)
    if !ok {
        v_max = -1
    }

    return QuantParam{v_min, v_max}
}

func Join(a interface{}, b interface{}) Exp {
    return Exp{opt:OPT_JOIN, v:[]interface{}{a, b}}
}

func Alter(a interface{}, b interface{}) Exp {
    return Exp{opt:OPT_ALTER, v:[]interface{}{a, b}}
}

func Quant(a interface{}, b interface{}) Exp {
    return Exp{opt:OPT_QUANT, v:[]interface{}{a, b}}
}
func Ident(a interface{}) Exp {
    return Exp{opt:OPT_IDENT, v:[]interface{}{a}}
}

func Str(a interface{}) Exp {
    return Exp{opt:OPT_STRING, v:[]interface{}{a}}
}

func CharSet(v interface{}) Exp {
    v_str, ok := v.(string)
    if !ok {
        return Exp{opt:OPT_CHARSET, v:[]interface{}{0}}
    }
    values := strings.Split(v_str, "-")

    var vint []int
    ivalue, err := strconv.ParseInt(values[0], 16, 8)
    if err != nil {
        ivalue = 0
    }
    vint = append(vint, int(ivalue))

    if len(values) > 1 {
        ivalue, err = strconv.ParseInt(values[1], 16, 8)
        if err != nil {
            ivalue = 0
        }
        vint = append(vint, int(ivalue))
    } else {
        vint = append(vint, -1)
    }

    return Exp{opt:OPT_CHARSET, v:[]interface{}{vint[0], vint[1]}}
}
