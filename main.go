package main

import (
	"flag"
	"log"
	"novpn/client"
	"novpn/config"
	"novpn/exchange"
	"novpn/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := flag.String("f", "config.json", "Config file")
	flag.Parse()
	psignal := make(chan os.Signal, 1)
	// ctrl+c->SIGINT, kill -9 -> SIGKILL
	signal.Notify(psignal, syscall.SIGINT, syscall.SIGKILL)
	config, err := config.FromJSONFile(*cfg)
	if err != nil {
		log.Println(err)
		return
	}
	go exchange.Run(config.Exchange)
	time.Sleep(time.Second)
	go service.Run(config.Service)
	time.Sleep(time.Second)
	go client.Run(config.Client)
	<-psignal
	log.Println("Bye~")
}
