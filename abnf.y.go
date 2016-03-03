//line abnf.y:8
package main

import __yyfmt__ "fmt"

//line abnf.y:9
import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var _parseResult interface{}

//line abnf.y:24
type yySymType struct {
	yys int
	val interface{}
}

const WSP = 57346
const NUMBER = 57347
const CHARSET = 57348
const STRING = 57349
const IDENT = 57350
const UNITE = 57351
const JOIN = 57352
const QUANT = 57353

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"WSP",
	"NUMBER",
	"CHARSET",
	"STRING",
	"IDENT",
	"UNITE",
	"'/'",
	"JOIN",
	"QUANT",
	"'*'",
	"'('",
	"')'",
	"'['",
	"']'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line abnf.y:81

/*  start  of  programs  */

type yyLex struct {
	s    string
	pos  int
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
	v   []interface{}
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
	return Exp{opt: OPT_JOIN, v: []interface{}{a, b}}
}

func Alter(a interface{}, b interface{}) Exp {
	return Exp{opt: OPT_ALTER, v: []interface{}{a, b}}
}

func Quant(a interface{}, b interface{}) Exp {
	return Exp{opt: OPT_QUANT, v: []interface{}{a, b}}
}
func Ident(a interface{}) Exp {
	return Exp{opt: OPT_IDENT, v: []interface{}{a}}
}

func Str(a interface{}) Exp {
	return Exp{opt: OPT_STRING, v: []interface{}{a}}
}

func CharSet(v interface{}) Exp {
	v_str, ok := v.(string)
	if !ok {
		return Exp{opt: OPT_CHARSET, v: []interface{}{0}}
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

	return Exp{opt: OPT_CHARSET, v: []interface{}{vint[0], vint[1]}}
}

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 2,
	1, 1,
	-2, 14,
}

const yyNprod = 16
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 80

var yyAct = [...]int{

	17, 2, 16, 22, 15, 5, 14, 1, 0, 0,
	0, 0, 0, 20, 21, 0, 19, 10, 8, 7,
	6, 0, 18, 0, 26, 9, 3, 11, 4, 28,
	0, 12, 13, 19, 10, 8, 7, 6, 0, 18,
	0, 0, 9, 3, 27, 4, 23, 0, 24, 25,
	19, 10, 8, 7, 6, 0, 18, 0, 0, 9,
	3, 0, 4, 19, 10, 8, 7, 6, 10, 8,
	7, 6, 9, 3, 0, 4, 9, 3, 0, 4,
}
var yyPact = [...]int{

	63, -1000, -1000, -1000, -1000, 63, -1000, -1000, -1000, -1,
	-11, 46, 59, 59, -1000, -1000, -2, -1000, -1000, -1000,
	-1000, -1000, -1000, 59, 29, 12, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 7, 0, 5, 27,
}
var yyR1 = [...]int{

	0, 1, 2, 2, 2, 2, 2, 2, 2, 2,
	3, 3, 3, 3, 4, 4,
}
var yyR2 = [...]int{

	0, 1, 5, 5, 2, 3, 5, 1, 1, 1,
	1, 2, 2, 3, 0, 2,
}
var yyChk = [...]int{

	-1000, -1, -2, 14, 16, -3, 8, 7, 6, 13,
	5, -4, -4, -4, -2, 5, 13, -2, 10, 4,
	-2, -2, 5, -4, -4, -4, -2, 15, 17,
}
var yyDef = [...]int{

	0, -2, -2, 14, 14, 0, 7, 8, 9, 10,
	0, 0, 0, 0, 4, 12, 11, 5, 14, 15,
	14, 14, 13, 0, 0, 0, 6, 2, 3,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	14, 15, 13, 3, 3, 3, 3, 10, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 16, 3, 17,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 11, 12,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lookahead func() int
}

func (p *yyParserImpl) Lookahead() int {
	return p.lookahead()
}

func yyNewParser() yyParser {
	p := &yyParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yytoken := -1 // yychar translated into internal numbering
	yyrcvr.lookahead = func() int { return yychar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yychar = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar, yytoken = yylex1(yylex, &yylval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yychar = -1
		yytoken = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar, yytoken = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yychar = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line abnf.y:45
		{
			_parseResult = yyDollar[1].val
		}
	case 2:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line abnf.y:49
		{
			yyVAL.val = yyDollar[3].val
		}
	case 3:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line abnf.y:51
		{
			yyVAL.val = Quant(NewQt(0, 1), yyDollar[3].val)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line abnf.y:53
		{
			yyVAL.val = Quant(yyDollar[1].val, yyDollar[2].val)
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line abnf.y:55
		{
			yyVAL.val = Join(yyDollar[1].val, yyDollar[3].val)
		}
	case 6:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line abnf.y:57
		{
			yyVAL.val = Alter(yyDollar[1].val, yyDollar[5].val)
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line abnf.y:59
		{
			yyVAL.val = Ident(yyDollar[1].val)
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line abnf.y:61
		{
			yyVAL.val = Str(yyDollar[1].val)
		}
	case 9:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line abnf.y:63
		{
			yyVAL.val = CharSet(yyDollar[1].val)
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line abnf.y:67
		{
			yyVAL.val = NewQt(-1, -1)
		}
	case 11:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line abnf.y:69
		{
			yyVAL.val = NewQt(yyDollar[1].val, -1)
		}
	case 12:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line abnf.y:71
		{
			yyVAL.val = NewQt(-1, yyDollar[2].val)
		}
	case 13:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line abnf.y:73
		{
			yyVAL.val = NewQt(yyDollar[1].val, yyDollar[3].val)
		}
	case 14:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line abnf.y:77
		{
			yyVAL.val = ""
		}
	case 15:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line abnf.y:79
		{
			yyVAL.val = yyDollar[1].val
		}
	}
	goto yystack /* stack new state and value */
}
