package main

type Command struct {
}

func parseCommand(msg []byte) (Command, error) {
	cp := msg
	return Command{}, nil
}
