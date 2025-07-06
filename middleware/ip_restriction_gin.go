/**
 * @Author: guxline zjguoxin@163.com
 * @Date: 2025/7/4 23:25:27
 * @LastEditors: guxline zjguoxin@163.com
 * @LastEditTime: 2025/7/4 23:25:27
 * Description:
 * Copyright: Copyright (©) 2025 中易综服. All rights reserved.
 */
package middleware

import (
	"net"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	ipNetCache = make(map[string]*net.IPNet)
	cacheMutex sync.RWMutex
)

// IPRestriction 创建IP限制中间件
func IPRestriction(enable bool, allowedIPs, trustedProxies []string) gin.HandlerFunc {
	// 预编译可信代理CIDR
	compiledProxies := compileCIDRs(trustedProxies)
	compiledAllowed := compileCIDRs(allowedIPs)

	return func(c *gin.Context) {
		// 如果未启用 IP 限制，直接放行
		if !enable {
			c.Next()
			return
		}
		clientIP := getClientIP(c, compiledProxies)

		// c.Set("clientIP", clientIP)
		// ip := c.MustGet("clientIP").(string)
		// fmt.Println("clientIP:", ip)
		if isIPAllowed(clientIP, allowedIPs, compiledAllowed) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(403, gin.H{
			"code":    403,
			"message": "Access denied for IP: " + clientIP,
			"your_ip": clientIP,
		})
	}
}

// 编译CIDR列表为IPNet对象
func compileCIDRs(cidrs []string) []*net.IPNet {
	var result []*net.IPNet
	for _, cidr := range cidrs {
		if strings.Contains(cidr, "/") {
			if ipNet := parseCIDR(cidr); ipNet != nil {
				result = append(result, ipNet)
			}
		}
	}
	return result
}

// 获取客户端真实IP
func getClientIP(c *gin.Context, compiledProxies []*net.IPNet) string {
	// 1. 使用Gin内置方法获取初步IP
	ipStr := c.ClientIP()

	// 2. 特殊处理IPv6本地地址
	if ipStr == "::1" {
		return "127.0.0.1"
	}

	// 3. 处理X-Forwarded-For头部
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		for i := len(ips) - 1; i >= 0; i-- {
			candidate := strings.TrimSpace(ips[i])
			if !isIPInCIDRs(candidate, compiledProxies) {
				ipStr = candidate
				break
			}
		}
	}

	// 4. 去除端口号
	if host, _, err := net.SplitHostPort(ipStr); err == nil {
		ipStr = host
	}

	return ipStr
}

// 检查IP是否被允许
func isIPAllowed(ipStr string, allowedIPs []string, compiledAllowed []*net.IPNet) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// 1. 检查精确匹配
	for _, allowed := range allowedIPs {
		if allowed == ipStr || allowed == "*" {
			return true
		}
	}

	// 2. 检查CIDR范围
	for _, ipNet := range compiledAllowed {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}

// 检查IP是否在CIDR列表中
func isIPInCIDRs(ipStr string, ipNets []*net.IPNet) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, ipNet := range ipNets {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

// 解析CIDR并缓存
func parseCIDR(cidr string) *net.IPNet {
	cacheMutex.RLock()
	ipNet, exists := ipNetCache[cidr]
	cacheMutex.RUnlock()

	if exists {
		return ipNet
	}

	_, ipNetTmp, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	cacheMutex.Lock()
	ipNetCache[cidr] = ipNetTmp
	cacheMutex.Unlock()

	return ipNetTmp
}
