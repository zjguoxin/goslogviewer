/**
 * @Author: guxline zjguoxin@163.com
 * @Date: 2025/7/3 07:17:54
 * @LastEditors: guxline zjguoxin@163.com
 * @LastEditTime: 2025/7/3 07:17:54
 * Description: Gin适配器
 * Copyright: Copyright (©) 2025 中易综服. All rights reserved.
 */
package adapter

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zjguoxin/goslogviewer"
	"github.com/zjguoxin/goslogviewer/middleware"
)

func RegisterGinRoutes(r *gin.Engine, lv *goslogviewer.LogViewer) {
	r.Static("/static", "./adapter/templates/source")

	// 1. 从嵌入的文件系统加载模板
	tmpl := template.Must(
		template.New("").
			ParseFS(templateFS, "templates/log.html"),
	)

	r.SetHTMLTemplate(tmpl)

	group := r.Group("/log", middleware.IPRestriction(
		lv.GetConfig().EnableIPRestriction,
		lv.GetConfig().AllowedIPs,
		lv.GetConfig().TrustedProxies,
	))
	{
		group.GET("", func(c *gin.Context) {
			c.HTML(http.StatusOK, "log.html", gin.H{
				"head":  "日志查看器",
				"title": "日志查看器",
			})
		})
		group.GET("/getLogFilesList", func(c *gin.Context) { lv.GetFilesHandler(c.Writer, c.Request) })
		group.GET("/getFileContent", func(c *gin.Context) { lv.GetContentHandler(c.Writer, c.Request) })
		group.POST("/clearFileContent", func(c *gin.Context) { lv.ClearFileContentHandler(c.Writer, c.Request) })
		group.POST("/deleteAllFiles", func(c *gin.Context) { lv.DeleteAllFilesHandler(c.Writer, c.Request) })
		group.GET("/exportFile", func(c *gin.Context) { lv.ExportFileHandler(c.Writer, c.Request) })

	}
}
