package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	m map[string]func(*state, command) error
}

func (c *commands) register(name string, handler func(*state, command) error) {
	c.m[name] = handler
}

func (c *commands) run(state *state, cmd command) error {
	handler, ok := c.m[cmd.name]
	if !ok {
		return fmt.Errorf("command %s not found", cmd.name)
	}

	return handler(state, cmd)
}
