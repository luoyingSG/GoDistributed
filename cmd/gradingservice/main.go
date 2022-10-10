package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/luoying/GoDistributed/grades"
	"github.com/luoying/GoDistributed/registry"
	"github.com/luoying/GoDistributed/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)

	r := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(context.Background(),
		host,
		port, r, grades.RegisterHandlers)
	if nil != err {
		stlog.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
