package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/tidwall/resp"
)

const (
	CommandSET = "SET"
	CommandGET = "GET"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type GetCommand struct {
	key []byte
}

func parseCommand(raw string) (Command, error) {
	var cmd Command

	rd := resp.NewReader(bytes.NewBufferString(raw))
	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, val := range v.Array() {
				switch val.String() {
				case CommandGET:
					if len(v.Array()) != 2 {
						return nil, errors.New("invalid number of arguments for GET command")
					}
					cmd = GetCommand{
						key: v.Array()[1].Bytes(),
					}

				case CommandSET:
					if len(v.Array()) != 3 {
						return nil, errors.New("invalid number of arguments for SET command")
					}
					cmd = SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
				}
				return cmd, nil
			}
		}
	}

	return nil, fmt.Errorf("invalid or unknown command: %s", raw)
}
