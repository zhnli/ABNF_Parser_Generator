package main

import (
	"fmt"
	"log"
	"strconv"
)

type (
	Expression struct {
		name  string
		value int
	}
)

func parse(s string) int {
	var v int
	var err error

	pos := 0
	buf := make([]byte, 0)

	var data Expression
	goto zExpression_Name_S0

zExpression_Name_S0:
	c := s[pos]
	pos += 1
	if c >= 'a' && c <= 'z' {
		buf = append(buf, c)
		goto zExpression_Name_S1
	}
	goto __zStateBad__

zExpression_Name_S1:
	c = s[pos]
	pos += 1
	if c >= 'a' && c <= 'z' {
		buf = append(buf, c)
		goto zExpression_Name_SS
	}
	goto __zStateBad__

zExpression_Name_SS:
	data.name = string(buf)
	buf = []byte{}
	goto zExpression_P2_S0

zExpression_P2_S0:
	c = s[pos]
	pos += 1
	if c == ':' {
		//buf = append(buf, c)
		goto zExpression_VALUE_S0
	}
	goto __zStateBad__

zExpression_VALUE_S0:
	c = s[pos]
	pos += 1
	if c >= '0' && c <= '9' {
		buf = append(buf, c)
		goto zExpression_VALUE_S1
	}
	goto __zStateBad__

zExpression_VALUE_S1:
	c = s[pos]
	pos += 1
	if c >= '0' && c <= '9' {
		buf = append(buf, c)
		goto zExpression_VALUE_SS
	}
	goto __zStateBad__

zExpression_VALUE_SS:
	v, err = strconv.Atoi(string(buf))
	if err != nil {
		log.Fatal("Error")
	}
	data.value = v
	buf = []byte{}
	goto zExpression_SS

zExpression_SS:
	fmt.Printf("Result: %#v\n", data)
	return 1

__zStateBad__:
	return -1
}

func main() {
	str := "ab:99"
	parse(str)
}
