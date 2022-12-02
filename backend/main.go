package main

import (
	"backend/config"
	"backend/router"
	"log"
)

func main() {
	addr, ok := config.Get("HTTP_ADDR")
	if !ok {
		log.Fatalln("HTTP_ADDR env not found.")
	}

	app := router.App()
	log.Fatalln(app.Listen(addr))
}
