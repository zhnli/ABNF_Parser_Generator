# Copyright 2014 The ebnf2y Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY: all clean

all: main.go abnf.y.go
	go fmt
	go build

clean:
	@go clean
	rm -f y.output

abnf.y.go: abnf.y
	go tool yacc -o $@ $<

