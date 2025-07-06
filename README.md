# Go Slog Viewer - Gin 框架安全增强版日志查看器

# 目录导航

- [功能概览](#功能概览)
  - [核心功能](#核心功能)
  - [安全增强特性](#安全增强特性)
- [安装方式](#安装方式)
- [使用方法](#使用方法)
  - [快速配置](#快速配置)
  - [配置选项详解](#配置选项详解)
  - [IP 拒绝响应](#ip拒绝响应)
- [生产环境建议](#生产环境建议)
  - [必须配置项](#必须配置项)
  - [处理器接口](#处理器接口)
  - [响应格式](#响应格式)
- [集成示例](#集成示例)
- [许可证](#许可证)
- [作者](#作者)

# 功能概览

## 核心功能：

- 指定日志文件目录
- 获取日志文件列表
- 查看日志文件内容
- 清空指定日志文件
- 删除所有日志文件
- 导出日志文件(支持 JSON 和纯文本格式)
- 安全的日志管理（清空/删除）
- 细粒度 IP 访问控制
- 代理感知的真实 IP 获取

## 安全增强特性

- IP 白名单限制（支持 CIDR 表示法）
- 可信代理配置（自动处理 X-Forwarded-For）
- 本地开发友好（自动转换 IPv6 的::1 地址）
- 生产环境安全警告

# 安装方式

```bash
go get github.com/zjguoxin/goslogviewer
```

# 使用方法

## 快速配置

```go
config := &goslogviewer.Config{
    LogDir:       "./logs",
    EnableDelete:        true, // 开启日志删除
	  EnableExport:        true, // 开启日志导出
		EnableClear:         true, // 开启日志清空
    // 安全配置
    EnableIPRestriction: true,
    AllowedIPs: []string{
        "192.168.1.0/24",  // 内网段
        "10.0.0.5",        // 特定IP
        "127.0.0.1",       // 本地
    },
    TrustedProxies: []string{
        "172.16.0.0/12",   // 内部代理
    },
}
```

## 配置选项详解

配置项 类型 默认值 说明
EnableIPRestriction bool false 是否启用 IP 限制
AllowedIPs []string nil 允许访问的 IP 列表（支持 CIDR）
TrustedProxies []string nil 可信代理 IP 段（自动跳过 XFF 检查）

| 配置项              | 类型     | 默认值 | 说明                                |
| ------------------- | -------- | ------ | ----------------------------------- |
| EnableClear         | bool     | false  | 是否启用日志清空                    |
| EnableDelete        | bool     | false  | 是否启用日志删除                    |
| EnableExport        | bool     | false  | 是否启用日志导出                    |
| EnableIPRestriction | bool     | false  | 是否启用 IP 限制                    |
| AllowedIPs          | []string | nil    | 允许访问的 IP 列表（支持 CIDR）     |
| TrustedProxies      | []string | nil    | 可信代理 IP 段（自动跳过 XFF 检查） |

## <span id="ip拒绝响应">IP 拒绝响应</span>

```json
{
  "code": 403,
  "message": "Access denied for IP: 203.0.113.1",
  "your_ip": "203.0.113.1",
  "allowed": ["192.168.1.0/24"]
}
```

# 生产环境建议

## 必须配置项

```go
EnableIPRestriction = true
AllowedIPs = ["your_management_ip"]
```

## 处理器接口

| 处理器                  | HTTP 方法 | 功能描述             | 参数说明                                          |
| ----------------------- | --------- | -------------------- | ------------------------------------------------- |
| GetFilesHandler         | GET       | 获取可用日志文件列表 | 无参数                                            |
| GetContentHandler       | GET       | 获取指定日志文件内容 | `name` - 文件名                                   |
| ClearFileContentHandler | POST      | 清空指定日志文件     | `name` - 文件名（表单数据）                       |
| DeleteAllFilesHandler   | POST      | 删除所有日志文件     | 无参数                                            |
| ExportFileHandler       | GET       | 导出日志文件         | `name` - 文件名，`format` - 可选参数（json/text） |

## 响应格式

### 除 ExportFileHandler 外，所有处理器返回统一格式的 JSON 响应：

```json
{
  "code": 200,
  "data": {}, // 或 null
  "msg": "success"
}
```

# 集成示例

```go
package main

import (
	"net/http"
    "github.com/zjguoxin/goslogviewer/adapter"
	"github.com/zjguoxin/goslogviewer/goslogviewer"
)

func main() {
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

	viewer := goslogviewer.New(config)

	// 注册Gin路由
	adapter.RegisterGinRoutes(r, lv)

	// 启动服务
	http.ListenAndServe(":8080", nil)
}
```

## 许可证

[MIT](https://github.com/zjguoxin/goslogviewer/blob/main/LICENSE)© zjguoxin

### 作者

[zjguoxin@163.com](https://github.com/zjguoxin)
