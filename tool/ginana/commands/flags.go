package commands

import "github.com/urfave/cli/v2"

type Flags struct {
	Name      string
	Path      string
	ModPrefix string
}

var f *Flags

func (f *Flags) ToNewAction() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "d",
			Value:       "",
			Usage:       "指定项目所在目录",
			Destination: &f.Path,
		},
	}
}
