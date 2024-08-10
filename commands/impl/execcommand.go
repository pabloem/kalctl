package impl

import (
	"github.com/pabloem/kalctl/commands/base"
)

type customRunCommand struct {
	name string
	desc string
	run  func(args base.CommandArgs) error
	args []base.Argument
}

func (c *customRunCommand) Name() string {
	return c.name
}

func (c *customRunCommand) Description() string {
	return c.desc
}

func (c *customRunCommand) Run(args base.CommandArgs) error {
	return c.run(args)
}

func (c *customRunCommand) Arguments() []base.Argument {
	return nil
}

// TODO: Add support for arguments
func NewCustomRunCommand(name, desc string, run func(args base.CommandArgs) error, args ...base.Argument) base.Command {
	return &customRunCommand{
		name: name,
		desc: desc,
		run:  run,
	}
}
