package main

import (
	"errors"

	"github.com/spf13/cobra"
)

type cmdImport struct {
	global *cmdGlobal
}

func (c *cmdImport) command() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "import"
	cmd.Short = `Command has been replaced with "lxd recover"`
	cmd.Long = `Description:
  This command has been replaced with "lxd recover". Please use that instead.
`
	cmd.RunE = c.run
	return cmd
}

func (c *cmdImport) run(cmd *cobra.Command, args []string) error {
	return errors.New(`Command has been replaced with "lxd recover"`)
}
