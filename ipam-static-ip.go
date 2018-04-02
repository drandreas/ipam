package main

import (
	"github.com/docker/go-plugins-helpers/ipam"
	"github.com/drandreas/ipam-static-ip/handler"
	"log"
)

func main() {
	log.Printf("Starting Up...")
	d := handler.NewHandler()
	h := ipam.NewHandler(d)
	err := h.ServeUnix("ipam-static-ip", 993)
	if err != nil {
		log.Printf("Error: ", err)
	}
}
