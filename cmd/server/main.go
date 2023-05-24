package main

import "github.com/Eviljeks/blurer/cmd/server/app"

func main() {
	cfg := app.DefaultConfig()

	cfg.Run()
}
