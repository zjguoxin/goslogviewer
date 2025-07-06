package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zjguoxin/goslogviewer"
	"github.com/zjguoxin/goslogviewer/adapter"
)

func main() {
	r := gin.Default()

	// 配置日志查看器
	config := &goslogviewer.Config{
		LogDir:              "./log",
		EnableDelete:        true, // 开启删除日志文件
		EnableExport:        true, // 开启日志导出
		EnableClear:         true, // 开启日志清空
		DevMode:             true, // 开启开发模式
		EnableIPRestriction: true, // 开启IP限制
		AllowedIPs: []string{
			"192.168.1.0/24", // 允许整个子网
			"10.0.0.5",       // 允许特定IP
			"127.0.0.1",      // 允许本地访问
		},
		TrustedProxies: []string{"172.16.0.0/12"}, // 内部代理网络
	}
	lv := goslogviewer.New(config)

	// 注册Gin路由
	adapter.RegisterGinRoutes(r, lv)

	r.Run(":8080")

}
