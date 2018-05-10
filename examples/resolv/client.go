package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ubiquiti/zeroconf"
)

var (
	service  = flag.String("service", "_workstation._tcp", "Set the service category to look for devices.")
	domain   = flag.String("domain", "local", "Set the search domain. For local networks, default is fine.")
	waitTime = flag.Int("wait", 10, "Duration in [s] to run discovery.")
)

func main() {
	flag.Parse()

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	offline := make(chan *zeroconf.ServiceEntry)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("online: ", entry)
		}
		log.Println("No more online entries.")
	}(entries)

	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			log.Println("offline: ", entry)
		}
		log.Println("No more offline entries.")
	}(offline)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(*waitTime))
	defer cancel()
	err = resolver.BrowseWithEx(ctx, *service, *domain, entries, offline, nil)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	time.Sleep(1 * time.Second)
}
