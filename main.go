package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/cadenmarchese/bottlecap/pkg/client"
)

var supportedArgs []string = []string{"ask", "image", "generate"}

func main() {
	app := &cli.Command{
		Usage: `Provide inputs to the bottelecap application using quotes`,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			args := os.Args[1:]

			if len(args) < 2 || !strings.Contains(strings.Join(supportedArgs, ""), args[0]) {
				return fmt.Errorf(`usage: ask "Your question in quotes" - or, image "your Image URL" - or, generate "your image description"`)
			}

			subcommand := args[0]
			argument := args[1]
			response, err := client.Client(subcommand, argument)
			if err != nil {
				return err
			}

			fmt.Println(response)
			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
