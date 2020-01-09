package main

import (
	"apiserver/router"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	//创建gin引擎
	g := gin.New()

	/*问题1:这里为什么要使用New?而不是使用gin.Default()

	通过源码分析可以看到：

	1.gin.Default() 函数会生成一个默认的 Engine 对象里面包含了 2 个默认的常用插件
	分别是 Logger 和 Recovery
	Logger 用于输出请求日志
	Recovery 确保单个请求发生 panic 时记录异常堆栈日志，输出统一的错误响应

	创建带有默认中间件的路由:日志与恢复中间件

	2.gin.New() 函数只生成一个默认的engine对象, 创建不带中间件的路由,所以gin.New()更精简
	*/

	//创建gin中间件
	middlewares := []gin.HandlerFunc{}
	log.Print(middlewares)

	/* 问题2: 中间件为啥这样写？
	[]xxx{} 是golang中的slice数组写法

	list := []int{1, 2, 3, 4}
	fmt.Println(list)
	*/

	//注册路由-将自己写的路由注册进去gin
	router.Load(
		g,
		middlewares...,
	)
	/* 问题3 这里为什么要用...? 不加...有什么问题？

	我们看到Load函数是这样写的, 参数里面也有...
	func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
		xxxx
	}

	如果最后一个函数参数的类型的是...T，那么在调用这个函数的时候，
	我们可以在参数列表的最后使用若干个类型为T的参数。
	这里...T在函数内部的类型实际是[]T.

	*/

	/* 问题4： 这里为什么要使用go 协程方式?

	使用go 协程来完成程序的自检过程,
	通过自检可以最大程度地保证启动后的 API 服务器处于健康状态

	*/
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Print("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", ":8080")
	log.Printf(http.ListenAndServe(":8080", g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < 2; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get("http://127.0.0.1:8080" + "/api/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
