package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/luoying/GoDistributed/registry"
)

func main() {
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var serv http.Server
	serv.Addr = registry.ServicePort

	go func() {
		log.Println(serv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Println("Registry service has been started. Press any key to terminate it")
		var s string
		fmt.Scanln(&s)
		serv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("Shutting down registry service")
}
