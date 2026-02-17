package main

import (
	"context"
	"fkmcps/cmd"
	"log"
	"os"
)

func main() {
	app := cmd.NewApp()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
