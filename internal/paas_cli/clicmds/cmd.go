package clicmds

import (
	"github.com/urfave/cli/v3"
)

var RootCmd = &cli.Command{
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Usage:   "Enable verbose output",
			Aliases: []string{"v"},
		},
	},
	Commands: []*cli.Command{
		AppCmd,
		EnvCmd,
	},
}
