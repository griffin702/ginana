package commands

import (
	"github.com/urfave/cli/v2"
)

func InitCommands() []*cli.Command {
	f = new(Flags)
	return []*cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "创建新项目",
			Flags:   f.ToNewAction(),
			Action:  NewAction(),
		},
		{
			Name:   "build",
			Usage:  "ginana build",
			Action: BuildAction,
		},
		{
			Name:   "run",
			Usage:  "ginana run",
			Action: RunAction,
		},
	}
}
