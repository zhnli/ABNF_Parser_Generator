package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var TrimSet string = " \t\n"

func doParse(buf []string) (string, *interface{}) {
	str := strings.Join(buf, " ")
	pos := strings.Index(str, "=")

	if pos < 0 {
		log.Fatalf("Error: '=' missing. Str='%v'", str)
	}

	statName := strings.Trim(str[:pos], TrimSet)
	str = str[pos+1:]
	str = strings.Trim(str, TrimSet)
	fmt.Printf("Parsing string \"%v\"\n", str)

	l := yyLex{s: str}

	if yyParse(&l) != 0 {
		log.Fatal(l.errs)
	}

	fmt.Printf("AST:%#v\n", _parseResult)
	return statName, &_parseResult
}

func removeInlineComment(text string) string {
	pos := strings.Index(text, " ;")
	if pos != -1 {
		text = text[:pos]
	}
	pos = strings.Index(text, "\t;")
	if pos != -1 {
		text = text[:pos]
	}
	return text
}

func main() {
	yyDebug = 0
	yyErrorVerbose = true

	debug := flag.Int(
		"d", 0, "Debug level [1-4]")
	infile := flag.String(
		"i", "syntax.def", "Input file that contains syntax definition in ABNF")
	outfile := flag.String(
		"o", "ast.output", "Output file that contains AST parsing results")

	flag.Parse()
	yyDebug = *debug

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Printf("Parsing file \"%v\"\n", *infile)

	infh, err := os.Open(*infile)
	if err != nil {
		log.Fatal(err)
	}
	defer infh.Close()

	outfh, err := os.Create(*outfile)
	if err != nil {
		log.Fatal(err)
	}
	defer outfh.Close()

	var buf []string
	var isNewStat bool = false
	trimStr := " \t\n"

	astMap := make(map[string]interface{})
	writer := bufio.NewWriter(outfh)
	scanner := bufio.NewScanner(infh)
	for scanner.Scan() {
		str := scanner.Text()
		strLen := len(str)

		str = strings.TrimLeft(str, trimStr)
		contentLen := len(str)

		if contentLen == 0 || str[0] == ';' {
			continue
		}

		if contentLen == strLen {
			isNewStat = true
		}

		str = strings.TrimRight(str, trimStr)
		str = removeInlineComment(str)

		if isNewStat && len(buf) > 0 {
			name, ast := doParse(buf)
			astMap[name] = *ast
			buf = []string{}
			isNewStat = false
		}

		buf = append(buf, str)
	}

	name, ast := doParse(buf)
	astMap[name] = *ast

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	result := fmt.Sprintf("AST:%#v\n", astMap)
	writer.WriteString(result)
	writer.Flush()
}
