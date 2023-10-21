package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/guyst16/mykube/pkg/virtualmachine"

	"github.com/urfave/cli"
)

var MAIN_DIR = "~/.mykube"
var LIBVIRT_MYKUBE_DIR = "/var/lib/libvirt/images/mykube"

func Cli() {
	app := &cli.App{
		Name:   "mykube",
		Usage:  "Manage single node K8S",
		Author: "guyst16 - Guy Steinberger",
		Commands: []cli.Command{
			{
				Name:  "create",
				Usage: "Create a single node K8S",
				Action: func(ctx *cli.Context) error {
					// Validate main directory existence
					_, dir_err := os.Stat(MAIN_DIR)
					if os.IsNotExist(dir_err) {
						os.Mkdir(MAIN_DIR, 0744)
					}

					// Validate virtual machines directory in Libvirt directory
					_, dir_err = os.Stat(MAIN_DIR)
					if os.IsNotExist(dir_err) {
						os.Mkdir(MAIN_DIR, 0744)
					}

					// Create mykube virtual machine
					myVM := virtualmachine.NewVirtualmachine("a", "a", 1, 1, "a")
					myVM.CreateVirtualmachine()
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "delete a single node K8S",
				Action: func(ctx *cli.Context) error {
					fmt.Printf("delete %s\n", ctx.Args().Get(0))
					return nil
				},
			},
			{
				Name:  "connect",
				Usage: "connect a single node K8S",
				Action: func(ctx *cli.Context) error {
					fmt.Printf("connect %s\n", ctx.Args().Get(0))
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "list all single nodes K8S",
				Action: func(ctx *cli.Context) error {
					virtualmachine.ListAllVirtualmachines()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
