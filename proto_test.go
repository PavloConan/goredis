package main

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	raw := "*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	_, err := parseCommand(raw)

	if err != nil {
		t.Fatal(err)
	}
}
