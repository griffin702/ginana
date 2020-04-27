package commands

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

func RunNew() cli.ActionFunc {
	return func(ctx *cli.Context) (err error) {
		for _, t := range toolList() {
			if !t.installed() || t.needUpdated() {
				fmt.Println(t.Install)
				t.install()
			}
		}
		f.Name = ctx.Args().First()
		if f.Path != "" {
			if f.Path, err = filepath.Abs(f.Path); err != nil {
				return
			}
			f.Path = filepath.Join(f.Path, f.Name)
		} else {
			pwd, _ := os.Getwd()
			f.Path = filepath.Join(pwd, f.Name)
		}
		f.ModPrefix = modPath(f.Path)
		// creata a project
		if err := create(); err != nil {
			return err
		}
		fmt.Printf("Project: %s\n", f.Name)
		fmt.Printf("Directory: %s\n\n", f.Path)
		fmt.Println("项目创建成功.")
		return nil
	}
}
