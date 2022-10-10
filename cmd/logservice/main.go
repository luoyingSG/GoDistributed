package main

import (
	"context"
	"fmt"
	stdlog "log"

	"github.com/luoying/GoDistributed/log"
	"github.com/luoying/GoDistributed/registry"
	"github.com/luoying/GoDistributed/service"
)

func main() {
	log.Run("./distributed.log") // 指定 log 文件的位置
	host, port := "localhost", "4000"
	reg := registry.Registration{
		ServiceName: registry.LogService,
		ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		reg,
		log.RegisterHandlers,
	)
	if nil != err {
		stdlog.Fatalln(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
