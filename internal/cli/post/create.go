package post

import (
	"log"

	"github.com/urfave/cli/v2"
)

var create = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "List Posts with specified query.",
	Action: func(c *cli.Context) error {
		log.Println("create")
		return nil
	},
}
