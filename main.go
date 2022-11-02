package main

import (
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("debugPProf") == "true" {
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	app := cli.NewApp()

	app.Name = "alyx_nft_backend"
	app.Version = "v0.1.0"
	app.Description = "alyx_nft_backend"
	server := NewService()
	app.Action = server.Start

	_ = app.Run(os.Args)
}
