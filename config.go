/**
 * @Author: guxline zjguoxin@163.com
 * @Date: 2025/7/3 07:12:26
 * @LastEditors: guxline zjguoxin@163.com
 * @LastEditTime: 2025/7/3 07:12:26
 * Description: 配置结构体
 * Copyright: Copyright (©) 2025 中易综服. All rights reserved.
 */
package goslogviewer

type Config struct {
	DevMode             bool     // 是否开发模式
	LogDir              string   // 日志目录路径
	EnableIPRestriction bool     // 是否启用IP限制（默认false）
	TrustedProxies      []string // 可信代理IP列表（用于获取真实客户端IP）
	AllowedIPs          []string // 允许访问的IP列表
	EnableDelete        bool     // 是否启用删除功能
	EnableExport        bool     // 是否启用导出功能
	EnableClear         bool     // 是否启用清除功能
	PageSize            int      // 每页显示条数
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		LogDir:       "./log",
		DevMode:      false,
		EnableDelete: false,
		EnableExport: true,
		EnableClear:  false,
		PageSize:     10,
		AllowedIPs:   []string{"127.0.0.1"},
	}
}
