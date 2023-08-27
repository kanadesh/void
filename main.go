package main

import (
	"fmt"
	"os"

	"github.com/kanadesh/void/internal/logger"
	"github.com/kanadesh/void/pkg/renderer"
	"github.com/urfave/cli"
)

func Handle(context *cli.Context) {
	link := context.String("link") // Get link from params

	// Put informations
	fmt.Print("\n")
	logger.Rich(logger.ColorBlue, "void", "requesting to "+link+"\n")

	renderer.Render(link, true) // Render screenshots with Chromedp
}

func main() {
	// Initialize the application
	app := &cli.App{
		Name:  "void",
		Usage: "generate videos on webdev",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "link",
				Usage: "Access link like http://localhost:5500/foo/app.html",
			},
		},
		Action: Handle,
	}

	app.Run(os.Args) // Run the app
}
