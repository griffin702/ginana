package commands

import "time"

var toolIndexs = []*Tool{
	{
		Name:      "ginana",
		Alias:     "ginana",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/griffin702/ginana/tool/ginana@" + Version,
		Summary:   "工具集",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "ginana",
	},
	{
		Name:      "wire",
		Alias:     "wire",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/google/wire/cmd/wire",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "google",
	},
	{
		Name:      "packr2",
		Alias:     "packr2",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/gobuffalo/packr/v2/packr2",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "gobuffalo",
	},
	{
		Name:      "swag",
		Alias:     "swag",
		BuildTime: time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local),
		Install:   "go get -u github.com/swaggo/swag/cmd/swag",
		Platform:  []string{"darwin", "linux", "windows"},
		Author:    "swaggo",
	},
}
