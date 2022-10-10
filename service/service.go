package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/luoying/GoDistributed/registry"
)

// 用于开启服务
func Start(ctx context.Context, host, port string, reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {
	// 调用该服务的请求处理函数
	registerHandlersFunc()
	// 启动该服务
	ctx = startService(ctx, host, port, reg.ServiceName)
	// 注册该服务
	err := registry.RegisterService(reg)
	if nil != err {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, host, port string, serviceName registry.ServiceName) context.Context {
	// 可取消的服务
	ctx, cancel := context.WithCancel(ctx)

	// 在本地启动该服务
	var serv http.Server
	serv.Addr = ":" + port

	go func() {
		log.Println(serv.ListenAndServe())
		err := registry.Shutdown(fmt.Sprintf("http://%s:%s", host, port))
		if nil != err {
			log.Println(err)
		}
		cancel()
	}()

	go func() {
		// 用户可以输入任意字符来停止该服务
		fmt.Printf("%v has been started. Press any key to terminate it\n", serviceName)
		// 等待用户输入，如果没有输入，则会在这里暂等，不会停止服务
		var s string
		fmt.Scanln(&s)

		err := registry.Shutdown(fmt.Sprintf("http://%s:%s", host, port))
		if nil != err {
			log.Println(err)
		}

		serv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
