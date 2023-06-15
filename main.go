package main

import (
	"fmt"
	"os"

	"github.com/kanadeishii/void/internal/logger"
	"github.com/kanadeishii/void/pkg/renderer"
	"github.com/urfave/cli"
)

func Handle(context *cli.Context) {
	link := context.String("link")

	fmt.Print("\n")
	logger.Rich(logger.ColorBlue, "void", "requesting to "+link+"\n")

	renderer.Render(link, true)
}

func main() {
	app := &cli.App{
		Name:  "void",
		Usage: "generate videos on webdev",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "link",
				Usage: "access link like http://localhost:5500/foo/app.html",
			},
		},
		Action: Handle,
	}

	app.Run(os.Args)
}
